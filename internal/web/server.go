package web

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/garnizeH/dimdim/embeded"
	"github.com/garnizeH/dimdim/pkg/domain"
	"github.com/garnizeH/dimdim/service"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	session "github.com/spazzymoto/echo-scs-session"
)

type Config struct {
	AppName     string
	Domain      string
	Port        string
	BindAddress string
	SessionsDSN string

	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration

	CORSAllowedOrigins []string
}

func (c Config) Address() string {
	return strings.Join([]string{c.BindAddress, ":", c.Port}, "")
}

func (c Config) AppURL() string {
	scheme := "https://"
	if strings.HasPrefix(c.Domain, "localhost") {
		scheme = "http://"
	}

	return strings.Join([]string{scheme, c.Domain, ":", c.Port}, "")
}

func (c Config) FullDomain() string {
	if strings.HasPrefix(c.Domain, "localhost") {
		return c.Domain + ":" + c.Port
	} else {
		return c.Domain
	}
}

func (c Config) IsLocalhost() bool {
	return strings.HasPrefix(c.Domain, "localhost")
}

type Server struct {
	*echo.Echo
	sessionManager *SessionManager
}

func (s Server) Stop() error {
	return s.sessionManager.Close()
}

func NewWebServer(
	cfg Config,
	service *service.Service,
) Server {
	e := echo.New()

	e.Debug = cfg.IsLocalhost()
	e.HideBanner = !cfg.IsLocalhost()
	e.HidePort = !cfg.IsLocalhost()

	templates := embeded.Templates()
	e.Renderer = templates
	e.HTTPErrorHandler = errorHandler(templates)

	e.Use(middleware.BodyLimit("1k"))

	// Setup CSRF protection.
	e.Use(middleware.CSRFWithConfig(middleware.CSRFConfig{
		TokenLookup:    "form:csrf_token",
		CookieMaxAge:   int(1 * time.Hour / time.Second),
		CookieHTTPOnly: true,
		CookieSecure:   true,
		CookieName:     "_csc",
		ContextKey:     "csc",
		CookiePath:     "/",
		CookieSameSite: http.SameSiteStrictMode,
		Skipper: func(c echo.Context) bool {
			path := c.Request().URL.Path
			return strings.HasPrefix(path, "/static") || path == "/payment_hook"
		},
	}))

	// Setup CORS.
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
	}))

	// Setup session handling.
	sm, err := newSessionManager(cfg.SessionsDSN)
	if err != nil {
		panic(err)
	}

	sessionManager := sm.SessionManager()
	e.Use(session.LoadAndSave(sessionManager))
	e.Use(sessionDataMiddleware(sessionManager, service.User(), cfg.AppName))

	// Setup handler.
	domain := domain.Domain(cfg.FullDomain())
	handlers := NewHandler(domain, sessionManager, service)
	handlers.LoadRoutes(e, templates)

	// Setup static page serving.
	staticG := e.Group("static")
	staticG.Use(middleware.Gzip())
	staticG.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Set cache control header to 1 year so we can cache for a long time any static file.
			// This means if we need to update a static file, we need to change its name.
			//
			// WARNING: This will cache 4xx-5XX responses as well. We should instead write our own Static
			// handler that caches only on success.
			c.Response().Header().Set(
				"Cache-Control",
				"max-age="+strconv.Itoa(int(365*24*time.Hour/time.Second)),
			)
			return next(c)
		}
	})
	staticG.StaticFS("/", embeded.Static())

	return Server{
		Echo:           e,
		sessionManager: sm,
	}
}
