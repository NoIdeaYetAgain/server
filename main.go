package main

import (
	"clove/api"
	db "clove/internal/db/sqlc"
	"clove/util"
	"clove/worker"
	"context"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog/log"
	stdlog "log"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		stdlog.Fatal("cannot load config:", err)
	}

	conn, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		stdlog.Fatal("cannot connect to db", err)
	}

	store := db.NewStore(conn)

	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}
	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	runTaskProcessor(redisOpt, *store)
	server, err := api.NewServer(config, store, taskDistributor) // Pass both config and store
	if err != nil {
		stdlog.Fatal("cannot create server:", err) // Handle the error from NewServer
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		stdlog.Fatal("cannot start server:", err)
	}
}

func runTaskProcessor(redisOpt asynq.RedisClientOpt, store db.SQLStore) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpt, store)
	log.Info().Msg("start task processor")
	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server")
	}
}
