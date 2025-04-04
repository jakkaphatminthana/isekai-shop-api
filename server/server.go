package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/jakkaphatminthana/isekai-shop-api/config"
	"github.com/jakkaphatminthana/isekai-shop-api/databases"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"

	_adminRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/admin/repository"
	_oauth2Controller "github.com/jakkaphatminthana/isekai-shop-api/pkg/oauth2/controller"
	_oauth2Service "github.com/jakkaphatminthana/isekai-shop-api/pkg/oauth2/service"
	_playerRepository "github.com/jakkaphatminthana/isekai-shop-api/pkg/player/repository"
)

type echoServer struct {
	app  *echo.Echo
	db   databases.Database
	conf *config.Config
}

var (
	// Singleton pattern (call once)
	once sync.Once
	// Instance of echoServer
	server *echoServer
)

func NewEchoServer(conf *config.Config, db databases.Database) *echoServer {
	echoApp := echo.New()
	echoApp.Logger.SetLevel(log.DEBUG) //Highest level

	once.Do(func() {
		server = &echoServer{
			app:  echoApp,
			db:   db,
			conf: conf,
		}
	})

	return server
}

// Start Server
func (s *echoServer) Start() {
	//> middlerware setup
	corsMiddleware := getCORSMiddleware(s.conf.Server.AllowOrigins)
	bodyLimitMiddleware := getBodyLimitMiddleware(s.conf.Server.BodyLimit)
	timeOutMiddleware := getTimeOutMiddleware(s.conf.Server.TimeOut)

	//> middleware checklist
	s.app.Use(middleware.Recover()) //prevent server broken when occur error
	s.app.Use(middleware.Logger())
	s.app.Use(corsMiddleware)
	s.app.Use(bodyLimitMiddleware)
	s.app.Use(timeOutMiddleware)

	authorizingMiddleware := s.getAuthorizingMiddlerware()

	//> routers
	s.app.GET("/v1/health", s.healthCheck)
	s.initOAuth2Router()
	s.initItemShopRouter(authorizingMiddleware)
	s.initItemManagingRouter(authorizingMiddleware)
	s.initPlayerCoinRouter(authorizingMiddleware)
	s.initInventoryRouter(authorizingMiddleware)

	//> shutdown listener
	quitCh := make(chan os.Signal, 1)
	signal.Notify(quitCh, syscall.SIGINT, syscall.SIGALRM)
	go s.gracefullyShutdown(quitCh)

	//> start and listen
	s.httpListening()
}

// Start HTTP server, listen and serve
func (s *echoServer) httpListening() {
	serverUrl := fmt.Sprintf(":%d", s.conf.Server.Port)

	if err := s.app.Start(serverUrl); err != nil && err != http.ErrServerClosed {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

// Server shutdown
func (s *echoServer) gracefullyShutdown(quitCh chan os.Signal) {
	ctx := context.Background()

	<-quitCh
	s.app.Logger.Info("Shutting down server...")

	if err := s.app.Shutdown(ctx); err != nil {
		s.app.Logger.Fatalf("Error: %s", err.Error())
	}
}

// Server health check
func (s *echoServer) healthCheck(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

// Middleware
func getTimeOutMiddleware(timeout time.Duration) echo.MiddlewareFunc {
	return middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Skipper:      middleware.DefaultSkipper,
		ErrorMessage: "Request Timeout",
		Timeout:      timeout * time.Second,
	})
}

func getCORSMiddleware(allowOrigins []string) echo.MiddlewareFunc {
	return middleware.CORSWithConfig(middleware.CORSConfig{
		Skipper:      middleware.DefaultSkipper,
		AllowOrigins: allowOrigins,
		AllowMethods: []string{echo.GET, echo.POST, echo.PUT, echo.PATCH, echo.DELETE},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	})
}

func getBodyLimitMiddleware(bodyLimit string) echo.MiddlewareFunc {
	return middleware.BodyLimit(bodyLimit)
}

func (s *echoServer) getAuthorizingMiddlerware() *authorizingMiddleware {
	playerRepository := _playerRepository.NewPlayerRepositoryImpl(s.db, s.app.Logger)
	adminRepository := _adminRepository.NewAdminRepositoryImpl(s.db, s.app.Logger)

	oauth2Service := _oauth2Service.NewGoogleOAuth2Service(playerRepository, adminRepository)
	oauth2Controller := _oauth2Controller.NewGoogleOAuth2Controller(
		oauth2Service,
		s.conf.OAuth2,
		s.app.Logger,
	)

	return &authorizingMiddleware{
		oauth2Controller: oauth2Controller,
		oauth2Conf:       s.conf.OAuth2,
		logger:           s.app.Logger,
	}
}
