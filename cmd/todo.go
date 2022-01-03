package main

import (
	"context"
	"github.com/adrianpk/godddtodo/internal/base"
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
	a   *app.App
	log base.Logger
)

func main() {
	log = base.NewLogger("debug", true)

	// App
	a := app.NewApp(name, version, log)
	cfg := a.LoadConfig()

	// Context
	ctx, cancel := context.WithCancel(context.Background())
	initExitMonitor(ctx, cancel)

	// Database
	mgo := db.NewMongoClient("mongo-client", db.Config{
		Host:       cfg.Mongo.Host,
		Port:       cfg.Mongo.Port,
		User:       cfg.Mongo.User,
		Pass:       cfg.Mongo.Pass,
		Database:   cfg.Mongo.Database,
		MaxRetries: cfg.Mongo.MaxRetriesUInt64(),
	}, log)

	// Repo
	lrr := repo.NewListRead("list-read-repo", mgo, repo.Config{}, log)

	lwr := repo.NewListWrite("list-write-repo", mgo, repo.Config{}, log)

	// Service
	ts, err := service.NewTodo("todo-app-service", lrr, lwr, service.Config{}, log)

	if err != nil {
		exit(err)
	}

	a.TodoService = &ts

	// Init & Start
	err = a.InitAndStart()
	if err != nil {
		exit(err)
	}

	log.Errorf("%s stopped: %s (%s)", a.Name(), a.Version(), err)
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
