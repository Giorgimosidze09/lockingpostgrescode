package handlers

import (
	"database/sql"
	admin "lockingpostgrescode/handlers/admin"
	auth "lockingpostgrescode/handlers/authHandler"
	client "lockingpostgrescode/handlers/client"
	support "lockingpostgrescode/handlers/support"
	"lockingpostgrescode/middleware"
	"net/http"

	"github.com/gorilla/mux"
)

func SetupRoutes(db *sql.DB) *mux.Router {
	r := mux.NewRouter()

	//auth
	r.HandleFunc("/api/v1/register", auth.RegisterHandler).Methods("POST", "OPTIONS")

	r.HandleFunc("/api/v1/login", auth.LoginHandler).Methods("POST", "OPTIONS")

	//client
	r.Handle("/api/v1/withdraw", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client.WithdrawHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/deposit", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client.DepositHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/exchange", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		client.ExchangeHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")

	//admin
	r.Handle("/api/v1/adminwithdraw", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin.AdminWithdrawHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/admindeposit", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin.AdminDepositHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")
	r.Handle("/api/v1/adminexchange", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin.AdminExchangeHandler(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/history/{user_id}", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin.HistoryHandler(w, r, db)
	}))).Methods("GET", "OPTIONS")

	r.Handle("/api/v1/adminregister", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		admin.AdminRegisterHandler(w, r)
	}))).Methods("POST", "OPTIONS")

	//support
	r.Handle("/api/v1/deposit/handle", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		support.AcceptRejectDepositRequest(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/exchange/handle", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		support.AcceptRejectExchangeRequest(w, r, db)
	}))).Methods("POST", "OPTIONS")

	r.Handle("/api/v1/withdraw/handle", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		support.AcceptRejectWithdrawRequest(w, r, db)
	}))).Methods("POST", "OPTIONS")

	return r
}
