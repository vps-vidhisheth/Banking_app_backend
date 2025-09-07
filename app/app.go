package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	accountSvc "banking-app/component/account/service"
	authSvc "banking-app/component/auth/service"
	bankSvc "banking-app/component/banks/service"
	customerSvc "banking-app/component/customer/service"
	ledgerSvc "banking-app/component/ledger/service"
	transactionSvc "banking-app/component/transactions/service"
	"banking-app/model"

	accountHandler "banking-app/component/account/controller"
	authHandler "banking-app/component/auth/controller"
	bankHandler "banking-app/component/banks/controller"
	customerHandler "banking-app/component/customer/controller"
	ledgerHandler "banking-app/component/ledger/controller"
	transactionHandler "banking-app/component/transactions/controller"

	"banking-app/db"
	"banking-app/middleware"
	"banking-app/repository"

	"github.com/gorilla/mux"
)

type App struct {
	sync.Mutex
	Name   string
	Router *mux.Router
	Server *http.Server
	WG     *sync.WaitGroup

	Repository struct {
		Account     *repository.Repository[model.Account]
		Ledger      *repository.Repository[model.Ledger]
		Transaction *repository.Repository[model.Transaction]
		Customer    *repository.Repository[model.Customer]
		Bank        *repository.Repository[model.Bank]
	}

	Service struct {
		Account     *accountSvc.AccountService
		Customer    *customerSvc.CustomerService
		Ledger      *ledgerSvc.LedgerService
		Transaction *transactionSvc.TransactionService
		Auth        *authSvc.AuthService
		Bank        *bankSvc.BankService
	}

	Handler struct {
		Account     *accountHandler.AccountHandler
		Customer    *customerHandler.CustomerHandler
		Ledger      *ledgerHandler.LedgerHandler
		Transaction *transactionHandler.TransactionHandler
		Auth        *authHandler.AuthHandler
		Bank        *bankHandler.BankHandler
	}
}

// ---------------- Database Initialization ----------------
func (a *App) initDatabase() {
	db.InitDB()
	log.Println("Database initialized successfully")
}

// ---------------- Repositories Initialization ----------------
func (a *App) initRepositories() {
	database := db.GetDB()
	a.Repository.Account = repository.NewRepository[model.Account](database)
	a.Repository.Ledger = repository.NewRepository[model.Ledger](database)
	a.Repository.Transaction = repository.NewRepository[model.Transaction](database)
	a.Repository.Customer = repository.NewRepository[model.Customer](database)
	a.Repository.Bank = repository.NewRepository[model.Bank](database)
}

// ---------------- Handlers Initialization ----------------
func (a *App) initHandlers() {
	a.Handler.Customer = customerHandler.NewCustomerHandler(a.Service.Customer)
	a.Handler.Account = accountHandler.NewAccountHandler(a.Service.Account)
	a.Handler.Ledger = ledgerHandler.NewLedgerHandler(a.Service.Ledger)
	a.Handler.Transaction = transactionHandler.NewTransactionHandler(a.Service.Transaction)
	a.Handler.Auth = authHandler.NewAuthHandler(a.Service.Auth)
	a.Handler.Bank = bankHandler.NewBankHandler(a.Service.Bank)
}

// ---------------- Services Initialization ----------------
func (a *App) initServices(jwtSecret string) {
	dbConn := db.GetDB()

	// Initialize UnitOfWork for AccountService
	uow := repository.NewUnitOfWork(dbConn)

	// Ledger and Transaction services
	a.Service.Ledger = ledgerSvc.NewLedgerService(dbConn)
	a.Service.Transaction = transactionSvc.NewTransactionService(dbConn)

	// Customer service
	a.Service.Customer = customerSvc.NewCustomerService(dbConn)

	// Account service now uses UnitOfWork
	a.Service.Account = accountSvc.NewAccountService(uow, a.Service.Ledger, a.Service.Transaction)

	// Auth service with CustomerService injected (no jwtSecret needed here)
	a.Service.Auth = authSvc.NewAuthService(a.Service.Customer)

	// Bank service
	bankRepo := repository.NewRepository[model.Bank](dbConn)
	a.Service.Bank = bankSvc.NewBankService(bankRepo, dbConn)
}

// ---------------- Router Initialization ----------------
func (a *App) initRouter() {
	a.Router = mux.NewRouter().StrictSlash(true)

	// Global login route
	if a.Handler.Auth != nil {
		a.Router.HandleFunc("/login", a.Handler.Auth.LoginHandler).Methods("POST")
	}

	// API subrouter with authentication middleware
	api := a.Router.PathPrefix("/api/v1").Subrouter()
	api.Use(CORSMiddleware)
	api.Use(middleware.AuthMiddleware)

	// Register routes
	customerHandler.RegisterCustomerRoutes(api, a.Handler.Customer)
	accountHandler.RegisterAccountRoutes(api, a.Handler.Account)
	ledgerHandler.RegisterLedgerRoutes(api, a.Handler.Ledger)
	transactionHandler.RegisterTransactionRoutes(api, a.Handler.Transaction)
	bankHandler.RegisterBankRoutes(api, a.Handler.Bank)
}

// ---------------- Server Initialization ----------------
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
	log.Printf("Server initialized on port %s", port)
}

func (a *App) Start() error {
	log.Println("🔹 Starting server...")
	if err := a.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return err
	}
	return nil
}

func (a *App) Stop() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Server.Shutdown(ctx); err != nil {
		log.Printf("Server shutdown failed: %v", err)
		return
	}

	log.Println("Server stopped gracefully")
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
