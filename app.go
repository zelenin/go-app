package app

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"syscall"
)

var ErrHandlerNotFound = errors.New("handler not found")

func NewApp() *App {
	return &App{
		handlers: []*handlerWrapper{},
		shutdownHandler: func(err error) error {
			return nil
		},
	}
}

type App struct {
	defaultHandler  Handler
	handlers        []*handlerWrapper
	shutdownHandler func(err error) error
}

func (app *App) Run() error {
	ctx := context.WithValue(context.Background(), keyAppPath, os.Args[0])
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	errChan := make(chan error)

	go func() {
		args := os.Args[1:]

		for _, handler := range app.handlers {
			if handler.match(args) {
				errChan <- handler.handler(ctx, args)
				return
			}
		}

		if app.defaultHandler != nil {
			errChan <- app.defaultHandler(ctx, args)
			return
		}

		errChan <- ErrHandlerNotFound
	}()

	select {
	case err := <-errChan:
		return app.shutdownHandler(err)
	case <-ctx.Done():
		return app.shutdownHandler(ctx.Err())
	}
}

func (app *App) AddDefaultHandler(handler Handler) {
	app.defaultHandler = handler
}

func (app *App) AddHandler(matcher Matcher, handler Handler) {
	app.handlers = append(app.handlers, &handlerWrapper{
		match:   matcher,
		handler: handler,
	})
}

func (app *App) AddSubCommand(name string, handler Handler) {
	app.handlers = append(app.handlers, &handlerWrapper{
		match: func(args []string) bool {
			return len(args) > 0 && args[0] == name
		},
		handler: handler,
	})
}

func (app *App) AddShutdownHandler(handler func(err error) error) {
	app.shutdownHandler = handler
}

type Matcher func(args []string) bool

type Handler func(ctx context.Context, args []string) error

type handlerWrapper struct {
	match   Matcher
	handler Handler
}
