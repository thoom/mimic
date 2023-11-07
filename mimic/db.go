package mimic

import (
	"database/sql"
	"log"
)

type MimicDB struct {
	DB *sql.DB
}

func (mdb *MimicDB) CreateSchema() {
	_, err := mdb.DB.Exec(`CREATE TABLE IF NOT EXISTS nonce (
		nonce TEXT PRIMARY KEY, 
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP, 
		used_at DATETIME DEFAULT NULL
		) WITHOUT ROWID;
	`)

	if err != nil {
		log.Fatal(err)
	}
}

func (mdb *MimicDB) CreateNonce() string {
	nonce := RandString(32)

	// Clean up old nonces (makeshift garbage collection)
	mdb.DB.Exec("DELETE FROM nonce WHERE used_at < datetime('now', '-1 week')")
	mdb.DB.Exec("DELETE FROM nonce WHERE created_at < datetime('now', '-1 hour') AND used_at IS NULL")
	mdb.DB.Exec("INSERT INTO nonce (nonce) VALUES (?)", nonce)

	return nonce
}
