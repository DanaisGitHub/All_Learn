package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

const fileName string = "messages.txt"

func main() {


	//raddress:=new(net.UDPAddr)

	udpConn, err := net.Dial("udp", ":42069") // opens up a tcp connection on this port
	if err != nil {
		fmt.Println("coudn't dial up udp connection: ", err)
		os.Exit(1)
	}
	defer udpConn.Close()

	for {
		udpConn.Write([]byte("Hello UDP"))
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {

	chanString := make(chan string, 1)
	go func() {
		defer close(chanString)

		eightBytes := make([]byte, 8)
		var fullLine []byte
		for {
			n, err := f.Read(eightBytes)
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
				char := eightBytes[i]
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
