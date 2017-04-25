package main

import (
	"fmt"
	"github.com/go-redis/redis"
	"os/exec"
	"encoding/json"
	"os"
	"github.com/urfave/cli"
	"strconv"
	"log"
)


var (
	host string
	port string
	db string
	password string
	channel string
)


func main() {

	app := cli.NewApp()
	app.Name = "signalred"
	app.Usage = "coming soon"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "host",
			Value:       "localhost",
			Usage:       "The redis host",
			EnvVar:      "REDIS_HOST",
			Destination: &host,
		},
		cli.StringFlag{
			Name:        "port",
			Value:       "6379",
			Usage:       "The redis host's listening port",
			EnvVar:      "REDIS_PORT",
			Destination: &port,
		},
		cli.StringFlag{
			Name:        "db",
			Value:       "0",
			Usage:       "The redis host's db",
			EnvVar:      "REDIS_DB",
			Destination: &db,
		},
		cli.StringFlag{
			Name:        "password",
			Value:       "",
			Usage:       "The redis host's password",
			EnvVar:      "REDIS_PASSWORD",
			Destination: &password,
		},
		cli.StringFlag{
			Name:        "channel",
			Value:       "signalred*",
			Usage:       "The redis host's publishing channel",
			EnvVar:      "REDIS_CHANNEL",
			Destination: &channel,
		},

	}

	app.Action = listenForCommands
	app.Run(os.Args)
}

func listenForCommands( c *cli.Context){

	dbNum, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("Failed to parse DB number: %s", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: password,
		DB:       dbNum,
	})

	pubsub := client.PSubscribe(channel)
	defer pubsub.Close()

	fmt.Fprintf(os.Stdout,"Connected! host[%s:%s] db[%d] channel[%s]\n\n", host, port, dbNum, channel)

	for {
		fmt.Fprintf(os.Stdout, "Waiting for a command %s", "\n")
		msg, _ := pubsub.ReceiveMessage()
		fmt.Fprintf(os.Stdout,"Received command: %s\n", msg)

		payload := &Payload{}
		err := json.Unmarshal([]byte(msg.Payload), payload)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed parsing JSON: %s\n", err)
			continue
		}

		fmt.Fprintf(os.Stdout,"Executing command: %s\n", payload)



		output, err := exec.Command( payload.Cmd, payload.Args...).Output();
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed executing command: %s", err)
			continue
		}
		fmt.Fprintf(os.Stdout, "\nSuccessfully executed %s -> \n%s\n\n", payload.Cmd, string(output) )
	}

}

type Payload struct {
	Cmd string `json:cmd`
	Args []string `json:args`
}