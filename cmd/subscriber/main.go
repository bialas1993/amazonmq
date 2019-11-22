package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/go-stomp/stomp"
)

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
		log.Fatalln(err.Error())
	}
	defer conn.Disconnect()
	log.Println(conn.Version().String())

	sub, err := conn.Subscribe(*queueName, stomp.AckAuto)
	defer conn.Disconnect()

	println("Reading from queue..")

	for {
		select {
		case msg := <-sub.C:
			println(string(msg.Body))
		}
	}
}
