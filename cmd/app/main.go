package main

import (
	"context"
	"errors"
	"expvar"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/garnizeH/dimdim/http/api/debug"
	server "github.com/garnizeH/dimdim/http/web"
	"github.com/garnizeH/dimdim/pkg/argon2id"
	"github.com/garnizeH/dimdim/pkg/logger"
	"github.com/garnizeH/dimdim/pkg/mailer"
	"github.com/garnizeH/dimdim/pkg/web"
	"github.com/garnizeH/dimdim/storage"
	"github.com/garnizeH/dimdim/storage/datastore"

	"github.com/ardanlabs/conf/v3"
)

const prefix = "DIMDIM"

// TODO: make the makefile build to insert this build value getting the git hash
var build = "develop"

func main() {
	// -------------------------------------------------------------------------
	// Logger Support

	var log *logger.Logger

	events := logger.Events{
		Error: func(ctx context.Context, r logger.Record) {
			log.Info(ctx, "******* SEND ALERT *******")
		},
	}

	traceIDFn := func(ctx context.Context) string {
		return web.GetTraceID(ctx)
	}

	log = logger.NewWithEvents(os.Stdout, logger.LevelInfo, prefix, traceIDFn, events)

	// -------------------------------------------------------------------------
	// Run

	ctx := context.Background()

	if err := run(ctx, log); err != nil {
		log.Error(ctx, "startup", "msg", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, log *logger.Logger) error {
	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Info(ctx, "startup", "GOMAXPROCS", runtime.GOMAXPROCS(0))

	// -------------------------------------------------------------------------
	// Configuration

	cfg := struct {
		conf.Version
		Web struct {
			AppName            string        `conf:"default:dimdim"`
			DomainName         string        `conf:"default:localhost"`
			Port               string        `conf:"default:3000"`
			BindAddress        string        `conf:"default:0.0.0.0"`
			ReadTimeout        time.Duration `conf:"default:5s"`
			WriteTimeout       time.Duration `conf:"default:10s"`
			IdleTimeout        time.Duration `conf:"default:120s"`
			ShutdownTimeout    time.Duration `conf:"default:20s"`
			CORSAllowedOrigins []string      `conf:"default:*,mask"`
		}
		Debug struct {
			Host string `conf:"default:0.0.0.0:3010"`
		}
		DB struct {
			DSN string `conf:"default:dimdim.db"`
		}
		Mailer struct {
			Addr     string `conf:"default:localhost:1025"`
			Identity string `conf:"default:dimdim@localhost"`
			Username string `conf:"default:test"`
			Password string `conf:"default:test"`
		}
		Argon struct {
			Time    uint32 `conf:"default:4"`
			SaltLen uint32 `conf:"default:32"`
			Memory  uint32 `conf:"default:65536"` // 64*1024
			Threads uint8  `conf:"default:4"`
			KeyLen  uint32 `conf:"default:256"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "DimDim",
		},
	}

	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// -------------------------------------------------------------------------
	// App Starting

	log.Info(ctx, "starting service", "version", cfg.Build)
	defer log.Info(ctx, "shutdown complete")

	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("failed to generate config for output: %w", err)
	}

	log.Info(ctx, "startup", "config", out)

	// -------------------------------------------------------------------------
	// Database Support

	log.Info(ctx, "startup", "status", "initializing database support", "dsn", cfg.DB.DSN)

	db, err := storage.NewDB(cfg.DB.DSN, datastore.Migrations, datastore.Factory)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %w", err)
	}
	defer func() {
		log.Info(ctx, "shutdown", "status", "stopping database support")
		if err := db.Close(); err != nil {
			log.Error(ctx, "shutdown", "status", "failed to close the database", "error", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Password Encryption Support

	log.Info(ctx, "startup", "status", "initializing password encryption support", "config", cfg.Argon)

	argon := argon2id.New(cfg.Argon.Time, cfg.Argon.SaltLen, cfg.Argon.Memory, cfg.Argon.Threads, cfg.Argon.KeyLen)

	// -------------------------------------------------------------------------
	// SMTP Support

	log.Info(ctx, "startup", "status", "initializing smtp support", "config", cfg.Mailer)

	u, err := url.Parse(cfg.Mailer.Addr)
	if err != nil {
		return fmt.Errorf("failed to parse the mailer address: %w", err)
	}

	host := u.Hostname()
	if host == "" {
		host = u.Scheme
	}
	if host == "" {
		return fmt.Errorf("failed to parse the mailer host: %q", cfg.Mailer.Addr)
	}

	mailer := mailer.New(
		cfg.Mailer.Addr,
		host,
		cfg.Mailer.Identity,
		cfg.Mailer.Username,
		cfg.Mailer.Password,
	)

	// -------------------------------------------------------------------------
	// Start Debug Service

	go func() {
		log.Info(ctx, "startup", "status", "debug router started", "host", cfg.Debug.Host)

		expvar.NewString("build").Set(cfg.Build)
		if err := http.ListenAndServe(cfg.Debug.Host, debug.Mux()); err != nil {
			log.Error(ctx, "shutdown", "status", "debug router closed", "error", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Start Web Service

	log.Info(ctx, "startup", "status", "initializing web service")

	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	serverCfg := server.Config{
		AppName:     cfg.Web.AppName,
		Domain:      cfg.Web.DomainName,
		Port:        cfg.Web.Port,
		BindAddress: cfg.Web.BindAddress,

		ReadTimeout:     cfg.Web.ReadTimeout,
		WriteTimeout:    cfg.Web.WriteTimeout,
		IdleTimeout:     cfg.Web.IdleTimeout,
		ShutdownTimeout: cfg.Web.ShutdownTimeout,

		CORSAllowedOrigins: cfg.Web.CORSAllowedOrigins,
	}

	server := server.NewWebServer(serverCfg, argon, db, mailer)

	serverErrors := make(chan error, 1)

	go func() {
		log.Info(ctx, "startup", "status", "api router started", "host", serverCfg.FullDomain())

		serverErrors <- server.Start(serverCfg.Address())
	}()

	// -------------------------------------------------------------------------
	// Shutdown

	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		log.Info(ctx, "shutdown", "status", "shutdown started", "signal", sig)
		defer log.Info(ctx, "shutdown", "status", "shutdown complete", "signal", sig)

		ctx, cancel := context.WithTimeout(ctx, cfg.Web.ShutdownTimeout)
		defer cancel()

		if err := server.Shutdown(ctx); err != nil {
			server.Close()
			return fmt.Errorf("failed to stop the server gracefully: %w", err)
		}
	}

	return nil
}
