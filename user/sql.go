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

CREATE TABLE IF NOT EXISTS t_accountinfo (
	id INTEGER(11)       PRIMARY KEY AUTO_INCREMENT,
	s_user               VARCHAR(42) NOT NULL,
	s_address            VARCHAR(42) NOT NULL,
	s_privatekey         longtext,
	i_issuspended        tinyint(1) default 0,
	i_isfrozen           tinyint(1) default 0,
	INDEX uniq_user (s_user),
	UNIQUE KEY uniq_address (s_address)
);

CREATE TABLE IF NOT EXISTS t_mainchain (
	id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	i_height int(11) NOT NULL,
	s_hash varchar(100) NOT NULL,
	i_tx_count int(11) NOT NULL,
	i_tx_index int(11) NOT NULL,
	i_created int(11) NOT NULL,
	UNIQUE KEY uniq_row(i_height)
);

CREATE TABLE IF NOT EXISTS t_txs (
	id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	s_hash varchar(100) NOT NULL,
	i_nonce int(11) NOT NULL,
	s_from_addr longtext NOT NULL,
	s_to_addr longtext NOT NULL,
	s_value varchar(100) NOT NULL,
	i_fee bigint(20) NOT NULL,
	i_height int(11) NOT NULL,
	i_created int(11) NOT NULL,
	s_error longtext
); 

CREATE TABLE IF NOT EXISTS t_address (
	s_address varchar(100) NOT NULL,
	i_tx_counts int(11) NOT NULL,
	s_mount varchar(100) NOT NULL,
	i_nonce int(11) NOT NULL,
	PRIMARY KEY (s_address)
);

CREATE TABLE IF NOT EXISTS t_history (
	id int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
	s_address varchar(100) NOT NULL,
	s_txs longtext NOT NULL,
	UNIQUE KEY uniq_address(s_address)
); 
`
