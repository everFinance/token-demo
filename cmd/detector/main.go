package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/everFinance/token-demo/detector"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "detector",
		Flags: []cli.Flag{
			&cli.StringFlag{Name: "token_symbol", Value: "SCP02", Usage: "token symbol", EnvVars: []string{"TOKEN_SYMBOL"}},
			&cli.StringFlag{Name: "token_owner", Value: "z1Jhn1rXBXWUvXbXhQaWOFMP3Swdq6STA36IPdQKo50", Usage: "token owner", EnvVars: []string{"TOKEN_OWNER"}},
			&cli.StringFlag{Name: "port", Value: ":80", EnvVars: []string{"PORT"}},
		},
		Action: run,
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)

	d := detector.New(c.String("token_symbol"), c.String("token_owner"))
	d.Run(c.String("port"))

	<-signals
	return nil
}
