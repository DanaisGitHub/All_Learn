package main

import (
	"fmt"
	"http-server/internal/request"
	"io"
	"net"
)

const fileName string = "messages.txt"

func main() {

	tcpConn, err := net.Listen("tcp", ":42069") // opens up a tcp connection on this port
	if err != nil {
		fmt.Println("couldn't listen: ", err)
	}
	defer tcpConn.Close()

	for { // to accept more that one tcp request
		tcpRequest, err := tcpConn.Accept() // accepts / takes in incoming tcp request from that port
		if err != nil {
			fmt.Println("couldn't accept: ", err)
		}
		fmt.Println("tcp server listening")

		// var lines <-chan string = getLinesChannel(tcpRequest)
		req, err := request.RequestFromReader(tcpRequest)
		if err != nil {
			fmt.Printf("couldn't get ")
			return
		}

		fmt.Println(*req)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	chanString := make(chan string, 1)
	go func() {
		defer close(chanString)

		eightB := make([]byte, 8)
		var fullLine []byte
		for {
			n, err := f.Read(eightB)
			if err != nil {
				if err == io.EOF {
					if len(fullLine) > 0 {
						chanString <- string(fullLine)
					}
					return
				}
				fmt.Printf("couidn't read 8 bytes: %v", err)
			}

			// for each char in the 8- bytes
			for i := range n {
				char := eightB[i]
				if char == '\n' {
					chanString <- string(fullLine)
					fullLine = []byte{}
					continue
				}
				fullLine = append(fullLine, char)
			}
		}
	}()

	return chanString // returning an open channel to main routine, whilst go-routine calculates, time is the same bu
}
