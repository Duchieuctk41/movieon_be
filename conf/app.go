package conf

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-contrib/cors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/caarlos0/env/v6"
	"github.com/gin-gonic/gin"
	"github.com/praslar/cloud0/ginext"
	"github.com/praslar/cloud0/logger"
)

type BaseApp struct {
	Config     *AppConfig
	Name       string
	Version    string
	Router     *gin.Engine
	HttpServer *http.Server

	listener       net.Listener
	initialized    bool
	healthDisabled bool
}

func NewApp(name, version string) *BaseApp {
	app := &BaseApp{
		Name:           name,
		Version:        version,
		Router:         gin.New(),
		Config:         NewAppConfig(),
		HttpServer:     &http.Server{},
		healthDisabled: false,
	}

	app.HttpServer.Handler = app.Router

	return app
}

func (app *BaseApp) DisableHealthEndpoint() {
	app.healthDisabled = true
}

func (app *BaseApp) Initialize() error {
	if err := env.Parse(app); err != nil {
		return err
	}

	app.HttpServer.ReadTimeout = time.Duration(app.Config.ReadTimeout) * time.Second

	// wait to gin 1.8 to support this
	//_ = app.Router.SetTrustedProxies(app.Config.TrustedProxy)
	app.Router.ForwardedByClientIP = true

	// register default middlewares
	app.Router.Use(
		ginext.RequestIDMiddleware,
		ginext.AccessLogMiddleware(app.Config.Env),
		ginext.CreateErrorHandler(app.Config.Debug),
		cors.New(cors.Config{
			AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3003"},
			AllowMethods:     []string{"PUT", "PATCH", "GET", "POST", "DELETE"},
			AllowHeaders:     []string{"Origin, Access-Control-Allow-Headers, Content-Type, Authorization, X-Requested-With"},
			ExposeHeaders:    []string{"Content-Length"},
			AllowCredentials: true,
			AllowOriginFunc: func(origin string) bool {
				return origin == "http://localhost:3000"
			},
			MaxAge: 12 * time.Hour,
		}),
	)

	// register routes
	if !app.healthDisabled {
		healthHandler := app.HealthHandler()
		app.Router.GET("/status", healthHandler)
		app.Router.GET("/status-q", healthHandler)
	}

	app.Router.NoRoute(ginext.NotFoundHandler)

	if app.Config.EnableDB {
		err := OpenDefault(app.Config.DB)
		if err != nil {
			return errors.New("failed to open default DB: " + err.Error())
		}
	}

	app.initialized = true

	return nil
}

// HealthHandler makes health check handler
func (app *BaseApp) HealthHandler() gin.HandlerFunc {
	rsp := struct {
		Name     string `json:"name"`
		Version  string `json:"version"`
		Hostname string `json:"hostname"`
	}{
		Name:    app.Name,
		Version: app.Version,
	}
	rsp.Hostname, _ = os.Hostname()

	return func(c *gin.Context) {
		c.JSON(http.StatusOK, rsp)
	}
}

func (app *BaseApp) StartTLS(ctx context.Context, certPath string, keyPath string) error {
	l := logger.Tag("BaseApp.Start")
	var err error

	if !app.initialized {
		if err = app.Initialize(); err != nil {
			return errors.New("failed to initialize app: " + err.Error())
		}
	}

	if app.listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%d", app.Config.Port)); err != nil {
		return errors.New("failed to listen: " + err.Error())
	}

	errCh := make(chan error, 1)

	go func() {
		l.Printf("start listening on %s", app.listener.Addr().String())
		if err := app.HttpServer.ServeTLS(app.listener, certPath, keyPath); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}

		// no error, close channel
		close(errCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		defer func() {
			l.Info("shutting down http server ...")
			shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_ = app.HttpServer.Shutdown(shutCtx)
			cancel()
		}()
		select {
		case gotSignal, ok := <-signalCh:
			if !ok {
				// channel close
				return
			}
			l.Printf("got signal: %v", gotSignal)
			return
		case <-ctx.Done():
			l.Printf("context has done")
			return
		}
	}()

	go func() {
		l.Printf("start listening debug server on port %d", app.Config.DebugPort)
		_ = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(app.Config.DebugPort), nil)
	}()

	return <-errCh
}

func (app *BaseApp) Start(ctx context.Context) error {
	l := logger.Tag("BaseApp.Start")
	var err error

	if !app.initialized {
		if err = app.Initialize(); err != nil {
			return errors.New("failed to initialize app: " + err.Error())
		}
	}

	if app.listener, err = net.Listen("tcp4", fmt.Sprintf("0.0.0.0:%s", app.Config.Port)); err != nil {
		return errors.New("failed to listen: " + err.Error())
	}

	errCh := make(chan error, 1)

	go func() {
		l.Printf("start listening on %s", app.listener.Addr().String())
		if err := app.HttpServer.Serve(app.listener); err != nil && err != http.ErrServerClosed {
			errCh <- err
			return
		}

		// no error, close channel
		close(errCh)
	}()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP)
	go func() {
		defer func() {
			l.Info("shutting down http server ...")
			shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			_ = app.HttpServer.Shutdown(shutCtx)
			cancel()
		}()

		select {
		case gotSignal, ok := <-signalCh:
			if !ok {
				// channel close
				return
			}
			l.Printf("got signal: %v", gotSignal)
			return
		case <-ctx.Done():
			l.Printf("context has done")
			return
		}
	}()

	go func() {
		l.Printf("start listening debug server on port %d", app.Config.DebugPort)
		_ = http.ListenAndServe("0.0.0.0:"+strconv.Itoa(app.Config.DebugPort), nil)
	}()

	return <-errCh
}

func (app *BaseApp) Listener() net.Listener {
	return app.listener
}

func NewAppConfig() *AppConfig {
	return &AppConfig{
		DB: &Config{},
	}
}
