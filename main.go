package main

import (
	"context"
	"log"

	"github.com/blanc08/stok-gas-management-backend/pkg/api"
	database "github.com/blanc08/stok-gas-management-backend/pkg/database/sqlc"
	"github.com/blanc08/stok-gas-management-backend/pkg/util"

	pgxuuid "github.com/jackc/pgx-gofrs-uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Viper || Load config
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config:", err)
	}

	// Database
	pgxConfig, err := pgxpool.ParseConfig(config.DBSource)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	pgxConfig.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		// do something with every new connection
		pgxuuid.Register(conn.TypeMap())
		return nil
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatal("Cannot connect to the database", err)
	}

	defer pool.Close()

	// ctx := context.Background()
	// conn, err := pgx.Connect(ctx, config.DBSource)
	// if err != nil {
	// 	log.Fatal("Cannot connect to the database", err)
	// }
	// defer conn.Close(ctx)

	store := database.NewStore()

	restApiServer, err := api.NewServer(config, store, pool)
	if err != nil {
		log.Fatal("cannot create the server :", err)
	}

	// Validator || bind to server

	err = restApiServer.Start(config.HTTPServerAddress)
	if err != nil {
		log.Fatal("cannot start the server :", err)
	}

}
