package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	ip := flag.String("ip", "127.0.0.1", "IP address of the UDP server")
	port := flag.String("p", "8000", "Port of the UDP server")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *ip, *port)
	conn, err := setupConnection(addr)
	if err != nil {
		log.Printf("Error setting up connection: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	outgoing := make(chan string)
	incoming := make(chan string)

	go readStdin(outgoing)
	go readConnection(conn, incoming)

	for {
		select {
		case out := <-outgoing:
			if _, err := send(conn, out); err != nil {
				log.Printf("Error sending data: %v\n", err)
				continue
			}
		case in := <-incoming:
			log.Print(in)
		}
	}
}

func setupConnection(addr string) (*net.UDPConn, error) {
	udpRAddr, err := net.ResolveUDPAddr("udp4", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialUDP("udp4", nil, udpRAddr)
	if err != nil {
		return nil, err
	}
	log.Printf("Listening on %v", addr)
	return conn, nil
}

func readStdin(out chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		text, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("Error reading from stdin: %v\n", err)
			close(out)
			return
		}
		out <- text
	}
}

func readConnection(conn *net.UDPConn, in chan<- string) {
	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Printf("Error reading from connection: %v\n", err)
			close(in)
			return
		}
		// Apply incoming middleware to the message before sending it to the channel
		modifiedMessage := string(buffer[:n])
		in <- modifiedMessage
	}
}

func send(conn *net.UDPConn, data string) (int, error) {
	return conn.Write([]byte(data))
}
