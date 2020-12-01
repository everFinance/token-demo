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
			&cli.StringFlag{Name: "token_symbol", Value: "ROL2", Usage: "token symbol", EnvVars: []string{"TOKEN_SYMBOL"}},
			&cli.StringFlag{Name: "token_owner", Value: "dQzTM9hXV5MD1fRniOKI3MvPF_-8b2XDLmpfcMN9hi8", Usage: "token owner", EnvVars: []string{"TOKEN_OWNER"}},
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
