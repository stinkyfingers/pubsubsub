package main

import (
	"github.com/stinkyfingers/pubsubsub/pubsub"
	// "github.com/rs/cors"

	"flag"
	"fmt"
	"log"
	"net/http"
)

var (
	port  = flag.String("port", "8081", "Listening Port --port=8081")
	topic = flag.String("topic", "", "If selected, will pull-wait indefinitely for specified topic")
)

func main() {

	flag.Parse()
	if port == nil {
		log.Fatal("failed to setup listening port")
	}

	if topic != nil && *topic != "" {
		go pubsub.PullWait(*topic)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", status)
	mux.HandleFunc("/pull", pull)
	mux.HandleFunc("/wait", pullwait)

	log.Printf("Starting server on %s", *port)
	http.ListenAndServe(fmt.Sprintf(":%s", *port), mux)
}

func status(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Status: OK")
	return
}

func pull(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	err := pubsub.Pull(topic)
	fmt.Println(err)
	return
}

func pullwait(w http.ResponseWriter, r *http.Request) {
	topic := r.URL.Query().Get("topic")
	err := pubsub.PullWait(topic)
	fmt.Println(err)
	return
}
