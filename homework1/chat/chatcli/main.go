package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	fmt.Println("Enter nickename:")
	readbuf := bufio.NewReader(os.Stdin)
	nickname, _, err := readbuf.ReadLine()
	if err != nil {
		fmt.Println("cannot read")
		return
	}

	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write(nickname)
	if err != nil {
		fmt.Printf("could not send nickname to server, %v", err)
		return
	}

	go func() {
		io.Copy(os.Stdout, conn)
	}()
	io.Copy(conn, os.Stdin) // until you send ^Z
	fmt.Printf("%s: exit", conn.LocalAddr())
}
