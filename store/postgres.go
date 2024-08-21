package store

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
	"recon/reporter"
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
			"address" TEXT NOT NULL, 
			"txid" TEXT NOT NULL,
			"timestamp" TEXT NOT NULL,
			"transferType" TEXT NOT NULL,
			"tokenType" TEXT NOT NULL,
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

func (s *PostgresTransferStore) Add(transfer reporter.Transfer) error {
	_, err :=
		s.pg.Exec(`INSERT INTO "transfers" ("chain", "address", "txid", "timestamp", "transferType", "tokenType", "from", "to", "amount") VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
			transfer.Chain, transfer.Address, transfer.Txid, transfer.Timestamp, transfer.TransferType, transfer.TokenType, transfer.From, transfer.To, transfer.Amount)
	return err
}

func (s *PostgresTransferStore) Get(_txid string) *reporter.Transfer {
	rows, err := s.pg.Query(`SELECT "chain", "address", "txid", "timestamp", "transferType", "tokenType", "from", "to", "amount" FROM "transfers" WHERE "txid" = $1`, _txid)
	if err != nil {
		log.Println(err) // TODO: error logging
		return nil
	}

	if !rows.Next() {
		return nil
	}

	transfer, err := rowToTransfer(rows)
	if err != nil {
		log.Println(err)
	}
	return transfer
}

func (s *PostgresTransferStore) ListByAddress(address string) []reporter.Transfer {
	rows, err := s.pg.Query(`SELECT "chain", "address", "txid", "timestamp", "transferType", "tokenType", "from", "to", "amount" FROM "transfers" WHERE "address" = $1`, address)
	if err != nil {
		log.Println("Error encountered executing query", err)
		return nil
	}

	transfers := make([]reporter.Transfer, 0)

	for rows.Next() {
		transfer, err := rowToTransfer(rows)
		if err != nil {
			log.Println("Error encountered mapping row to Transfer struct", err)
			continue
		}
		transfers = append(transfers, *transfer)
	}

	return transfers
}

func rowToTransfer(rows *sql.Rows) (*reporter.Transfer, error) {
	var transfer reporter.Transfer
	if err := rows.Scan(&transfer.Chain, &transfer.Address, &transfer.Txid, &transfer.Timestamp, &transfer.TransferType, &transfer.TokenType, &transfer.From, &transfer.To, &transfer.Amount); err != nil {
		return nil, err
	}
	return &transfer, nil
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
