package main

import (
	"banking-ledger/app"
	"banking-ledger/server"
)

func main() {
	// running the app
	banking_ledger := app.Initialize(server.New())
	banking_ledger.Start()
}
