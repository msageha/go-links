package main

import (
	"context"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

const (
	DEV_ORIGIN           = "http://localhost:3000"
	SESSION              = "X-SESSION"
	SESSION_TOKEN_COOKIE = "token"
	API_TOKEN            = "X-API-TOKEN"
)

var excludeAuth = map[*regexp.Regexp][]string{
	regexp.MustCompile(`^/api$`):         {"GET"},
	regexp.MustCompile(`^/api/session$`): {"POST", "GET"},
}

type HttpServer struct {
	addr   string
	logger *zap.Logger
	server *http.Server
}

func NewHttpServer(addr string, logger *zap.Logger) *HttpServer {
	svr := &HttpServer{
		addr:   addr,
		logger: logger,
		server: &http.Server{
			Addr: addr,
		},
	}

	return svr
}

func (s *HttpServer) Start(ctx context.Context) error {
	s.server.Handler = s.handler(ctx)
	return s.server.ListenAndServe()
}

func (s *HttpServer) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *HttpServer) handler(ctx context.Context) http.Handler {
	router := mux.NewRouter()
	// Admin Session
	router.HandleFunc("/api", healthCheck).Methods("GET")
	return wrapHeader(router, s.logger)
}

func wrapHeader(next http.Handler, logger *zap.Logger) http.Handler {
	return wrapLoggerHeader(wrapAllowOrigin(wrapContentHeader(wrapAuth(next, logger))), logger)
}

func wrapContentHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func wrapAllowOrigin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin == DEV_ORIGIN {
			w.Header().Add("Access-Control-Allow-Origin", DEV_ORIGIN)
		}
		next.ServeHTTP(w, r)
	})
}

func isExcluded(method, path string) bool {
	for r, methods := range excludeAuth {
		if r.MatchString(path) {
			for _, m := range methods {
				if method == m {
					return true
				}
			}
		}
	}
	return false
}

func wrapAuth(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if isExcluded(r.Method, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		token := r.Header.Get(SESSION)
		if token == "" {
			var cookie *http.Cookie
			var err error
			cookie, err = r.Cookie(SESSION_TOKEN_COOKIE)
			if err == nil && cookie != nil {
				token = cookie.Value
			}
		}

		var ctx context.Context

		if token == "" {
			token = r.Header.Get(API_TOKEN)
			if token == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			//apiToken, err := ss.FindAPITokenByToken(r.Context(), token)
			//	if err != nil {
			//		logger.Sugar().Errorf("Failed to get api token: %v", err)
			//		w.WriteHeader(http.StatusInternalServerError)
			//		return
			//	}
			//
			//	if apiToken == nil {
			//		w.WriteHeader(http.StatusUnauthorized)
			//		return
			//	}
			//
			//	permission := &pm.AdminPermission{
			//		AdminEmail: apiToken.Email,
			//		Role:       pm.Admin,
			//	}
			//
			//	ctx = context.WithValue(r.Context(), ss.API_TOKEN, apiToken)
			//	ctx = context.WithValue(ctx, pm.PERMISSION, permission)
			//} else if strings.HasPrefix(r.URL.Path, ADMIN_PATH) {
			//	session, err := ss.FindOAuthAdminSessionByToken(r.Context(), token)
			//	if err != nil {
			//		logger.Sugar().Errorf("Failed to get admin session: %v", err)
			//		w.WriteHeader(http.StatusInternalServerError)
			//		return
			//	}
			//
			//	if session == nil {
			//		w.WriteHeader(http.StatusUnauthorized)
			//		return
			//	}
			//
			//	user, err := u.FindAdminUserByID(ctx, session.AdminUserID)
			//
			//	if err != nil {
			//		logger.Sugar().Errorf("Failed to get user: %v", err)
			//		w.WriteHeader(http.StatusInternalServerError)
			//		return
			//	}
			//
			//	if user == nil {
			//		w.WriteHeader(http.StatusUnauthorized)
			//		return
			//	}
			//
			//	permission := &pm.AdminPermission{
			//		AdminUserID: session.AdminUserID,
			//		AdminEmail:  user.Email,
			//		Role:        pm.Admin,
			//	}
			//
			//	ctx = context.WithValue(r.Context(), ss.SESSION, session)
			//	ctx = context.WithValue(ctx, pm.PERMISSION, permission)
			//} else {
			//	session, err := ss.FindSessionByToken(r.Context(), token)
			//	if err != nil {
			//		logger.Sugar().Errorf("Failed to get session: %v", err)
			//		w.WriteHeader(http.StatusInternalServerError)
			//		return
			//	}
			//
			//	if session == nil {
			//		w.WriteHeader(http.StatusUnauthorized)
			//		return
			//	}
			//
			//	permission := &pm.Permission{
			//		UserID: session.UserID,
			//		Role:   pm.User,
			//	}
			//
			//	ctx = context.WithValue(r.Context(), ss.SESSION, session)
			//	ctx = context.WithValue(ctx, pm.PERMISSION, permission)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func wrapLoggerHeader(next http.Handler, logger *zap.Logger) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lw := statusLoggingResponseWriter{
			ResponseWriter: w,
		}

		next.ServeHTTP(&lw, r)

		logger.Info("request",
			zap.Int("status", lw.status),
			zap.Duration("duration", time.Since(start)),
			zap.String("method", r.Method),
			zap.String("agent", r.UserAgent()),
			zap.String("url", r.URL.String()),
			zap.String("proto", r.Proto),
		)
	})
}

type statusLoggingResponseWriter struct {
	status int
	http.ResponseWriter
}

func (w *statusLoggingResponseWriter) WriteHeader(code int) {
	w.status = code
	w.ResponseWriter.WriteHeader(code)
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
