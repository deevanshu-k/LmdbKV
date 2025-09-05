package main

import (
	"flag"
	"log"

	"github.com/deevanshu-k/lmdbkv/config"
	"github.com/deevanshu-k/lmdbkv/store"
	"github.com/deevanshu-k/lmdbkv/subscriber"
	"github.com/deevanshu-k/lmdbkv/writer"
)

var log_key = "MAIN"

func init() {
	flag.StringVar(&config.STORE_PATH, "lmdbpath", "", "lmdb storage folder path (eg. /tmp/lmdbPath)")
	flag.StringVar(&config.WRITER_HTTP_SERVER_ADDRESS, "waddress", "", "address for writer http server (eg. 127.0.0.1:5500)")
	flag.StringVar(&config.SUBSCRIBER_TCP_SERVER_ADDRESS, "saddress", "", "address for subscriber tcp server (eg. 127.0.0.1:5500)")
}

func main() {
	log.Printf("[%s] Starting LmdbKV Server", log_key)

	flag.Parse()
	if config.STORE_PATH == "" {
		log.Fatalf("[%s] Error: --lmdbpath is required", log_key)
	}
	if config.WRITER_HTTP_SERVER_ADDRESS == "" {
		log.Fatalf("[%s] Error: --waddress is required", log_key)
	}
	if config.SUBSCRIBER_TCP_SERVER_ADDRESS == "" {
		log.Fatalf("[%s] Error: --saddress is required", log_key)
	}

	/* Create Store */
	s := store.NewStore(config.STORE_PATH)

	/* Start writer http server */
	go writer.StartHttpWriterServer(s, config.WRITER_HTTP_SERVER_ADDRESS)

	/* Start subscriber tcp server */
	go subscriber.StartTcpSubscriberServer(s, config.SUBSCRIBER_TCP_SERVER_ADDRESS)

	select {}
}
