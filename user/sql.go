package user

// -- address char(34)
// -- hash char(64)
// -- value bigint
//BTC 区块hash存在相同的

var (
	DBUser = "root"
	DBPwd  = "Bochen@123"
	DBHost = "http://127.0.0.1:3306"
	DBName = "chain"
)

var createSQL = `
CREATE TABLE IF NOT EXISTS t_userinfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_username           VARCHAR(42) NOT NULL,
	s_pwd                VARCHAR(100),
	s_phonenum           VARCHAR(42),
	UNIQUE KEY uniq_username (s_username)
);

CREATE TABLE IF NOT EXISTS t_tokeninfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_verificationcode   VARCHAR(42) NOT NULL,
	s_phonenum           VARCHAR(42),
	s_timestamp          VARCHAR(42),
	UNIQUE KEY uniq_phonenum (s_phonenum)
);
`
