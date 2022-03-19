package app

import (
	"context"
	"errors"
	"os"
)

var ErrHandlerNotFound = errors.New("handler not found")

func NewApp() *App {
	return &App{
		handlers: []*handlerWrapper{},
	}
}

type App struct {
	handlers []*handlerWrapper
}

func (app *App) Run() error {
	ctx := context.Background()
	ctx = context.WithValue(ctx, keyAppPath, os.Args[0])

	args := os.Args[1:]

	for _, handler := range app.handlers {
		if handler.checker(args) {
			return handler.handler(ctx, args)
		}
	}

	return ErrHandlerNotFound
}

func (app *App) AddHandler(checker Checker, handler Handler) {
	app.handlers = append(app.handlers, &handlerWrapper{
		checker: checker,
		handler: handler,
	})
}

type Checker func(args []string) bool

type Handler func(ctx context.Context, args []string) error

type handlerWrapper struct {
	checker Checker
	handler Handler
}

func CommandChecker(command string) Checker {
	return func(args []string) bool {
		return len(args) > 0 && args[0] == command
	}
}
