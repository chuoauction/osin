package osin

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

// Server is an OAuth2 implementation
type Server struct {
	Config            *ServerConfig
	Storage           Storage
	AuthorizeTokenGen AuthorizeTokenGen
	AccessTokenGen    AccessTokenGen
	Now               func() time.Time
	Logger            Logger
}

// NewServer creates a new server instance
func NewServer(config *ServerConfig, storage Storage) *Server {
	return &Server{
		Config:            config,
		Storage:           storage,
		AuthorizeTokenGen: &AuthorizeTokenGenDefault{},
		AccessTokenGen:    &AccessTokenGenDefault{},
		Now:               time.Now,
		Logger:            &LoggerDefault{},
	}
}

// NewResponse creates a new response for the server
func (s *Server) NewResponse() *Response {
	r := NewResponse(s.Storage)
	r.ErrorStatusCode = s.Config.ErrorStatusCode
	return r
}

func (s *Server) ValidatorMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bearer := CheckBearerAuth(c)
		if bearer == nil {
			s.logError(E_INVALID_REQUEST, nil, "handle_info_request=%s", "bearer is nil")
			return echo.NewHTTPError(http.StatusUnauthorized, "bearer is nil")
		}

		if bearer.Code == "" {
			s.logError(E_INVALID_REQUEST, nil, "handle_info_request=%s", "code is nil")
			return echo.NewHTTPError(http.StatusUnauthorized, "bearer code is invalid")
		}

		accessData, err := s.Storage.LoadAccess(bearer.Code)
		if err != nil {
			s.logError(E_INVALID_REQUEST, err, "handle_info_request=%s", "failed to load access data")
			return echo.NewHTTPError(http.StatusUnauthorized, "access token is invalid")
		}

		if accessData == nil {
			s.logError(E_INVALID_REQUEST, nil, "handle_info_request=%s", "access data is nil")
			return echo.NewHTTPError(http.StatusUnauthorized, "access data is invalid")
		}
		if accessData.Client == nil {
			s.logError(E_UNAUTHORIZED_CLIENT, nil, "handle_info_request=%s", "access data client is nil")
			return echo.NewHTTPError(http.StatusUnauthorized, "client data is invalid")
		}
		if accessData.Client.GetRedirectUri() == "" {
			s.logError(E_UNAUTHORIZED_CLIENT, nil, "handle_info_request=%s", "access data client redirect uri is empty")
			return echo.NewHTTPError(http.StatusUnauthorized, "client redirect uri is invalid")
		}
		if accessData.IsExpiredAt(s.Now()) {
			s.logError(E_INVALID_GRANT, nil, "handle_info_request=%s", "access data is expired")
			return echo.NewHTTPError(http.StatusUnauthorized, "access token is expired")
		}

		c.Set("AccessData", accessData)
		return next(c)
	}
}
