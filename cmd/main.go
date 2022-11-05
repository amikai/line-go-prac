package main

import (
	"log"

	"github.com/amikai/line-go-prac/cmd/bot"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "app [module]",
		Short: "Line bot system entrypoints",
	}

	cmd.AddCommand(bot.NewBotCommand())

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
