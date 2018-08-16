package user

var DBClient *DB

func Init() {
	db := &DB{}
	if err := db.Open(); err != nil {
		panic(err)
	}
	DBClient = db
}
