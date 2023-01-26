package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

type spy struct {
	prefix string
	raddr  net.Addr
}

var CommitSha string

var clients = make(chan *net.TCPConn)

func (l spy) Write(b []byte) (n int, err error) {
	log.Printf(fmt.Sprintf("\n%s %s \n %s\n", l.prefix, l.raddr.String(), b))
	return len(b), nil
}

func (l spy) Close() error {
	return nil
}

func env(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		log.Fatalf("Error, %s environment variable is not set.", key)
	}
	return v
}

func transfer(dst, src *net.TCPConn, prefix string) {
	defer dst.CloseWrite()
	defer src.CloseRead()
	_, err := io.Copy(dst, io.TeeReader(src, spy{prefix, dst.RemoteAddr()}))
	if err != nil {
		log.Printf("Error transfering data : %s", err.Error())
	}
}

func forward(raddr *net.TCPAddr) {
	for {
		go func(src *net.TCPConn) {
			dst, err := net.DialTCP("tcp", nil, raddr)
			if err == nil {
				go transfer(dst, src, ">>>>>")
				go transfer(src, dst, "<<<<<")
			} else {
				log.Printf("error dialing %s : %s", raddr.IP, raddr.Port, err.Error())
			}
		}(<-clients)
	}
}

func init() {
	target := fmt.Sprintf("%s:%s",
		env("TARGET_HOST"),
		env("TARGET_PORT"),
	)
	raddr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		log.Fatalf("Error, can not establish a connection to %s : %s", target, err.Error())
	}
	go forward(raddr)
}

func main() {
	port := env("TCP_LOGGER_PORT")
	laddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		log.Fatalf("Error, TCP_LOGGER_PORT is not a valid address : %s", err.Error())
	}
	server, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		log.Fatalf("Error, could not start server : %s", err.Error())
	}
	log.Printf("Tcp Logger version %s, listening on %s", CommitSha, port)
	for {
		client, err := server.AcceptTCP()
		if err == nil {
			clients <- client
		} else {
			log.Printf("Error processing request : %s", err.Error())
		}
	}
}
