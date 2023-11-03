package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"pos/internal/account"
	custommiddleware "pos/internal/custom_middleware"
	"pos/internal/oauth"
	"pos/internal/permission"
	"pos/internal/protected"
	"pos/internal/role"
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
	custommiddleware.SetJwtSecret(cfg.JwtCfg.Secret)

	r := chi.NewRouter()
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.CleanPath)

	permissionRepo := permission.NewRepo(pool)
	permissionReadModel := permission.NewReadModel(pool)
	roleRepo := role.NewRepo(pool)
	roleReadModel := role.NewReadModel(pool)
	accountRepo := account.NewRepo(pool)
	accountReadModel := account.NewReadModel(pool)
	oauthRepo := oauth.NewRepo(pool)
	oauthReadModel := oauth.NewReadModel(pool)
	rolePermissionRepo := role.NewRepoRolePermission(pool)
	rolePermissionReadModel := role.NewReadModelRolePermission(pool)
	accountRoleRepo := account.NewRepoAccountRole(pool)
	accountRoleReadModel := account.NewReadModelAccountRole(pool)

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
		rolePermissionRepo,
		rolePermissionReadModel,
	)
	mutateDataRole := role.NewMutationData(
		roleRepo,
		roleReadModel,
	)
	readDataAccount := account.NewReadData(
		accountRepo,
		accountReadModel,
	)
	mutateDataAccount := account.NewMutationData(
		accountRepo,
		accountReadModel,
	)
	oauthSvc := oauth.NewServiceOAuth(
		accountReadModel,
		oauthRepo,
		oauthReadModel,
		cfg.JwtCfg.Secret,
		cfg.JwtCfg.RefreshExpTime,
		cfg.JwtCfg.AccessExpTime,
		roleReadModel,
		permissionReadModel,
	)

	rolePermissionSvc := role.NewRolePermissionService(
		roleRepo,
		roleReadModel,
		rolePermissionRepo,
		rolePermissionReadModel,
	)
	accountRoleSvc := account.NewAccountRoleService(
		accountRepo,
		accountReadModel,
		accountRoleRepo,
		accountRoleReadModel,
	)

	permissionRoute := permission.NewRoute(
		mutateDataPermission,
		readDataPermission,
	)
	roleRoute := role.NewRoute(
		mutateDataRole,
		readDataRole,
		rolePermissionSvc,
	)
	rolePermissionRoute := role.NewPermissionRoute(
		mutateDataRole,
		readDataRole,
		rolePermissionSvc,
	)
	accountRoute := account.NewRoute(
		mutateDataAccount,
		readDataAccount,
	)
	accountRoleRoute := account.NewRoleRoute(
		mutateDataAccount,
		readDataAccount,
		accountRoleSvc,
	)
	oauthRoute := oauth.NewRoute(
		oauthSvc,
		cfg.JwtCfg.Secret,
	)
	r.Mount("/api", oauthRoute.Routes())

	r.Group(func(r chi.Router) {
		r.Use(custommiddleware.AuthJwtMiddleware)
		r.Mount("/api/permission", permissionRoute.Routes())
		r.Mount("/api/role-permission", rolePermissionRoute.Routes())
		r.Mount("/api/role", roleRoute.Routes())
		r.Mount("/api/account", accountRoute.Routes())
		r.Mount("/api/account-role", accountRoleRoute.Routes())
		r.Mount("/api/dashboard", protected.Routes())
	})

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
