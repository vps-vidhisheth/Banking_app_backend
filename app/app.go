package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"banking-app/db"
	"banking-app/handler"
	"banking-app/middleware"
	"banking-app/repository"
	"banking-app/service"

	"github.com/gorilla/mux"
)

type App struct {
	sync.Mutex
	Name       string
	Router     *mux.Router
	Server     *http.Server
	WG         *sync.WaitGroup
	Repository struct {
		Account     *repository.AccountRepository
		Ledger      *repository.LedgerRepository
		Transaction *repository.TransactionRepository
		Customer    *repository.CustomerRepository
		Bank        *repository.BankRepository
	}
	Service struct {
		Account     *service.AccountService
		Customer    *service.CustomerService
		Ledger      *service.LedgerService
		Transaction *service.TransactionService
		Auth        *service.AuthService
		Bank        *service.BankService
	}
	Handler struct {
		Account     *handler.AccountHandler
		Customer    *handler.CustomerHandler
		Ledger      *handler.LedgerHandler
		Transaction *handler.TransactionHandler
		Auth        *handler.AuthHandler
		Bank        *handler.BankHandler
	}
}

func NewApp(name string, wg *sync.WaitGroup, jwtSecret string) *App {
	a := &App{
		Name: name,
		WG:   wg,
	}

	a.initDatabase()
	a.initRepositories()
	a.initServices(jwtSecret)
	a.initHandlers()
	a.initRouter()
	a.initServer()

	return a
}

func (a *App) initDatabase() {
	db.InitDB()
	log.Println("✅ Database initialized successfully")
}

func (a *App) initRepositories() {
	database := db.GetDB()
	a.Repository.Account = repository.NewAccountRepository(database)
	a.Repository.Ledger = repository.NewLedgerRepository(database)
	a.Repository.Transaction = repository.NewTransactionRepository(database)
	a.Repository.Customer = repository.NewCustomerRepository(database)
	a.Repository.Bank = repository.NewBankRepository(database)
}

func (a *App) initServices(jwtSecret string) {
	a.Service.Ledger = service.NewLedgerService(a.Repository.Ledger)
	a.Service.Transaction = service.NewTransactionService(a.Repository.Transaction)
	a.Service.Customer = service.NewCustomerService(a.Repository.Customer)
	a.Service.Account = service.NewAccountService(a.Repository.Account, a.Service.Ledger, a.Service.Transaction)
	a.Service.Auth = service.NewAuthService(a.Service.Customer)

	a.Service.Bank = service.NewBankService(a.Repository.Bank)
}

func (a *App) initHandlers() {
	a.Handler.Customer = handler.NewCustomerHandler(a.Service.Customer)
	a.Handler.Account = handler.NewAccountHandler(a.Service.Account)
	a.Handler.Ledger = handler.NewLedgerHandler(a.Service.Ledger)
	a.Handler.Transaction = handler.NewTransactionHandler(a.Service.Transaction)
	a.Handler.Auth = handler.NewAuthHandler(a.Service.Auth)
	a.Handler.Bank = handler.NewBankHandler(a.Service.Bank)
}

func (a *App) initRouter() {
	a.Router = mux.NewRouter().StrictSlash(true)

	a.Router.HandleFunc("/login", a.Handler.Auth.LoginHandler).Methods("POST")

	api := a.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(CORSMiddleware)
	api.Use(middleware.AuthMiddleware)

	handler.RegisterCustomerRoutes(api, a.Handler.Customer)
	handler.RegisterAccountRoutes(api, a.Handler.Account)
	handler.RegisterLedgerRoutes(api, a.Handler.Ledger)
	handler.RegisterTransactionRoutes(api, a.Handler.Transaction)
	handler.RegisterBankRoutes(api, a.Handler.Bank)
}

func (a *App) initServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	a.Server = &http.Server{
		Addr:         ":" + port,
		Handler:      a.Router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}
	log.Printf("✅ Server initialized on port %s", port)
}

func (a *App) Start() error {
	log.Println("🔹 Starting server at port 8080...")
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf(" Server shutdown failed: %v", err)
		return
	}

	log.Println(" Server stopped gracefully")
}

func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4200")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
