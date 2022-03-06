package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

type client struct {
	ch      chan<- string
	nickCli string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string)
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

func broadcaster() {
	clients := make(map[client]bool)
	for {
		select {
		case msg := <-messages:
			for cli := range clients {
				cli.ch <- msg
			}

		case cli := <-entering:
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.ch)
		}
	}
}

func handleConn(conn net.Conn) {
	ch := make(chan string)
	go clientWriter(conn, ch)
	//realization client nickname
	who := conn.RemoteAddr().String()
	buffer := make([]byte, 200)
	_, err := conn.Read(buffer)
	if err != nil {
		fmt.Printf("could not read nickname from %s", who)
		return
	}
	nName := string(buffer)

	ch <- "You are " + nName
	messages <- nName + " has arrived"
	entering <- client{ch, nName}

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- nName + ": " + input.Text()
	}
	leaving <- client{ch, nName}
	messages <- nName + " has left"
	conn.Close()
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg)
	}
}
