package router

import (
	"github.com/denisecase/go-hunt-sql/api/middleware"
)

func (s *Server) initializeRoutes() {

	// Home Route
	s.Router.HandleFunc("/", s.ShowHome).Methods("GET")

	// Auth Routes (see auth controller)
	s.Router.HandleFunc("/user/login", s.ShowLogin).Methods("GET")
	s.Router.HandleFunc("/user/logout", s.ShowLogout).Methods("GET")
	s.Router.HandleFunc("/user/register", s.ShowRegister).Methods("GET")

	s.Router.HandleFunc("/user/login", middleware.SetMiddlewareJSON(s.PostLogin)).Methods("POST")
	// s.Router.HandleFunc("/forgot-password", middleware.SetMiddlewareJSON(s.Login)).Methods("GET")

	//Users routes
	s.Router.HandleFunc("/user/register", middleware.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user", middleware.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/user", middleware.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/user/{id}", middleware.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/user/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/user/{id}", middleware.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Team routes
	s.Router.HandleFunc("/team", middleware.SetMiddlewareJSON(s.CreateTeam)).Methods("POST")
	s.Router.HandleFunc("/team", middleware.SetMiddlewareJSON(s.GetTeams)).Methods("GET")
	s.Router.HandleFunc("/team/{id}", middleware.SetMiddlewareJSON(s.GetTeam)).Methods("GET")
	s.Router.HandleFunc("/team/{id}", middleware.SetMiddlewareJSON(middleware.SetMiddlewareAuthentication(s.UpdateTeam))).Methods("PUT")
	s.Router.HandleFunc("/team/{id}", middleware.SetMiddlewareAuthentication(s.DeleteTeam)).Methods("DELETE")
}
