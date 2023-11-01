package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"pos/permission"
	"pos/role"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {
	// load configuration
	var configFileName string
	flag.StringVar(&configFileName, "c", "config.yml", "Config file name")
	flag.Parse()

	cfg := loadConfig(configFileName)
	log.Debug().Any("config", cfg).Msg("config loaded")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, cfg.DBCfg.ConnStr())
	if err != nil {
		log.Error().Err(err).Msg("unable to connect to database")
	}

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	permissionRepo := permission.NewRepo(pool)
	permissionReadModel := permission.NewReadModel(pool)
	roleRepo := role.NewRepo(pool)
	roleReadModel := role.NewReadModel(pool)

	readDataPermission := permission.NewReadData(
		permissionRepo,
		permissionReadModel,
	)
	mutateDataPermission := permission.NewMutationData(
		permissionRepo,
		permissionReadModel,
	)

	readDataRole := role.NewReadData(
		roleRepo,
		roleReadModel,
	)
	mutateDataRole := role.NewMutationData(
		roleRepo,
		roleReadModel,
	)

	permissionRoute := permission.NewRoute(
		mutateDataPermission,
		readDataPermission,
	)

	roleRoute := role.NewRoute(
		mutateDataRole,
		readDataRole,
	)

	permissionRoute.Routes(r)
	roleRoute.Routes(r)

	log.Info().Msg(fmt.Sprintf("starting up server on: %s", cfg.Listen.Addr()))
	server := &http.Server{
		Handler:      r,
		Addr:         cfg.Listen.Addr(),
		ReadTimeout:  time.Second * time.Duration(cfg.Listen.ReadTimeout),
		WriteTimeout: time.Second * time.Duration(cfg.Listen.WriteTimeout),
		IdleTimeout:  time.Second * time.Duration(cfg.Listen.IdleTimeout),
	}
	if err := server.ListenAndServe(); err != nil {
		log.Fatal().Err(err).Msg("failed to start the server")
		return
	}
	log.Info().Msg("server stop")
}
