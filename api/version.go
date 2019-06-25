package api

import (
	"net/http"

	"github.com/go-park-mail-ru/2018_2_DeadMolesStudio/middleware"
)

func newApiV1Subrouter(hs ...http.Handler) http.Handler {
	v1 := http.NewServeMux()
	for _, h := range hs {
		v1.Handle("/v1/", http.StripPrefix("/v1", h))
	}
	return v1
}

func (s *Server) initApiV1() http.Handler {
	v1 := http.NewServeMux()

	auth := http.NewServeMux()
	auth.HandleFunc("/log_in", middleware.SessionMiddleware(http.HandlerFunc(s.logIn), s.sm))
	auth.HandleFunc("/log_out", middleware.SessionMiddleware(http.HandlerFunc(s.logOut), s.sm))
	auth.HandleFunc("/sign_up", s.signUp)
	auth.HandleFunc("/get_profile", middleware.SessionMiddleware(http.HandlerFunc(s.getProfile), s.sm))
	v1.Handle("/auth/", http.StripPrefix("/auth", auth))

	money := http.NewServeMux()
	money.HandleFunc("/get_categories", s.GetCategories)
	money.HandleFunc("/add_category", s.AddCategory)
	money.HandleFunc("/add_category_pic", s.AddCategoryPic)
	money.HandleFunc("/add_money", s.AddMoney)
	money.HandleFunc("/add_money_pic", s.AddMoneyPhoto)
	money.HandleFunc("/get_history", s.GetHistory)
	v1.Handle("/money/", middleware.SessionMiddleware(http.StripPrefix("/money", money), s.sm))

	return newApiV1Subrouter(v1)
}
