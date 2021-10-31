package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/zzJinux/tcp-piercer/client"
	"github.com/zzJinux/tcp-piercer/server"
	"github.com/zzJinux/tcp-piercer/share"
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
	serverPort, err := share.PortValidate(*port)
	if err != nil {
		log.Fatal(err)
	}

	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	if err := s.StartContext(ctx, serverPort); err != nil {
		log.Fatal(err)
	}
}

func clientMain(args []string) {
	subFlag := flag.NewFlagSet("client", flag.ContinueOnError)
	subFlag.Parse(args)

	args = subFlag.Args()
	if len(args) < 2 {
		log.Fatalf("Specify a server and a host port to be serviced")
	}

	server := args[0]
	srvPort, err := share.PortValidate(args[1])
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	c, err := client.NewClient(server, srvPort)
	if err != nil {
		log.Fatalf("main: %v", err)
	}

	ctx := context.Background()
	if err := c.StartContext(ctx); err != nil {
		log.Fatal(err)
	}
}
