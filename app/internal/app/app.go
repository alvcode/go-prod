package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
	"golang.org/x/sync/errgroup"
	"net"
	"net/http"
	"prod/pkg/client/postgresql"
	"prod/pkg/metric"
	"time"

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

	//productStorage := storage.NewProductStorage(pgClient)

	return App{
		cfg:     cfg,
		router:  router,
		pgxPool: pgClient,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	grp, ctx2 := errgroup.WithContext(ctx)
	grp.Go(func() error {
		return a.startHTTP(ctx2)
	})

	logging.GetLogger(ctx).Info("Application initialized and started")
	return grp.Wait()
}

func (a *App) startHTTP(ctx context.Context) error {
	//logging.GetLogger(ctx).WithFields(map[string]interface{}{
	//	"IP":   a.cfg.HTTP.IP,
	//	"Port": a.cfg.HTTP.Port,
	//})

	logging.GetLogger(ctx).Printf("IP: %s, Port: %d", a.cfg.HTTP.IP, a.cfg.HTTP.Port)

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", a.cfg.HTTP.IP, a.cfg.HTTP.Port))
	if err != nil {
		logging.GetLogger(ctx).WithError(err).Fatal("failed to create http listener")
	}

	//logging.GetLogger(ctx).WithFields(map[string]interface{}{
	//	"AllowedMethods":     a.cfg.HTTP.CORS.AllowedMethods,
	//	"AllowedOrigins":     a.cfg.HTTP.CORS.AllowedOrigins,
	//	"AllowedHeaders":     a.cfg.HTTP.CORS.AllowedHeaders,
	//	"AllowCredentials":   a.cfg.HTTP.CORS.AllowCredentials,
	//	"OptionsPassthrough": a.cfg.HTTP.CORS.OptionsPassthrough,
	//	"ExposedHeaders":     a.cfg.HTTP.CORS.ExposedHeaders,
	//	"Debug":              a.cfg.HTTP.CORS.Debug,
	//})
	logging.GetLogger(ctx).Printf("CORS: %+v", a.cfg.HTTP.CORS)

	c := cors.New(cors.Options{
		AllowedMethods:     a.cfg.HTTP.CORS.AllowedMethods,
		AllowedOrigins:     a.cfg.HTTP.CORS.AllowedOrigins,
		AllowedHeaders:     a.cfg.HTTP.CORS.AllowedHeaders,
		AllowCredentials:   a.cfg.HTTP.CORS.AllowCredentials,
		OptionsPassthrough: a.cfg.HTTP.CORS.OptionsPassthrough,
		ExposedHeaders:     a.cfg.HTTP.CORS.ExposedHeaders,
		Debug:              a.cfg.HTTP.CORS.Debug,
	})

	handler := c.Handler(a.router)

	a.httpServer = &http.Server{
		Handler:      handler,
		WriteTimeout: a.cfg.HTTP.WriteTimeout,
		ReadTimeout:  a.cfg.HTTP.ReadTimeout,
	}

	logging.GetLogger(ctx).Println("http server started")
	if err = a.httpServer.Serve(listener); err != nil {
		switch {
		case errors.Is(err, http.ErrServerClosed):
			logging.GetLogger(ctx).Warningln("Server shutdown")
		default:
			logging.GetLogger(ctx).Fatalln(err)
		}
	}
	err = a.httpServer.Shutdown(context.Background())
	if err != nil {
		logging.GetLogger(ctx).Fatalln(err)
	}
	return err
}
