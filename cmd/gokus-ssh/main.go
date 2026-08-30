package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"gokus/internal/room"
	"gokus/internal/tui"

	tea "charm.land/bubbletea/v2"
	"charm.land/ssh"
	"charm.land/wish/v2"
	"charm.land/wish/v2/activeterm"
	"charm.land/wish/v2/bubbletea"
)

const (
	address     = "localhost:23234"
	hostKeyPath = ".ssh/gokus_ed25519"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	registry := room.NewRegistry()
	defer registry.Close()
	server, err := wish.NewServer(
		wish.WithAddress(address),
		wish.WithHostKeyPath(hostKeyPath),
		wish.WithMiddleware(bubbletea.Middleware(newTeaHandler(registry)), activeterm.Middleware()),
	)
	if err != nil {
		return fmt.Errorf("create SSH server: %w", err)
	}
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	log.Println("Starting SSH server", "address", address)
	go func() {
		if err = server.ListenAndServe(); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
			log.Println("could not start server", "error", err)
			done <- nil
		}
	}()
	<-done
	log.Println("stopping SSH server")
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil && !errors.Is(err, ssh.ErrServerClosed) {
		return fmt.Errorf("could not stop server: %w", err)
	}
	return nil
}

func newTeaHandler(registry *room.Registry) bubbletea.Handler {
	return func(sess ssh.Session) (tea.Model, []tea.ProgramOption) {
		roomID := sess.User()
		client, err := registry.Join(sess.Context(), roomID)
		if err != nil {
			wish.Fatalln(sess, "could not join room:", err)
			return nil, nil
		}
		return tui.NewRemoteModel(client), nil
	}
}
