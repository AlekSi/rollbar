package main

import (
	"errors"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/AlekSi/rollbar"
)

func main() {
	token := flag.String("token", "", "API token (env: ROLLBAR_TOKEN)")
	flag.Parse()

	if *token == "" {
		*token = os.Getenv("ROLLBAR_TOKEN")
	}
	message := strings.Join(flag.Args(), " ")

	client := &rollbar.Client{Token: *token}
	err := client.Post(&rollbar.Payload{Error: errors.New(message)})
	if err != nil {
		log.Fatal(err)
	}
}
