package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"../src/chatroom"
	"log"
)

const (
	CONN_HOST = "localhost"

	CONN_TYPE = "tcp"
	MAX_DATA_RECV = 9999
	BACKLOG = 50
	IP = "134.226.214.254"
)

var (
	port = "8070"
)

func main() {

	port = os.Args[1]
	fmt.Println("Start!\nIP:", getIpAddress())

	listen, err := net.Listen("tcp4", ":"+port)
	defer listen.Close()
	if err != nil {
		log.Fatalf("Socket listen port %d failed,%s", port, err)
		os.Exit(1)
	}
	log.Printf("Begin listen port: %d", port)

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go handleRequest(conn)
	}

}

func handleRequest(conn net.Conn) {

	for {
		// Make a buffer to hold incoming data.
		buf := make([]byte, MAX_DATA_RECV)

		// Read the incoming connection into the request.
		reqLen, err := conn.Read(buf)
		request := string(buf[:reqLen])

		if err == nil {
			fmt.Println(request)
			processRequest(conn, request)
		}
	}

}


// Handles incoming requests.
func processRequest(conn net.Conn, request string) {




	if(request == "KILL_SERVICE\n") {

		fmt.Printf("Kill the service\n" )
		chatroom.Kill()
		conn.Close()

	} else if(request == "HELO BASE_TEST\n") {

		fmt.Printf("Send back Hello\n" )
		//"HELO text\nIP:[ip address]\nPort:[port number]\nStudentID:[your student ID]\n"
		ip := getIpAddress()
		returnMessage := "HELO BASE_TEST\nIP:" + ip + "\nPort:" + port + "\nStudentID:" + "13329643" + "\n"
		fmt.Printf(returnMessage)
		conn.Write([]byte(returnMessage))

	} else if(strings.Contains(request, "JOIN_CHATROOM")) {

		fmt.Printf("It is a JOIN CHATROOM REQUEST\n")
		chatroom.RequestJoinChatroom(request, conn, port)

	} else if(strings.Contains(request, "LEAVE_CHATROOM")) {

		fmt.Printf("It is a LEAVE CHATROOM REQUEST\n")
		chatroom.RequestLeavingChatroom(request, conn, port)

	} else if(strings.Contains(request, "CHAT")) {

		fmt.Printf("It is a CHAT REQUEST\n")
		chatroom.RequestSendMessage(request, conn, port)

	} else if(strings.Contains(request, "DISCONNECT")) {

		fmt.Printf("It is a Disconnect REQUEST\n")
		chatroom.RequestDisconnect(request, conn, port)

	} else {

		fmt.Printf("Nothing interesting\n")
		conn.Write([]byte("Nothing interesting."))

	}

	fmt.Printf("Task Complete\n")
	fmt.Printf("---------------------------------------------\n\n\n")

}

func getIpAddress() string{

	netInterfaceAddresses, err := net.InterfaceAddrs()

	if err != nil { return "" }

	for _, netInterfaceAddress := range netInterfaceAddresses {

		networkIp, ok := netInterfaceAddress.(*net.IPNet)

		if ok && !networkIp.IP.IsLoopback() && networkIp.IP.To4() != nil {

			ip := networkIp.IP.String()

			fmt.Println("Resolved Host IP: " + ip)

			return ip
		}
	}
	return ""
}