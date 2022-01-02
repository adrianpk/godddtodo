package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/adrianpk/godddtodo/internal/app"
	repo "github.com/adrianpk/godddtodo/internal/app/adapter/driver/repo/mongo"
	"github.com/adrianpk/godddtodo/internal/app/service"
	db "github.com/adrianpk/godddtodo/internal/base/db/mongo"
)

type contextKey string

const (
	name    = "todo"
	version = "v0.0.1"
)

var (
	a *app.App
)

func main() {
	// App
	a := app.NewApp(name, version)
	cfg := a.LoadConfig()

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	initExitMonitor(ctx, cancel)

	// Database
	mgo := db.NewMongoClient("mongo-client", db.Config{
		Host:         cfg.Mongo.Host,
		Port:         cfg.Mongo.Port,
		User:         cfg.Mongo.User,
		Pass:         cfg.Mongo.Pass,
		Database:     cfg.Mongo.Database,
		MaxRetries:   cfg.Mongo.MaxRetriesUInt64(),
		TracingLevel: cfg.Tracing.Level,
	})

	// Repo
	lrr := repo.NewListRead("list-read-repo", mgo, repo.Config{
		TracingLevel: cfg.Tracing.Level,
	})

	lwr := repo.NewListWrite("list-write-repo", mgo, repo.Config{
		TracingLevel: cfg.Tracing.Level,
	})

	// Service
	ts, err := service.NewTodo("todo-app-service", lrr, lwr, service.Config{
		TracingLevel: cfg.Tracing.Level,
	})

	if err != nil {
		exit(err)
	}

	a.TodoService = &ts

	// Init & Start
	err = a.InitAndStart()
	if err != nil {
		exit(err)
	}

	log.Fatalf("%s stoped: %s (%s)", a.Name(), a.Version(), err)
}

func exit(err error) {
	log.Fatal(err)
}

func initExitMonitor(ctx context.Context, cancel context.CancelFunc) {
	go checkSigterm(cancel)
	go checkCancel(ctx)
}

func checkSigterm(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	cancel()
}

func checkCancel(ctx context.Context) {
	<-ctx.Done()
	a.Stop()
	os.Exit(1)
}
