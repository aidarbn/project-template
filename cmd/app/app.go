package main

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	"bitbucket.org/creativeadvtech/project-template/internal/config"
	"bitbucket.org/creativeadvtech/project-template/internal/object-module"
	"bitbucket.org/creativeadvtech/project-template/internal/repositories"
	"bitbucket.org/creativeadvtech/project-template/pkg/database"
	"bitbucket.org/creativeadvtech/project-template/pkg/logging"
	"bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
	logs "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// load app configuration from environment variables
	var cfg config.Config
	if err := envconfig.Process("", &cfg); err != nil {
		logs.Fatal(err)
		os.Exit(-1)
	}
	debug := cfg.LogLevel == "debug"
	logging.Init(cfg.LogLevel)

	// set up postgres connection
	logs.Info("Setting up Postgres database connection.")
	db, err := database.NewDatabase(fmt.Sprintf(database.URLTemplate, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName), debug)
	if err != nil {
		logs.Fatal(err)
		os.Exit(-1)
	}
	// perform schema migration
	logs.Info("Performing schema migration.")
	if err = database.MigrateDBSchema(db.SqlDB(), "./migrations", cfg.DBName); err != nil {
		logs.Errorf("Can't perform migration; error: %v", err)
	}

	// configure router
	router := chi.NewRouter()
	router.Use(rest.RequestLogger(rest.NewLogger(cfg.LogLevel)))
	router.Use(middleware.Recoverer)
	router.Use(rest.SentryMiddleware(cfg.SentryDSN, cfg.SentryENV, debug, internal.AppVersion))

	router.Handle("/", http.RedirectHandler("/docs/", http.StatusMovedPermanently))
	router.Get("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./static"))).ServeHTTP)
	router.Get("/status", rest.APIHandlerFunc(internal.Status(internal.AppVersion)))

	router.Mount("/v1", router.Group(func(r chi.Router) {
		r.Mount("/objects", object_module.NewRest(object_module.NewModule(repositories.NewObjectRepository(db))))
	}))

	// initialize server
	httpServer := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.ServerPort),
		Handler:      router,
		WriteTimeout: cfg.ServerWriteTimeout,
		ReadTimeout:  cfg.ServerReadTimeout,
	}

	logs.Info("Serving the web application...")
	logs.Fatal(httpServer.ListenAndServe())
}
