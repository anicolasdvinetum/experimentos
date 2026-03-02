package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
)

type Client struct {
	ch   chan string
	conn net.Conn
	name string
}

type Server struct {
	clients map[string]*Client
	mu      sync.Mutex
}

func (server *Server) serverStart(addr string) error {
	listener, err := net.Listen("tcp", addr)

	if err != nil {
		return err
	}

	defer listener.Close()

	fmt.Println("Servidor escuchando en:", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error aceptando conexion")
			continue
		}
		go server.newConnection(conn)

	}

}

func (server *Server) newConnection(conn net.Conn) {
	reader := bufio.NewReader(conn)
	name, err := reader.ReadString('\n')
	if err != nil {
		conn.Close()
		return
	}

	name = strings.TrimSpace(name)

	server.mu.Lock()
	original := name
	i := 1
	for {
		if _, exists := server.clients[name]; !exists {
			break
		}
		name = fmt.Sprintf("%s%d", original, i)
		i++
	}
	client := &Client{
		name: name,
		conn: conn,
		ch:   make(chan string, 10),
	}
	server.clients[name] = client
	server.mu.Unlock()

	if client.name != original {
		conn.Write([]byte("/name " + name))
	}
	fmt.Printf("%s se ha unido\n", name)

	go server.receive(client, reader)
	go server.send(client)

}

func (server *Server) removeClient(client *Client) {
	server.mu.Lock()
	delete(server.clients, client.name)
	close(client.ch)
	client.conn.Close()
	server.mu.Unlock()
}

func (server *Server) receive(client *Client, reader *bufio.Reader) {
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Cliente desconectado:", client.name)
			server.removeClient(client)
			return
		}
		msg = strings.TrimSpace(msg)
		fullMsg := fmt.Sprintf("%s: %s", client.name, msg)
		fmt.Println(fullMsg)
		server.broadcast(fullMsg)

	}

}

func (server *Server) send(client *Client) {
	for msg := range client.ch {
		_, err := client.conn.Write([]byte(msg + "\n"))
		if err != nil {
			fmt.Println("Error enviando a", client.name, err)
			server.removeClient(client)
			return
		}
	}
}

func (server *Server) broadcast(msg string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	//si el canal del cliente está lleno, nos lo saltamos sin más
	for _, client := range server.clients {
		select {
		case client.ch <- msg:
		default:
		}
	}

}

func main() {
	serverAddr := flag.String("addr", "localhost:9000", "server adress, <host>:<port>")
	flag.Parse()

	server := &Server{
		clients: make(map[string]*Client),
	}

	if err := server.serverStart(*serverAddr); err != nil {
		log.Fatal(err)
	}
}
