package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
)

var appID = ""
var trustedFacets = []string{appID}

func main() {
	var _appID = flag.String("app-id", "localhost", "appID (the host)")
	var port = flag.Int("port", 3000, "the port to listen on")
	flag.Parse()
	appID = fmt.Sprintf("https://%v:%v", *_appID, *port)
	trustedFacets[0] = appID

	log.Println("Lets go!")
	http.HandleFunc("/register", Register)
	http.HandleFunc("/login", Login)
	http.Handle("/", http.FileServer(http.Dir("./static")))

	log.Fatal(http.ListenAndServeTLS(fmt.Sprintf(":%v", *port), "tls.crt", "tls.key", nil))
}
