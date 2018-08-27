package user

// -- address char(34)
// -- hash char(64)
// -- value bigint
//BTC 区块hash存在相同的

var (
	DBUser = "root"
	DBPwd  = "Bochen@123"
	DBHost = "127.0.0.1:3306"
	DBName = "chain"
)

var createSQL = `
CREATE TABLE IF NOT EXISTS t_userinfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_username           VARCHAR(42) NOT NULL,
	s_pwd                VARCHAR(100),
	s_phonenum           VARCHAR(11),
	i_issuspended        TINYINT(1) default 0,
	i_auth               TINYINT(1) default 0,
	i_isapproved         TINYINT(1) default 0,
	UNIQUE KEY uniq_username (s_username),
	INDEX index_phonenum (s_phonenum)
);

CREATE TABLE IF NOT EXISTS t_verificationinfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_verificationcode   VARCHAR(4) NOT NULL,
	s_phonenum           VARCHAR(11) NOT NULL,
	i_timestamp          INTEGER(11),
	UNIQUE KEY uniq_verificationcode (s_verificationcode),
	INDEX index_phonenum (s_phonenum)
);

CREATE TABLE IF NOT EXISTS t_accoutninfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_address            VARCHAR(42) NOT NULL,
	s_privatekey         longtext,
	i_issuspended        tinyint(1) default 0,
	i_isfrozen           tinyint(1) default 0,
	UNIQUE KEY uniq_address (s_address),
	INDEX uniq_privatekey (s_privatekey)
);
`
