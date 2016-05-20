package main

import (
	"bufio"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func send(host string, buff []byte, recvAmt int, timeout time.Duration) (string, error) {
	conn, err := net.DialTimeout("tcp", host, timeout)
	var n int
	if err != nil {
		return "", err
	}
	defer conn.Close()
	_, err = conn.Write(buff)
	if err != nil {
		return "", err
	}
	recvbuff := make([]byte, recvAmt)
	conn.SetReadDeadline(time.Now().Add(timeout))
	n, err = conn.Read(recvbuff)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(recvbuff[:n]), nil
}

func main() {
	inputfile := flag.String("i", "", "Input file containing newline separated list of hosts. Default: stdin.")
	outfile := flag.String("o", "", "Output file to write the results to. Default: stdout")
	payload := flag.String("send", "id\n", "Payload to send after connection.")
	isHex := flag.Bool("hex", false, "Payload provided to -send is hex encoded.")
	recvAmt := flag.Int("recv", 1024, "How much data to recv before closing the connection.")
	timeoutMS := flag.Int("timeout", 500, "Timeout in millseconds to wait for connection.")
	port := flag.Int("p", 0, "Port to connect to.")
	concurrency := flag.Int("c", 100, "Amount of worker goroutines to spawn.")
	flag.Parse()

	if *port == 0 || *port < 1 || *port > 65535 {
		log.Fatal("Invalid port specification. Please see -h.")
	}
	portstr := strconv.Itoa(*port)

	var file *os.File
	var fileout *os.File
	var err error
	if *inputfile != "" {
		file, err = os.Open(*inputfile)
		checkError(err)
	} else {
		file = os.Stdin
	}

	timeout := time.Duration(*timeoutMS) * time.Millisecond

	if *outfile != "" {
		fileout, err = os.Create(*outfile)
		checkError(err)
	} else {
		fileout = os.Stdout
	}

	var sendbuff []byte
	if *isHex {
		sendbuff, err = hex.DecodeString(*payload)
		checkError(err)
	} else {
		sendbuff = []byte(*payload)
	}

	hosts := make(chan string)
	writer := make(chan string, *concurrency)
	wg := &sync.WaitGroup{}

	for i := 0; i < *concurrency; i++ {
		go func() {
			for host := range hosts {
				recv, err := send(host, sendbuff, *recvAmt, timeout)
				if err == nil {
					writer <- fmt.Sprintf("%s: %s\n", host, recv)
				} else {
					wg.Done()
				}
			}
		}()
	}

	go func(out *os.File) {
		for output := range writer {
			out.WriteString(output)
			wg.Done()
		}
	}(fileout)

	fmt.Fprintf(os.Stderr, "Starting scan\n")

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		ip := fmt.Sprintf("%s:%s", scanner.Text(), portstr)
		wg.Add(1)
		hosts <- ip
	}
	file.Close()

	wg.Wait()

	close(hosts)
	close(writer)
	fileout.Close()

	fmt.Fprintf(os.Stderr, "Scan complete.\n")
}
