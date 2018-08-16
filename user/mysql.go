package user

import (
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/UserSystemL0/common"
	"github.com/UserSystemL0/log"
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

func (db *DB) GetUserInfo(name string, phoneNum string) (*common.User, error) {
	user := &common.User{}
	row := db.sqlDB.QueryRow("SELECT id, s_username, s_pwd, s_phonenum from t_userinfo where s_username=? and s_phonenum=?;", name, phoneNum)
	err := row.Scan(&user.ID, &user.UserName, &user.HashPwd, &user.PhoneNum)
	if err == sql.ErrNoRows {
		return nil, fmt.Errorf("user does not exist")
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DB) UpdateUserInfo(user *common.User) error {
	sqlStr := fmt.Sprintf("UPDATE t_userinfo set s_username='%s', s_pwd='%s', s_phonenum='%s' where id=%d;",
		user.UserName, user.HashPwd, user.PhoneNum, user.ID)
	return db.execSQL(sqlStr)
}

func (db *DB) GetTokenInfo(phoneNum string, verificationode string) (*common.TokenInfo, error) {
	var timestamp string
	row := db.sqlDB.QueryRow("SELECT s_timestamp from t_tokeninfo where s_phonenum=? and s_verificationcode=? order by id desc limit 1;", phoneNum, verificationode)
	if err := row.Scan(&timestamp); err != nil {
		return nil, err
	}
	token := &common.TokenInfo{
		PhoneNum:         phoneNum,
		VerificationCode: verificationode,
		TimeStamp:        timestamp,
	}
	return token, nil
}

func (db *DB) InsertUser(user *common.User) error {
	sqlStr := fmt.Sprintf("INSERT INTO t_userinfo(s_username, s_pwd, s_phonenum) values('%s','%s','%s');",
		user.UserName, user.HashPwd, user.PhoneNum)

	return db.execSQL(sqlStr)
}
