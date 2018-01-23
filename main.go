package main

import (
	"flag"
	"log"
	"net/http"

	test "github.com/curtiscovington/curtiscovington.com/test/pb"

	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var (
	port    string
	tlsPort string
	client  test.TestClient
)

func main() {

	flag.StringVar(&port, "port", "8080", "The port to use for the http connection.")
	flag.StringVar(&tlsPort, "tlsport", "4430", "The port to use for the https connection.")
	flag.Parse()
	fs := http.FileServer(http.Dir("static"))

	conn, err := grpc.Dial(":36060", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to dial: %v", err)
	}
	defer conn.Close()

	client = test.NewTestClient(conn)

	http.Handle("/", fs)

	go func() {
		log.Println("Listening on http://localhost:" + port)
		err := http.ListenAndServe(":"+port, http.HandlerFunc(redir))
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("Listening on https://localhost:" + tlsPort)
	err2 := http.ListenAndServeTLS(":"+tlsPort, "cert.pem", "key.pem", nil)
	if err2 != nil {
		log.Fatal(err2)
	}
}

// Redirect all http requests to https
func redir(w http.ResponseWriter, req *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()
	res, err := client.Test(ctx, &test.Request{Query: "test"})
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("res: %v", res)
	http.Redirect(w, req, "https://localhost:"+tlsPort+req.RequestURI, http.StatusMovedPermanently)
}
