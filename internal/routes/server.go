package routes

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/joho/godotenv"
)

var Environment = "dev"

func init() {
    os.Setenv("env", Environment)
}

func Serve() {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

    err := godotenv.Load()
    if err != nil {
        logger.Error(fmt.Sprintf("Error loading .env file: %v", err))
    }

    router := newRouter()

    host := os.Getenv("HOST")
    port := os.Getenv("PORT")
    addr := fmt.Sprintf("%s:%s", host, port)


    killSig := make(chan os.Signal, 1)
    signal.Notify(killSig, os.Interrupt, syscall.SIGTERM)

    srv := &http.Server{
        Addr:    addr,
        Handler: router,
    }

    go func() {
        err := srv.ListenAndServe()

        if errors.Is(err, http.ErrServerClosed) {
            logger.Info("Server shutdown complete")
        } else if err != nil {
            logger.Error("Server error", slog.Any("err", err))
            os.Exit(1)
        }
    }()

    logger.Info("Server started", slog.String("host", host), slog.String("port", port), slog.String("env", Environment))
    <-killSig

    logger.Info("Shutting down server")


    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        logger.Error("Server shutdown failed", slog.Any("err", err))
        os.Exit(1)
    }

    logger.Info("Server shutdown complete")
}
