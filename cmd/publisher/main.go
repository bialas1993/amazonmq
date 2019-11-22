package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"time"

	"github.com/go-stomp/stomp"
)

const defaultPort = ":61613"

var (
	serverAddr                             = flag.String("server", "127.0.0.1:61613", "STOMP server endpoint")
	messageCount                           = flag.Int("count", 10, "Number of messages to send/receive")
	queueName                              = flag.String("queue", "premium.paywall.subscriber.payment", "Destination queue")
	helpFlag                               = flag.Bool("help", false, "Print help text")
	stop                                   = make(chan bool)
	options      []func(*stomp.Conn) error = []func(*stomp.Conn) error{}
)

func main() {
	flag.Parse()
	if *helpFlag {
		fmt.Fprintf(os.Stderr, "Usage of %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(1)
	}

	netConn, err := tls.Dial("tcp", *serverAddr, &tls.Config{InsecureSkipVerify: true})
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer netConn.Close()

	conn, err := stomp.Connect(netConn)
	if err != nil {
		println("cannot connect to server", err.Error())
		return
	}
	defer conn.Disconnect()

	timer := time.NewTicker(time.Second * 3)

	println("Sending messages: ")

	for {
		select {
		case t := <-timer.C:
			print(".")
			if err = conn.Send(*queueName, "text/plain", []byte(t.String()), nil); err != nil {
				println("failed to send to server", err)
				return
			}
		}
	}
}
