package app

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"net"
	"net/http"
	"prod/internal/domain/product/storage"
	"prod/pkg/client/postgresql"
	"prod/pkg/metric"
	"time"

	"os"
	"path/filepath"
	_ "prod/docs"
	"prod/internal/config"
	"prod/pkg/logging"
)

type App struct {
	cfg        *config.Config
	router     *httprouter.Router
	httpServer *http.Server
	pgxPool    *pgxpool.Pool
}

func NewApp(ctx context.Context, cfg *config.Config) (App, error) {
	logging.GetLogger(ctx).Println("router init")
	router := httprouter.New()

	logging.GetLogger(ctx).Println("swagger init")
	router.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	router.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)

	metricHandler := metric.Handler{}
	metricHandler.Register(router)

	pgConfig := postgresql.NewPgConfig(cfg.PostgreSQL.Host, cfg.PostgreSQL.Port, cfg.PostgreSQL.Username, cfg.PostgreSQL.Password, cfg.PostgreSQL.Database)
	pgClient, err := postgresql.NewClient(ctx, 5, time.Second*5, pgConfig)
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}

	productStorage := storage.NewProductStorage(pgClient)
	all, err := productStorage.All(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}
	logging.GetLogger(ctx).Println(all)
	logging.GetLogger(ctx).Fatalln(all)

	return App{
		cfg:     cfg,
		router:  router,
		pgxPool: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) {
	a.startHTTP(ctx)
}

func (a *App) startHTTP(ctx context.Context) {
	logging.GetLogger(ctx).Infoln("start http")

	var listener net.Listener

	if a.cfg.Listen.Type == "sock" {
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
		}
		socketPath := filepath.Join(appDir, a.cfg.Listen.SocketFile)
		logging.GetLogger(ctx).Infof("socket path: %s", socketPath)

		listener, err = net.Listen("unix", socketPath)
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
		}
	} else {
		logging.GetLogger(ctx).Infof("bind app to host: %s and port: %s", a.cfg.Listen.BindIP, a.cfg.Listen.Port)
		var err error
		listener, err = net.Listen("tcp", a.cfg.Listen.BindIP+":"+a.cfg.Listen.Port)
		if err != nil {
			logging.GetLogger(ctx).Fatalln(err)
		}
	}

	c := cors.New(cors.Options{
		AllowedMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowedOrigins:     []string{"*"},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		ExposedHeaders:     []string{"*"},
		Debug:              false,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logging.GetLogger(ctx).Println("http server started")
	if err := a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warningln("Server shutdown")
		default:
			logging.GetLogger(ctx).Fatalln(err)
		}
	}
	err := a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}
}
