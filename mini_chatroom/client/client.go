package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Client struct {
	ch   chan string
	conn net.Conn
	name string
	done chan struct{}
}

func newClient(name string, server string) (*Client, error) {
	conn, err := net.Dial("tcp", server)

	if err != nil {
		return nil, err
	}
	cl := &Client{
		ch:   make(chan string, 100),
		name: name,
		conn: conn,
		done: make(chan struct{}),
	}

	return cl, nil
}

func (client *Client) writeToServer() {
	for {
		select {
		case msg := <-client.ch:
			_, err := client.conn.Write([]byte(msg + "\n"))
			if err != nil {
				fmt.Println("Error enviando al servidor", err)
				client.closeDone()
				return
			}
		case <-client.done:
			return
		}
	}
}

func (client *Client) readFromServer() error {
	reader := bufio.NewReader(client.conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Conexión fallida:", err)
			return err
		}
		fmt.Print(msg)
	}
}

func (client *Client) input() {
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		msg := scanner.Text()
		select {
		case client.ch <- msg:
		case <-client.done:
			return
		}
	}
}

func (client *Client) closeDone() {
	select {
	case <-client.done:
	default:
		close(client.done)
	}
}

func (client *Client) announceName() {
	client.conn.Write([]byte(client.name + "\n"))

}

func main() {

	name := flag.String("name", "anon", "username")
	serverAddr := flag.String("addr", "localhost:9000", "server adress, <host>:<port>")

	flag.Parse()
	client, err := newClient(*name, *serverAddr)

	if err != nil {
		log.Fatal(err)
	}

	defer client.conn.Close()

	client.announceName()

	go client.readFromServer()
	go client.writeToServer()

	client.input()
}
