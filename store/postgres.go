package store

import (
	"chaintx/reporter"
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

type PostgresTransferStore struct {
	pg *sql.DB
}

func NewPostgresTransferStore() *PostgresTransferStore {
	pg := getPgClient()
	return &PostgresTransferStore{pg}
}

func (s *PostgresTransferStore) Add(address string, transfer reporter.Transfer) error {
	_, err :=
		s.pg.Exec(`INSERT INTO "transfers" ("txid", "timestamp", "transferType", "from", "to", "amount") VALUES ($1, $2, $3, $4, $5, $6)`,
			transfer.Txid, transfer.Timestamp, transfer.TransferType, transfer.From, transfer.To, transfer.Amount)
	return err
}

func (s *PostgresTransferStore) Get(_txid string) *reporter.Transfer {
	rows, err := s.pg.Query(`SELECT "txid", "timestamp", "transferType", "from", "to", "amount" FROM "transfers" WHERE "txid" = $1`, _txid)
	if err != nil {
		log.Fatal(err)
	}

	if !rows.Next() {
		return nil
	}

	var transfer reporter.Transfer
	if err := rows.Scan(&transfer.Txid, &transfer.Timestamp, &transfer.TransferType, &transfer.From, &transfer.To, &transfer.Amount); err != nil {
		log.Fatal(err)
	}
	log.Println("Found transfer", transfer)
	return &transfer
}

var pgClient *sql.DB

func getPgClient() *sql.DB {
	if pgClient != nil {
		return pgClient
	}

	connStr := "user=postgres password=postgres dbname=chaintx sslmode=disable" // TODO: use env creds
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		// We panic here since we definitely need the client to be connected on start,
		// and there's no recovering from this
		log.Fatal(err)
	}
	pgClient = db
	return db
}

func BootstrapTables(client *sql.DB) {
	_, err := client.Exec(`
		CREATE TABLE IF NOT EXISTS transfers (
			-- TODO: address
			"txid" TEXT PRIMARY KEY,
			"timestamp" TEXT NOT NULL,
			"transferType" TEXT NOT NULL,
			"from" TEXT NOT NULL,
			"to" TEXT NOT NULL,
			"amount" TEXT NOT NULL,
		);

		CREATE TABLE IF NOT EXISTS accounts (
			"chain" TEXT NOT NULL,
			"address" TEXT PRIMARY KEY,
			"state" TEXT NOT NULL,
			"createdAt" DATE NOT NULL,
			"lastUpdated" DATE NOT NULL
		);
	`)
	if err != nil {
		log.Fatal("Failed to bootstrap postgres db", err)
	}
}
