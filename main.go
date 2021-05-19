package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zzJinux/tcp-piercer/client"
	"github.com/zzJinux/tcp-piercer/server"
)

var help = `<main help message here>`

func main() {

	var subcmd string
	args := flag.Args()
	if len(args) > 0 {
		subcmd = args[0]
		args = args[1:]
	}

	switch subcmd {
	case "server":
		serverMain(args)
	case "client":
		clientMain(args)
	default:
		fmt.Print(help)
		os.Exit(1)
	}
}

func serverMain(args []string) {
	subFlag := flag.NewFlagSet("server", flag.ContinueOnError)
	p := subFlag.String("p", "", "")
	port := subFlag.String("port", "", "")
	subFlag.Parse(args)

	if *port == "" {
		*port = *p
	}
	if *port == "" {
		log.Fatalf("Specify a listening port")
	}

	s, err := server.NewServer(*port)
	if err != nil {
		// TODO: error handling
		log.Fatalln("NewServer fails")
	}

	if s.Start() != nil {
		// TODO: error handling
		log.Fatalln("server error!")
	}
}

func clientMain(args []string) {
	subFlag := flag.NewFlagSet("client", flag.ContinueOnError)
	subFlag.Parse(args)

	args = subFlag.Args()
	if len(args) < 2 {
		log.Fatalf("Specify a server and at least one host port")
	}

	serverAddr := args[0]
	port := args[1]

	c, err := client.NewClient(serverAddr, port)
	if err != nil {
		// TODO: error handling
		log.Fatalln("NewClient fails")
	}

	if c.Start() != nil {
		// TODO: error handling
		log.Fatal("client error!")
	}
}
