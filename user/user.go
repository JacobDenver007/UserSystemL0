package user

var DBClient *DB
var RPCClient *RPC

func Init() {
	db := &DB{}
	if err := db.Open(); err != nil {
		panic(err)
	}
	DBClient = db

	RPCClient = &RPC{
		rpchost: "http://127.0.0.1:8881",
	}

	Scanning(nil, nil)
}
