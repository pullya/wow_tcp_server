package main

import (
	"bufio"
	"fmt"
	"net"
	"time"
)

func main() {
	//ctx := context.Background()

	for i := 0; i < 10; i++ {
		go func() {
			conn, err := net.Dial("tcp", "tcp_server:8081")
			if err != nil {
				fmt.Println("Error while establishing connection to tcp-server: ", err)
				return
			}
			defer conn.Close()

			request := "request1\n"
			fmt.Fprintf(conn, request+"\n")
			message, _ := bufio.NewReader(conn).ReadString('\n')
			fmt.Print("Message from server: " + message)
		}()
		time.Sleep(5 * time.Second)
	}
}
