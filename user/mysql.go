package user

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/JacobDenver007/UserSystemL0/common"
	"github.com/JacobDenver007/UserSystemL0/log"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	sqlDB *sql.DB
}

type Address struct {
	address string
	count   uint32
	mount   string
	nonce   uint32
}

type History struct {
	txIds []uint64
}

func (db *DB) execSQL(sqlStr string) error {
	//log.Debugf("ExecSQL: %s", sqlStr)
	sqlStrs := strings.Split(sqlStr, ";")
	tx, err := db.sqlDB.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if tx != nil {
			tx.Rollback()
		}
	}()
	for _, sqlStr := range sqlStrs {
		sqlStr = strings.TrimSpace(sqlStr)
		if len(sqlStr) != 0 {
			if _, err := tx.Exec(fmt.Sprintf("%s;", sqlStr)); err != nil {
				return fmt.Errorf("%s - %s", sqlStr, err)
			}
		}
	}
	err = tx.Commit()
	if err == nil {
		tx = nil
	}
	return err
}

func (db *DB) Open() error {
	//mysql
	sdb, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&loc=%s&parseTime=true",
		DBUser, DBPwd, DBHost, DBName, url.QueryEscape("Asia/Shanghai")))
	if err != nil {
		return err
	}
	sdb.SetMaxOpenConns(2000)
	sdb.SetMaxIdleConns(2000)
	sdb.SetConnMaxLifetime(60 * time.Second)
	db.sqlDB = sdb
	if err := db.execSQL(createSQL); err != nil {
		sdb.Close()
		return err
	}
	return nil
}

func (db *DB) IfExistUserName(name string) bool {
	var username string
	row := db.sqlDB.QueryRow("SELECT s_username from t_userinfo where s_username=?", name)
	if err := row.Scan(&username); err == sql.ErrNoRows {
		return false
	}
	return true
}

func (db *DB) InsertUser(user *common.User) error {
	sqlStr := fmt.Sprintf("INSERT INTO t_userinfo(s_username, s_pwd, s_phonenum) values('%s','%s','%s');",
		user.UserName, user.HashPwd, user.PhoneNum)

	return db.execSQL(sqlStr)
}

func (db *DB) GetUserInfo(name string) (*common.User, error) {
	user := &common.User{}
	row := db.sqlDB.QueryRow("SELECT id, s_username, s_pwd, s_phonenum, i_issuspended, i_auth, i_isapproved from t_userinfo where s_username=?;", name)
	err := row.Scan(&user.ID, &user.UserName, &user.HashPwd, &user.PhoneNum, &user.IsSuspended, &user.Auth, &user.IsApproved)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("用户不存在")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) UpdateUserInfo(user *common.User) error {
	sqlStr := fmt.Sprintf("UPDATE t_userinfo set s_username='%s', s_pwd='%s', s_phonenum='%s', i_issuspended=%d, i_auth=%d, i_isapproved=%d where id=%d;",
		user.UserName, user.HashPwd, user.PhoneNum, user.IsSuspended, user.Auth, user.IsApproved, user.ID)
	return db.execSQL(sqlStr)
}

func (db *DB) InsertToken(phoneNum string, verificationCode string) error {
	sqlStr := fmt.Sprintf("INSERT INTO t_verificationinfo(s_verificationcode, s_phonenum, i_timestamp) values('%s','%s',%d);",
		verificationCode, phoneNum, time.Now().Unix())

	return db.execSQL(sqlStr)
}

func (db *DB) GetTokenInfo(phoneNum string) (*common.TokenInfo, error) {
	token := &common.TokenInfo{}
	row := db.sqlDB.QueryRow("SELECT s_phonenum, s_verificationcode, i_timestamp from t_verificationinfo where s_phonenum=? order by id desc limit 1;", phoneNum)
	err := row.Scan(&token.PhoneNum, &token.VerificationCode, &token.TimeStamp)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("此手机号未发送过验证码，请检查手机号是否正确")
	}
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (db *DB) InsertAccount(user string, address string, privateKey string) error {
	sqlStr := fmt.Sprintf("INSERT INTO t_accountinfo(s_user, s_address, s_privatekey) values('%s','%s','%s');",
		user, address, privateKey)

	return db.execSQL(sqlStr)
}

func (db *DB) GetAccountInfo(address string) (*common.Account, error) {
	account := &common.Account{}
	row := db.sqlDB.QueryRow("SELECT s_user, s_address, s_privatekey, i_issuspended, i_isfrozen from t_accountinfo where s_address=?;", address)
	err := row.Scan(&account.User, &account.Address, &account.PrivateKey, &account.IsSuspended, &account.IsFrozen)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("account does not exist")
	}
	if err != nil {
		return nil, err
	}
	return account, nil
}

