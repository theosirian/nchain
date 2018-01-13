package main

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/kthomas/go-db-config"
)

var (
	migrateOnce sync.Once
)

func migrateSchema() {
	migrateOnce.Do(func() {
		db := dbconf.DatabaseConnection()

		db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
		db.Exec("CREATE EXTENSION IF NOT EXISTS \"pgcrypto\";")

		initial := !db.HasTable(&Network{})

		db.AutoMigrate(&Network{})
		db.Model(&Network{}).AddForeignKey("sidechain_id", "networks(id)", "SET NULL", "CASCADE")

		db.AutoMigrate(&Wallet{})
		db.Model(&Wallet{}).AddIndex("idx_wallets_application_id", "application_id")
		db.Model(&Wallet{}).AddForeignKey("network_id", "networks(id)", "SET NULL", "CASCADE")

		db.AutoMigrate(&Transaction{})
		db.Model(&Transaction{}).AddIndex("idx_transactions_application_id", "application_id")
		db.Model(&Transaction{}).AddForeignKey("network_id", "networks(id)", "SET NULL", "CASCADE")
		db.Model(&Transaction{}).AddForeignKey("wallet_id", "wallets(id)", "SET NULL", "CASCADE")

		db.AutoMigrate(&Contract{})
		db.Model(&Contract{}).AddIndex("idx_contracts_application_id", "application_id")
		db.Model(&Contract{}).AddForeignKey("network_id", "networks(id)", "SET NULL", "CASCADE")
		db.Model(&Contract{}).AddForeignKey("transaction_id", "transactions(id)", "SET NULL", "CASCADE")

		db.AutoMigrate(&Token{})
		db.Model(&Token{}).AddIndex("idx_tokens_application_id", "application_id")
		db.Model(&Token{}).AddForeignKey("network_id", "networks(id)", "SET NULL", "CASCADE")
		db.Model(&Token{}).AddForeignKey("contract_id", "contracts(id)", "SET NULL", "CASCADE")
		db.Model(&Token{}).AddForeignKey("sale_contract_id", "contracts(id)", "SET NULL", "CASCADE")

		if initial {
			populateInitialNetworks()
		}
	})
}

func populateInitialNetworks() {
	db := dbconf.DatabaseConnection()

	var btcMainnet = &Network{}
	db.Raw("INSERT INTO networks (created_at, name, description, is_production) values (NOW(), 'Bitcoin', 'Bitcoin mainnet', true) RETURNING id").Scan(&btcMainnet)

	var btcTestnet = &Network{}
	db.Raw("INSERT INTO networks (created_at, name, description, is_production) values (NOW(), 'Bitcoin testnet', 'Bitcoin testnet', false) RETURNING id").Scan(&btcTestnet)

	var lightningMainnet = &Network{}
	db.Raw("INSERT INTO networks (created_at, name, description, is_production) values (NOW(), 'Lightning Network', 'Lightning Network mainnet', true) RETURNING id").Scan(&lightningMainnet)

	var lightningTestnet = &Network{}
	db.Raw("INSERT INTO networks (created_at, name, description, is_production) values (NOW(), 'Lightning Network testnet', 'Lightning Network testnet', false) RETURNING id").Scan(&lightningTestnet)

	db.Exec("UPDATE networks SET sidechain_id = ? WHERE id = ?", lightningMainnet.Id, btcMainnet.Id)
	db.Exec("UPDATE networks SET sidechain_id = ? WHERE id = ?", lightningTestnet.Id, btcTestnet.Id)

	db.Exec("INSERT INTO networks (created_at, name, description, is_production, config) values (NOW(), 'Ethereum', 'Ethereum mainnet', true, '{\"json_rpc_url\": \"http://ethereum-mainnet-json-rpc.provide.services\"}')")
	db.Exec("INSERT INTO networks (created_at, name, description, is_production, config) values (NOW(), 'Ethereum testnet', 'Ropsten (Revival) testnet', false, '{\"json_rpc_url\": \"http://ethereum-ropsten-testnet-json-rpc.provide.services\", \"testnet\": \"ropsten\"}')")
}

func DatabaseConnection() *gorm.DB {
	return dbconf.DatabaseConnection()
}
