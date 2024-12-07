package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/fx/fxtest"
	"go.uber.org/zap"
)

var (
	serverURL string
)

func TestFxMain(t *testing.T) {

	assert := assert.New(t)

	testApp := fxtest.New(
		t,
		fx.Provide(
			newLogger,
			fx.Annotate(newHealthHandler, fx.As(new(http.Handler))),
			httptest.NewServer,
		),
		fx.Invoke(httpTestServerLifecycle),
		fx.WithLogger(func(log *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: log}
		}),
	)

	assert.NoError(testApp.Start(context.TODO()))
	assert.NoError(httpbin())
	assert.NoError(testApp.Stop(context.TODO()))
}

func httpTestServerLifecycle(lc fx.Lifecycle, logger *zap.Logger, s *httptest.Server) {
	serverURL = s.URL
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			logger.Info("shutdown http/1 test server")
			s.Close()
			return nil
		},
	})
}

func httpbin() error {

	r, err := http.NewRequest(http.MethodGet, serverURL+"/healthz", nil)
	if err != nil {
		return fmt.Errorf("http.NewRequest: %v", err)
	}

	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		return fmt.Errorf("failed to execute http request: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d", resp.StatusCode)
	}

	return nil
}
