package main

import (
	"bitbucket.org/creativeadvtech/project-template/internal"
	"bitbucket.org/creativeadvtech/project-template/internal/config"
	"bitbucket.org/creativeadvtech/project-template/internal/module"
	"bitbucket.org/creativeadvtech/project-template/pkg/database"
	"bitbucket.org/creativeadvtech/project-template/pkg/logging"
	rest2 "bitbucket.org/creativeadvtech/project-template/pkg/rest"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kelseyhightower/envconfig"
	logs "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

func main() {
	// load app configuration from environment variables
	var cfg config.Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		logs.Fatal(err)
		os.Exit(-1)
	}
	debug := cfg.LogLevel == "debug"
	logging.Init(cfg.LogLevel)

	// set up postgres connection
	logs.Info("Setting up Postgres database connection.")
	db, err := database.NewDatabase(fmt.Sprintf(database.DbURLTemplate, cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName), debug)
	if err != nil {
		logs.Fatal(err)
		os.Exit(-1)
	}
	router := chi.NewRouter()
	router.Use(rest2.RequestLogger(rest2.NewLogger(cfg.LogLevel)))
	router.Use(middleware.Recoverer)
	router.Use(rest2.SentryMiddleware(cfg.SentryDSN, cfg.SentryENV, debug, internal.AppVersion))

	router.Handle("/", http.RedirectHandler("/docs/", http.StatusMovedPermanently))
	router.Get("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./static"))).ServeHTTP)
	router.Get("/status", rest2.APIHandlerFunc(internal.Status(internal.AppVersion)))

	router.Mount("/v1", router.Group(func(r chi.Router) {
		r.Mount("/objects", module.NewRest(module.NewModule(db)))
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
