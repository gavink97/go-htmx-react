package routes

import (
	"log"
	"log/slog"
	"net/http"
	"os"

	h "github.com/gavink97/gavin-site/internal/handlers"
	"github.com/gavink97/gavin-site/internal/hash/passwordhash"
	m "github.com/gavink97/gavin-site/internal/middleware"
	database "github.com/gavink97/gavin-site/internal/store/db"
	"github.com/gavink97/gavin-site/internal/store/dbstore"
	"github.com/joho/godotenv"
	"github.com/justinas/alice"
)

func newRouter() http.Handler {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	databaseName := os.Getenv("DATABASE_NAME")
	sessionCookieName := os.Getenv("SESSION_COOKIE_NAME")

	mux := http.NewServeMux()

	db := database.MustOpen(databaseName)
	passwordhash := passwordhash.NewHPasswordHash()

	userStore := dbstore.NewUserStore(
		dbstore.NewUserStoreParams{
			DB:           db,
			PasswordHash: passwordhash,
		})

	sessionStore := dbstore.NewSessionStore(
		dbstore.NewSessionStoreParams{
			DB: db,
		})

	// changes in version 1.23
	publicFiles := http.FileServer(http.Dir("./public"))
	mux.Handle("/public/*", http.StripPrefix("/public/", publicFiles))

	publicAssets := http.FileServer(http.Dir("./assets"))
	mux.Handle("/assets/*", http.StripPrefix("/assets/", publicAssets))

	authMiddleware := m.NewAuthMiddleware(sessionStore, sessionCookieName)

	// make a function to create + rotate log files
	err = os.Mkdir("logs", os.ModePerm)
	if err != nil {
		log.Printf("failed to create logs directory: %v", err)
	}

	logFile, err := os.OpenFile("logs/server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	logger := slog.New(slog.NewJSONHandler(logFile, nil))
	slog.SetDefault(logger)

	lm := loggingMiddleware(logger)

	authChain := alice.New(
		lm,
		m.RemoveTrailingSlashMiddleware,
		perClientRateLimiter,
		m.TextHTMLMiddleware,
		m.CSPMiddleware,
		authMiddleware.AddUserToContext,
	)

	endpointChain := alice.New(
		lm,
		jsonMiddleware,
		perClientRateLimiter,
		authMiddleware.AddUserToContext,
	)

	mux.Handle("GET /about", authChain.Then(http.HandlerFunc(h.NewAboutHandler().ServeHTTP)))

	mux.Handle("GET /register", authChain.Then(http.HandlerFunc(h.NewGetRegisterHandler().ServeHTTP)))

	mux.Handle("POST /register", authChain.Then(http.HandlerFunc(h.NewPostRegisterHandler(h.PostRegisterHandlerParams{
		UserStore: userStore,
	}).ServeHTTP)))

	mux.Handle("GET /login", authChain.Then(http.HandlerFunc(h.NewGetLoginHandler().ServeHTTP)))

	mux.Handle("POST /login", authChain.Then(http.HandlerFunc(h.NewPostLoginHandler(h.PostLoginHandlerParams{
		UserStore:         userStore,
		SessionStore:      sessionStore,
		PasswordHash:      passwordhash,
		SessionCookieName: sessionCookieName,
	}).ServeHTTP)))

	mux.Handle("POST /logout", authChain.Then(http.HandlerFunc(h.NewPostLogoutHandler(h.PostLogoutHandlerParams{
		SessionCookieName: sessionCookieName,
	}).ServeHTTP)))

	mux.Handle("/", authChain.Then(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			notFound := h.NewNotFoundHandler()
			notFound.ServeHTTP(w, r)
			return
		}
		h.NewHomeHandler().ServeHTTP(w, r)
	})))

	// mux.HandleFunc("/", h.NewHomeHandler().ServeHTTP)
	mux.Handle("/api/data", endpointChain.Then(http.HandlerFunc(endpointHandler)))

	return mux
}
