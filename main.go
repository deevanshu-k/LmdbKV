package main

import (
	"log"

	"github.com/deevanshu-k/lmdbkv/config"
	"github.com/deevanshu-k/lmdbkv/store"
	"github.com/deevanshu-k/lmdbkv/subscriber"
	"github.com/deevanshu-k/lmdbkv/writer"
)

func main() {
	log.Print("Starting LmdbKV Server")

	/* Create Store */
	s := store.NewStore(config.STORE_PATH)

	/* Start writer http server */
	go writer.StartHttpWriterServer(s, config.WRITER_HTTP_SERVER_ADDRESS)

	/* Start subscriber tcp server */
	go subscriber.StartTcpSubscriberServer(s, config.SUBSCRIBER_TCP_SERVER_ADDRESS)

	select {}
}