func (db *DB) UpdateAccountInfo(account *common.Account) error {
	sqlStr := fmt.Sprintf("UPDATE t_accountinfo set s_address='%s', s_privatekey='%s', i_issuspended=%d, i_isfrozen=%d where s_address='%s';",
		account.Address, account.PrivateKey, account.IsSuspended, account.IsFrozen, account.Address)
	return db.execSQL(sqlStr)
}

func (db *DB) GetUserAccount(user string) ([]string, error) {
	sqlStr := fmt.Sprintf("SELECT s_address from t_accountinfo where s_user='%s';", user)
	rows, err := db.sqlDB.Query(sqlStr)
	defer rows.Close()
	if err != nil {
		return nil, err
	}
	accounts := make([]string, 0)
	for rows.Next() {
		var address string
		rows.Scan(&address)
		accounts = append(accounts, address)
	}
	return accounts, nil
}

func (db *DB) DeleteAccountInfo(address string) error {
	sqlStr := fmt.Sprintf("DELETE from t_accountinfo where s_address='%s';", address)
	return db.execSQL(sqlStr)
}

func (db *DB) GetBestBlock() (*Block, error) {
	block := &Block{}
	var height uint32
	sqlstr := "SELECT i_height FROM t_mainchain order by i_height desc limit 1"
	row := db.sqlDB.QueryRow(sqlstr)
	err := row.Scan(&height)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	blockHeader := &BlockHeader{}
	blockHeader.Height = height
	block.Header = blockHeader
	return block, nil
}

func (db *DB) InsertBlock(block *Block) error {
	txIndex, _ := db.getTxIndex(block.Header.Height - 1)

	tmpTxIndex := txIndex

	var sqlStr string
	address := make(map[string]*Address)
	history := make(map[string]*History)

	for _, tx := range block.Txs {
		tTxName := "t_txs"

		txID := (txIndex + 1) % 10000000
		sqlStr += fmt.Sprintf("INSERT INTO %s(id, s_hash, i_nonce, s_from_addr, s_to_addr, s_value, i_fee, i_height, i_created) values(%d, '%s',%d,'%s','%s','%s','%s',%d,%d);", tTxName, txID, tx.Hash, tx.Data.Nonce, tx.Data.Sender, tx.Data.Recipient, tx.Data.Amount.String(), tx.Data.Fee.String(), block.Header.Height, tx.Data.CreateTime)

		address = db.formAddress(address, tx.Data.Sender, tx.Data.Nonce)
		history = db.formHistory(history, tx.Data.Sender, txIndex)
		if tx.Data.Sender != tx.Data.Recipient {
			address = db.formAddress(address, tx.Data.Recipient, 0)
			history = db.formHistory(history, tx.Data.Recipient, txIndex)
		}
		txIndex++
	}

	sqlStr += db.updateAddress(address)
	sqlStr += db.updateHistory(history)

	sqlStr += fmt.Sprintf("INSERT INTO t_mainchain(i_height, s_hash, i_tx_count, i_tx_index, i_created) values(%d,'%s', %d, %d, %d);",
		block.Header.Height, "", len(block.Txs), tmpTxIndex, block.Header.TimeStamp)
	return db.execSQL(sqlStr)
}

func (db *DB) getTxIndex(blockNumber uint32) (uint64, error) {
	row := db.sqlDB.QueryRow("SELECT i_tx_count, i_tx_index FROM t_mainchain where i_height=?", blockNumber)
	var txCount, txIndex uint64
	err := row.Scan(&txCount, &txIndex)
	if err == sql.ErrNoRows {
		log.Errorf("db.getTxIndex: ERROR, can not get txindex %d", blockNumber)
		return 0, nil
	}
	if err != nil {
		log.Errorf("db.getTxIndex: ERROR, block:%d %s", blockNumber, err)
		panic(err)
	}
	return txCount + txIndex, nil
}

func (db *DB) formHistory(history map[string]*History, key string, txIndex uint64) map[string]*History {
	if v, ok := history[key]; !ok {
		txIDs := make([]uint64, 0)
		history[key] = &History{txIds: append(txIDs, txIndex)}
	} else {
		history[key] = &History{txIds: append(v.txIds, txIndex)}
	}
	return history
}

func (db *DB) updateHistory(history map[string]*History) string {

	var sqlStr string

	for key, v := range history {
		tName := "t_history"
		isExist, _, _, _ := db.getAddressInfo(key)

		if !isExist {
			addressRow := key
			txs, _ := json.Marshal(v.txIds)
			sqlStr += fmt.Sprintf("REPLACE INTO %s(s_address, s_txs) values('%s','%s');", tName, addressRow, string(txs))
		} else {
			addressRow := key
			id, txs, _ := db.getHistory(tName, addressRow)
			txsID := make([]uint64, 0)
			err := json.Unmarshal([]byte(txs), &txsID)
			if err != nil {
				log.Errorf("updateHistory Unmarshal fail %s", err.Error())
			}
			for _, txIndex := range v.txIds {
				txsID = append(txsID, txIndex)
			}
			newTxs, _ := json.Marshal(txsID)
			sqlStr += fmt.Sprintf("REPLACE INTO %s(id,s_address, s_txs) values(%d,'%s','%s');", tName, id, addressRow, string(newTxs))
		}
	}
	return sqlStr
}

