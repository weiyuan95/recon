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

// Bootstrap creates the necessary tables if they don't exist
// The executed statements must be idempotent
func (s *PostgresTransferStore) Bootstrap() {
	_, err := s.pg.Exec(`
		CREATE TABLE IF NOT EXISTS transfers (
			"id" SERIAL PRIMARY KEY, -- surrogate key
			"chain" TEXT NOT NULL, -- denormalized
			"txid" TEXT NOT NULL,
			"timestamp" TEXT NOT NULL,
			"transferType" TEXT NOT NULL,
			"from" TEXT NOT NULL,
			"to" TEXT NOT NULL,
			"amount" TEXT NOT NULL
		);
		-- Non unique index on chain and txid since 1 tx can have multiple transfers
		CREATE INDEX IF NOT EXISTS "idx_transfers_chain_txid" ON "transfers" ("chain", "txid");

		CREATE TABLE IF NOT EXISTS accounts (
			"id" SERIAL PRIMARY KEY, -- surrogate key
			"chain" TEXT NOT NULL, -- denormalized
			"address" TEXT NOT NULL,
			"createdAt" DATE NOT NULL,
			"lastUpdated" DATE NOT NULL
		);
		CREATE UNIQUE INDEX IF NOT EXISTS "idx_accounts_chain_address" ON "accounts" ("chain", "address");
	`)
	if err != nil {
		log.Fatal("Failed to bootstrap postgres db", err)
	}

	// region Migrations go here
	// ...
	// endregion
}

func (s *PostgresTransferStore) Add(address string, transfer reporter.Transfer) error {
	_, err :=
		s.pg.Exec(`INSERT INTO "transfers" ("chain", "txid", "timestamp", "transferType", "from", "to", "amount") VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			transfer.Chain, transfer.Txid, transfer.Timestamp, transfer.TransferType, transfer.From, transfer.To, transfer.Amount)
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
