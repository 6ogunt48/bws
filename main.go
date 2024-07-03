package main

import (
	"bufio"
	"bytes"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	ln, err := net.Listen("tcp", ":80")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	gid := getGID()
	fmt.Println("Handling connection in goroutine: ", gid)

	// Simulate a delay for concurrency in action
	time.Sleep(20 * time.Second)

	reader := bufio.NewReader(conn)
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading:", err)
		return
	}

	path := extractPath(requestLine)
	if path == "/" { // Convention to serve index.html for root path
		path = "/index.html"
	}

	// construct file path
	filePath := "www" + path

	// Ensure the file path is withing the www directory
	absoluteFilePath, err := filepath.Abs(filePath)
	if err != nil {
		fmt.Println("Error getting absolute path:", err)
		notFoundResponse := "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\n404 Not Found\n"
		conn.Write([]byte(notFoundResponse))
		return
	}

	wwwDir, err := filepath.Abs("www")
	if err != nil {
		fmt.Println("Error getting www directory absolute path:", err)
		notFoundResponse := "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\n404 Not Found\n"
		conn.Write([]byte(notFoundResponse))
		return
	}

	// Check if the file is within the www directory
	if !strings.HasPrefix(absoluteFilePath, wwwDir) {
		fmt.Println("Attempt to access file outside www directory:", absoluteFilePath)
		notFoundResponse := "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\n404 Not Found\n"
		conn.Write([]byte(notFoundResponse))
		return
	}

	fileContent, err := os.ReadFile(absoluteFilePath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		notFoundResponse := "HTTP/1.1 404 Not Found\r\nContent-Type: text/plain\r\n\r\n404 Not Found\n"
		conn.Write([]byte(notFoundResponse))
		return
	}

	// write HTTP response
	header := "HTTP/1.1 200 OK\r\nContent-Type: text/html\r\n\r\n"
	conn.Write([]byte(header + string(fileContent)))

}

func extractPath(requestLine string) string {
	parts := strings.Split(requestLine, " ")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[1])
	}
	return "/"
}

// getGID returns the goroutine ID
func getGID() uint64 {
	b := make([]byte, 64)                         // Allocate a byte slice of size 64
	b = b[:runtime.Stack(b, false)]               // Capture the current stack trace and re-slice to contain only bytes written by runtime
	b = bytes.TrimPrefix(b, []byte("goroutine ")) // Trim the "goroutine " prefix from the stack trace
	b = b[:bytes.IndexByte(b, ' ')]               // Isolate the goroutine ID from the rest of the information
	n, _ := strconv.ParseUint(string(b), 10, 64)  // Convert the isolated goroutine ID to a uint64
	return n                                      // Return the goroutine ID
}