func (db *DB) formAddress(address map[string]*Address, hash string, inNonce uint32) map[string]*Address {
	_, dbcount, nonce, _ := db.getAddressInfo(hash)
	if inNonce > nonce {
		nonce = inNonce
	}

	if v, ok := address[hash]; !ok {
		address[hash] = &Address{address: hash, count: dbcount + 1, mount: "0", nonce: nonce}
	} else {
		address[hash] = &Address{address: hash, count: v.count + 1, mount: "0", nonce: nonce}
	}
	return address
}

func (db *DB) updateAddress(address map[string]*Address) string {
	var sqlStr string

	for k, v := range address {
		sqlStr += fmt.Sprintf("REPLACE INTO t_address(s_address, i_tx_counts, s_mount, i_nonce) values('%s', %d, '%s', %d);",
			k, v.count, "0", v.nonce)
	}
	return sqlStr
}

func (db *DB) getAddressInfo(addressHash string) (bool, uint32, uint32, string) {
	row := db.sqlDB.QueryRow("SELECT i_tx_counts, i_nonce, s_mount FROM t_address where s_address=?", addressHash)
	var count, nonce uint32
	var amount string
	err := row.Scan(&count, &nonce, &amount)
	if err == sql.ErrNoRows {
		log.Debugf("db.getAddressInfo: can not find address %s", addressHash)
		return false, 0, 0, "0"
	}
	if err != nil {
		log.Errorf("db.getAddressInfo: ERROR, address: %s %s", addressHash, err)
		panic(err)
	}
	return true, count, nonce, amount
}

func (db *DB) getHistory(tName string, addressRow string) (uint64, string, error) {
	sqlStr := fmt.Sprintf("SELECT id, s_txs FROM %s where s_address='%s'", tName, addressRow)
	row := db.sqlDB.QueryRow(sqlStr)
	var txs string
	var id uint64
	err := row.Scan(&id, &txs)
	if err != nil {
		log.Errorf("mysql.getHistory: ERROR, address: %s %s", addressRow, err)
		panic(err)
	}
	return id, txs, nil
}

func (db *DB) GetHistory(addr string) ([]*Transaction, error) {
	var reslutTxs []*Transaction

	dbTxs, err := db.getHistoryFromDB(addr)
	if err != nil {
		return nil, err
	}
	reslutTxs = append(reslutTxs, dbTxs...)
	return reslutTxs, err
}

func (db *DB) getHistoryFromDB(addr string) ([]*Transaction, error) {
	tName := "t_history"
	isExist, txcount, _, _ := db.getAddressInfo(addr)

	//mysql not found txcount by addr, contractAddr
	if !isExist {
		return nil, nil
	}

	var (
		txs []*Transaction
	)

	for i := (txcount - 1) / 1000; i >= 0; i-- {
		addressRow := addr
		txIDsStr, err := db.getTxIDs(tName, addressRow)
		if err != nil {
			return nil, err
		}
		txIDs := make([]uint64, 0)
		json.Unmarshal([]byte(txIDsStr), &txIDs)

		for j := len(txIDs) - 1; j >= 0; j-- {
			tx, err := db.getTxByID(txIDs[j])
			if err != nil {
				return nil, err
			}
			txs = append(txs, tx)
		}
	}

	return txs, nil
}

func (db *DB) getTxIDs(tName string, addressRow string) (string, error) {
	sqlStr := fmt.Sprintf("SELECT s_txs FROM %s where s_address='%s'", tName, addressRow)
	row := db.sqlDB.QueryRow(sqlStr)
	var txs string
	err := row.Scan(&txs)
	if err != nil {
		return "", err
	}
	return txs, nil
}

func (db *DB) getTxByID(id uint64) (*Transaction, error) {
	sqlStr := fmt.Sprintf("select s_hash,s_from_addr,s_to_addr,s_value,i_fee,i_nonce,i_height,i_created from %s where id=%d ", "t_txs", id)
	row := db.sqlDB.QueryRow(sqlStr)

	tx := &Transaction{}
	var value, fee string
	err := row.Scan(&tx.Hash, &tx.Data.Sender, &tx.Data.Recipient, &value, &fee, &tx.Data.Nonce, &tx.Data.BlockNumber, &tx.Data.CreateTime)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	tx.Data.Amount.SetString(value, 10)
	tx.Data.Fee.SetString(fee, 10)
	return tx, nil
}
