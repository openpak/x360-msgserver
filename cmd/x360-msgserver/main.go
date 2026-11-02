// Command x360-msgserver serves the Xbox Live message server, Live service 2 (see internal/msgserver).
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/openpak/x360-common/wire"
	"github.com/openpak/x360-msgserver/internal/msgserver"
)

func main() {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	addr := os.Getenv("X360_MSGSERVER_ADDR")
	if addr == "" {
		addr = ":10002"
	}

	titles, err := msgserver.LoadTitles(os.Getenv("X360_MSGSERVER_TITLES"))
	if err != nil {
		log.Error("titles", "err", err)
		os.Exit(1)
	}

	srv := &http.Server{Addr: addr, Handler: wire.Guard(os.Getenv("X360_GATEWAY_KEY"), &msgserver.Handler{Log: log, Titles: titles}), ReadHeaderTimeout: 10 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		<-ctx.Done()
		sc, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		srv.Shutdown(sc)
	}()
	log.Info("x360-msgserver listening", "addr", addr)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("serve", "err", err)
		os.Exit(1)
	}
}
