# go-app

Simple multi-command application

## Usage
```go
package main

import (
	"context"
	goapp "github.com/zelenin/go-app"
	"log"
)

func main() {
	app := goapp.NewApp()

	app.AddSubCommand("command1", func(ctx context.Context, args []string) error {
		log.Printf("[%s] args: %v", args[0], args)
		return nil
	})

	app.AddSubCommand("command2", func(ctx context.Context, args []string) error {
		log.Printf("[%s] args: %v", args[0], args)
		return nil
	})

	err := app.Run()
	if err != nil {
		log.Fatalf("app.Run: %s", err)
	}
}
```

[![asciicast](https://asciinema.org/a/478476.svg)](https://asciinema.org/a/478476)

## Author

[Aleksandr Zelenin](https://github.com/zelenin/), e-mail: [aleksandr@zelenin.me](mailto:aleksandr@zelenin.me)
