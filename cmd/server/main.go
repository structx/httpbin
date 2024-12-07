package main

import (
	"context"
	"errors"
	"net/http"
	"time"

	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type healthHandler struct {
	logger *zap.SugaredLogger
}

// interface compliance
var _ route = (*healthHandler)(nil)
var _ http.Handler = (*healthHandler)(nil)

func newHealthHandler(logger *zap.Logger) *healthHandler {
	return &healthHandler{
		logger: logger.Sugar().Named("HealthHandler"),
	}
}

func (hh *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hh.logger.Debug("Healthz")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte("OK"))
	if err != nil {
		hh.logger.Errorf("w.Write: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func (hh *healthHandler) pattern() string {
	return "/healthz"
}

func newLogger() (*zap.Logger, error) {
	return zap.NewDevelopment()
}

func newHttpServeMux(r route) *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle(r.pattern(), r)
	return mux
}

type route interface {
	http.Handler
	pattern() string
}

func asRoute(f any) any {
	return fx.Annotate(
		f,
		fx.As(new(route)),
	)
}

func newHttpServer(handler http.Handler) *http.Server {
	s := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  time.Second * 15,
		WriteTimeout: time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	return s
}

func httpServeLifecycle(lc fx.Lifecycle, logger *zap.Logger, s *http.Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				logger.Info("start http/1 server")
				if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Fatal("http/1 listen and serve: %v", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			logger.Info("shutdown http/1 server")
			return s.Close()
		},
	})
}

func main() {
	fx.New(
		fx.Provide(
			newLogger,
			asRoute(newHealthHandler),
			fx.Annotate(newHttpServeMux, fx.As(new(http.Handler))),
			newHttpServer,
		),
		fx.Invoke(httpServeLifecycle),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	).Run()
}
