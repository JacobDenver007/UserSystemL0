package user

import (
	"database/sql"
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

func (db *DB) execSQL(sqlStr string) error {
	log.Debugf("ExecSQL: %s", sqlStr)
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
	row := db.sqlDB.QueryRow("SELECT s_username, s_pwd, s_phonenum, i_issuspended, i_auth, i_isapproved from t_userinfo where s_username=?;", name)
	err := row.Scan(&user.UserName, &user.HashPwd, &user.PhoneNum, &user.IsSuspended, &user.Auth, &user.IsApproved)
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
		return nil, fmt.Errorf("this phone did not send verification code")
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
	row := db.sqlDB.QueryRow("SELECT s_address, s_privatekey, i_issuspended, i_isfrozen from t_accountinfo where s_address=?;", address)
	err := row.Scan(&account.Address, &account.PrivateKey, &account.IsSuspended, &account.IsFrozen)
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
