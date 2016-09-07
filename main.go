package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jessevdk/go-flags"
)

type options struct {
	ServerPort   int    `short:"s" long:"server" description:"Port to run server on!"`
	ClientString string `short:"c" long:"client" description:"Where to connect to the server!"`
}

func main() {
	var opts options
	_, err := flags.Parse(&opts)
	check(err)
	// check for invalid program state
	if opts.ServerPort > 0 && len(opts.ClientString) > 0 {
		panic("Cannot be both client and server")
	}
	if opts.ServerPort > 0 {
		runServer(opts.ServerPort)
	} else if len(opts.ClientString) > 0 {
		runClient(opts.ClientString)
	} else {
		panic("No client or server flag specified or invalid server string")
	}
}

func runClient(serverIP string) {
	// do client stuff
	log.Println("I am a client! Connecting to:", serverIP)
	conn, err := net.Dial("tcp", serverIP)
	check(err)
	go recieveMessages(conn)
	linesOfText := make(chan string, 10)
	go recieveText(linesOfText)
	for {
		select {
		case line := <-linesOfText:
			fmt.Fprintln(conn, line)
		}
	}
}

func runServer(port int) {
	// do server stuff
	log.Println("I am a server!")
	ln, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	check(err)
	newConns := make(chan net.Conn)
	clientMsgs := make(chan clientMsg, 100)
	go sendServerMsgs(newConns, clientMsgs)
	for {
		conn, err := ln.Accept()
		check(err)
		newConns <- conn
		go recieveServerMessages(conn, clientMsgs)
	}
}

func sendServerMsgs(newConns chan net.Conn, clientMsgs chan clientMsg) {
	var conns []net.Conn
	linesOfText := make(chan string, 10)
	go recieveText(linesOfText)
	for {
		select {
		case line := <-linesOfText:
			for _, conn := range conns {
				fmt.Fprintln(conn, line)
			}
		case clientMsg := <-clientMsgs:
			for _, conn := range conns {
				// don't send msg back to the client that sent the msg
				if clientMsg.conn == conn {
					continue
				}
				fmt.Fprint(conn, clientMsg.msg)
			}
			log.Printf("%v: %v", clientMsg.conn.RemoteAddr().String(), clientMsg.msg)
		case newConn := <-newConns:
			conns = append(conns, newConn)
		}
	}
}

type clientMsg struct {
	conn net.Conn
	msg  string
}

func recieveServerMessages(conn net.Conn, msgs chan clientMsg) {
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// send the msg along the channel.
		msgs <- clientMsg{conn, msg}
	}
}

func recieveText(linesOfText chan string) {
	var msg string
	for {
		_, err := fmt.Fscanln(os.Stdin, &msg)
		check(err)
		linesOfText <- msg
	}
}

func recieveMessages(conn net.Conn) {
	clientName := conn.RemoteAddr().String()
	reader := bufio.NewReader(conn)
	for {
		msg, err := reader.ReadString('\n')
		if err != nil {
			break
		}
		// print out the msg
		log.Printf("%v: %v", clientName, msg)
	}
}

func check(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
