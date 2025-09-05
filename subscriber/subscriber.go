package subscriber

import (
	"bufio"
	"encoding/json"
	"log"
	"net"
	"strings"

	"github.com/deevanshu-k/lmdbkv/store"
)

var log_key = "SUBSCRIBER-SERVER"

func StartTcpSubscriberServer(s *store.Store, address string) {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("[%s] Error from listner: %v", log_key, err)
	}
	defer listener.Close()

	log.Printf("[%s] Listening on %s", log_key, address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("[%s] Error while accepting connection: %v", log_key, err)
		}

		go handleConnection(conn, s)
	}
}

func handleConnection(conn net.Conn, s *store.Store) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	ch := make(chan store.State)
	clientId := s.GetUniqueClient()

	s.RegisterUserChannel(clientId, ch)

	go func() {
		for {
			data, ok := <-ch
			if !ok {
				break
			}

			log.Printf("[%s] %v", clientId, data)

			if bytes, err := json.Marshal(data); err == nil {
				bytes = append(bytes, '\n')
				_, _ = conn.Write(bytes)
			}
		}
	}()

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[%s] [%s] Connection closed: %v", log_key, clientId, err)
			s.UnRegisterUser(clientId)
			return
		}

		message = strings.TrimSpace(message)
		commands := strings.Split(message, " ")

		if len(commands) != 2 || (commands[0] != "SUB" && commands[0] != "UNSUB") {
			continue
		}

		/* SUB <string> */
		if commands[0] == "SUB" {
			s.Subscribe(clientId, commands[1])
		}

		/* UNSUB <string> */
		if commands[0] == "UNSUB" {
			s.UnSubscribe(clientId, commands[1])
		}

		// send back ACK
		response := "OK\n"
		_, _ = conn.Write([]byte(response))
	}
}
