package main

import (
	"context"
	"log"

	"github.com/blanc08/stok-gas-management-backend/pkg/api"
	database "github.com/blanc08/stok-gas-management-backend/pkg/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5"
)

func main() {
	// Viper || Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Database
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}
	defer conn.Close(ctx)

	store := database.NewStore(conn)

	restApiServer, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create the server :", err)
	}

	// Validator || bind to server

	err = restApiServer.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the server :", err)
	}

}
