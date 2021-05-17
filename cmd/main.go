package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	handler2 "github.com/card-deck/internal/api/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/card-deck/internal/config"
	"github.com/card-deck/internal/repository"
	"go.uber.org/zap"
)

var (
	ctx        context.Context
	cancelFunc context.CancelFunc
)

func main() {
	var cfgFile string

	flag.StringVar(&cfgFile, "config", "config.hcl", "The config file")
	flag.Parse()

	cfg, err := config.NewConfig(cfgFile)
	if err != nil {
		log.Fatal(err)
	}
	logger, err := cfg.App.InitLogger()
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancelFunc = initAppContext(logger)
	defer cancelFunc()

	db, err := cfg.NewDBConnection(ctx)
	if err != nil {
		logger.Fatal("failed init DB connection", zap.String("error", err.Error()))
	}
	repo := repository.NewRepository(db, logger)

	router := chi.NewRouter()
	addMiddlewares(router)

	// init handlers
	cardGameHandler := handler2.NewCardGameHandler(repo, logger)

	// mount routes to handlers
	cardGameHandler.MountRoutes(router)

	if err = listenAndServe(ctx, cfg.App.Listening, logger.Sugar(), router); err != nil {
		logger.Fatal("failed start server", zap.String("error", err.Error()))
	}
}

// initAppContext inits app context for graceful shutdown
// os.Interrupt: Ctrl-C
// syscall.SIGTERM: kill PID, docker stop
func initAppContext(log *zap.Logger) (context.Context, context.CancelFunc) {
	ctx, cancelFunc = context.WithCancel(context.Background())
	go func() {
		signals := []os.Signal{os.Interrupt, syscall.SIGTERM}
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, signals...)
		defer signal.Reset(signals...)
		gotSignal := <-sigChan
		log.Info("stopping app", zap.String("signal", gotSignal.String()))
		cancelFunc()
	}()

	return ctx, cancelFunc
}

// addMiddlewares adds default middlewares to router
func addMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}

// listenAndServe starts server and handle correct shutdown
func listenAndServe(ctx context.Context, port int, logger *zap.SugaredLogger, router *chi.Mux) (err error) {
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		logger.Infof("server starting on port: %v", port)
		if err = srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen and serve failed: %v", err)
		}
	}()

	<-ctx.Done()

	logger.Info("stopping server")

	ctxShutDown, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err = srv.Shutdown(ctxShutDown); err != nil {
		logger.Fatalf("server shutdown failed: %v", err)
	}

	logger.Info("server stopped correctly")

	if err == http.ErrServerClosed {
		err = nil
	}

	return
}
