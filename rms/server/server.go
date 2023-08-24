package server

import (
	"context"
	"github.com/go-chi/chi/v5"
	"net/http"
	"rms/handler"
	"rms/middlewares"
	"rms/utils"
	"time"
)

type Server struct {
	chi.Router
	server *http.Server
}

const (
	readTimeout       = 5 * time.Minute
	readHeaderTimeout = 30 * time.Second
	writeTimeout      = 5 * time.Minute
)

func SetupRouter() *Server {
	router := chi.NewRouter()
	router.Route("/rms", func(r chi.Router) {
		r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
			utils.RespondJSON(w, http.StatusOK, map[string]interface{}{
				"status": "server is running",
			})
		})
		r.Route("/public", func(p chi.Router) {
			p.Post("/register", handler.Registration)
			p.Post("/login", handler.Login)
			p.Group(userRestaurant)

		})
		r.Route("/admin", func(a chi.Router) {
			a.Use(middlewares.AuthMiddleware)
			a.Use(middlewares.AuthAdmin)
			a.Group(adminUser)
			a.Group(adminSubAdmin)
			a.Group(adminRestaurant)
		})

		r.Route("/sub-admin", func(s chi.Router) {
			s.Use(middlewares.AuthMiddleware)
			s.Use(middlewares.AuthSubAdmin)
			s.Group(sub)
		})

		r.Route("/user", func(u chi.Router) {
			u.Use(middlewares.AuthMiddleware)
			u.Use(middlewares.AuthUser)
			u.Group(userAddress)
			u.Delete("/delete", handler.DeleteUser)
			u.Group(userRestaurant)
		})
	})
	return &Server{
		Router: router,
	}
}

func (svc *Server) Run(port string) error {
	svc.server = &http.Server{
		Addr:              port,
		Handler:           svc.Router,
		ReadTimeout:       readTimeout,
		ReadHeaderTimeout: readHeaderTimeout,
		WriteTimeout:      writeTimeout,
	}
	return svc.server.ListenAndServe()
}

func (svc *Server) Shutdown(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	return svc.server.Shutdown(ctx)
}
