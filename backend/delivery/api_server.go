package delivery

import (
	"gandalf-data-aggregator/auth"
	"gandalf-data-aggregator/config"
	token "gandalf-data-aggregator/pkg/jwt"
	"gandalf-data-aggregator/service"
	"net/http"
	"net/url"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type ServerInterface interface {
	UserActivity(ctx echo.Context) error
	CompleteAuth(ctx echo.Context) error
	BeginAuth(ctx echo.Context) error
	CurrentUser(ctx echo.Context) error
	RegisterUserDataKey(ctx echo.Context) error
	GenerateGandalfCallback(ctx echo.Context) error
	HandleTwitterCallback(ctx echo.Context) error
}

type Server struct {
	service  *service.Service
	router   *echo.Echo
	jwtMaker token.Maker
	cfg      config.Config
}

func NewServer(e *echo.Echo, cfg config.Config, service *service.Service, jwtMaker token.Maker) *Server {
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	server := &Server{
		cfg:      cfg,
		router:   e,
		jwtMaker: jwtMaker,
		service:  service,
	}

	server.registerUnAuthHandlers()
	server.registerAuthHandlers()
	return server
}

func (s Server) Run(address string) error {
	return s.router.Start(address)
}

func (s *Server) registerAuthHandlers() {
	authGroup := s.router.Group("/user")

	authGroup.Use(auth.JWTMiddleware(s.jwtMaker))

	authGroup.GET("/me", s.CurrentUser)
	authGroup.GET("/activity", s.UserActivity)
	authGroup.GET("/generate-callback", s.GenerateGandalfCallback)
}

func (s *Server) registerUnAuthHandlers() {
	s.router.POST("/auth/twitter/complete", s.CompleteAuth)
	// TODO: callback should be a frontend url
	s.router.GET("/auth/twitter/callback", s.HandleTwitterCallback)
	s.router.GET("/auth/twitter", s.BeginAuth)
	s.router.GET("/gandalf/callback/:source/:state", s.RegisterUserDataKey)
}

func (s *Server) CurrentUser(c echo.Context) error {
	userID, ok := c.Get("UserID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	user, err := s.service.GetUserByID(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, "Unable to get user")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"username":   user.Username,
		"avatar_url": user.AvatarURL,
	})
}

type AuthRequest struct {
	AuthURL string `json:"authURL"`
}

func (s *Server) CompleteAuth(c echo.Context) error {
	var req AuthRequest
	if err := c.Bind(&req); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid request format")
	}

	user, token, err := s.service.CompleteAuth(c.Request().Context(), req.AuthURL)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to complete authentication")
	}

	return c.JSON(http.StatusOK, map[string]string{
		"token":    token,
		"username": user.Username,
	})
}

func (s *Server) BeginAuth(c echo.Context) error {
	authURL, err := s.service.BeginAuth(c.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to begin authentication")
	}

	return c.Redirect(http.StatusFound, authURL)
}

func (s *Server) HandleTwitterCallback(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"url": c.Request().URL.String()})
}

func (s *Server) RegisterUserDataKey(c echo.Context) error {
	parsedURL, err := url.Parse(c.Request().URL.String())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to register data key")
	}

	_ = s.service.RegisterUserDataKey(
		c.Request().Context(),
		c.Param("state"),
		parsedURL.Query().Get("dataKey"),
		c.Param("source"),
	)

	return c.JSON(http.StatusOK, nil)
}

func (s *Server) GenerateGandalfCallback(c echo.Context) error {
	userID, ok := c.Get("UserID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	callbackURL, err := s.service.GenerateGandalfCallback(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Unable to generate callback url")
	}

	return c.JSON(http.StatusOK, map[string]string{"callbackURL": callbackURL})
}

type UserActivityParams struct {
	Limit int `query:"limit"`
	Page  int `query:"page"`
}

func (s *Server) UserActivity(c echo.Context) error {
	userID, ok := c.Get("UserID").(uuid.UUID)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "Invalid user")
	}

	var params UserActivityParams
	err := c.Bind(&params)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid parameter")
	}

	activities, err := s.service.GetActivitySetByUser(c.Request().Context(), userID, params.Limit, params.Page)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	stats, err := s.service.GenerateUserYearlyData(c.Request().Context(), userID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to fetch activities")
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"activities": activities,
		"stats":      stats,
	})
}
