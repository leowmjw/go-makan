package main

import (
	"fmt"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	fmt.Println("GO Makan!")
	app := &cli.App{
		Name:  "go-makan",
		Usage: "fight the loneliness!",
		Action: func(c *cli.Context) error {
			fmt.Println("Hello friend!")
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:  "start",
				Usage: "start",
				Action: func(c *cli.Context) error {
					fmt.Println("TestReceiveFirstItemStartOrderWorkflow")
					// Create the client
					// start the api server
					// Call http ..
					// tear it down ..
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
