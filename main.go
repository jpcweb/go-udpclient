package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

func errorHandling(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	ip := flag.String("ip", "", "")
	port := flag.String("p", "", "")
	flag.Parse()

	udpRAddr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%s", *ip, *port))
	errorHandling(err)
	/*Connect to UDP server*/
	conn, err := net.DialUDP("udp4", nil, udpRAddr)
	errorHandling(err)
	defer conn.Close()
	conn.Write([]byte(""))

	/*outcoming channel > stdin*/
	outcoming := make(chan string)
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			errorHandling(err)
			outcoming <- text
		}
	}()
	/*incoming channel > read from conn*/
	incoming := make(chan string)
	go func() {
		buffer := make([]byte, 1024)
		for {
			n, err := conn.Read(buffer[0:])
			errorHandling(err)
			incoming <- string(buffer[0:n])
		}
	}()
	/*Infinite loop > execute according to channels */
	for {
		select {
		case out := <-outcoming:
			conn.Write([]byte(out))
		case in := <-incoming:
			fmt.Print(in)
		}
	}

}
