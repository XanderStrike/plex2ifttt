package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/xanderstrike/plexhooks"
	"github.com/xanderstrike/plexlights/handler"
)

func hook(w http.ResponseWriter, r *http.Request) {
	log.Println("Serving request...")

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	re := regexp.MustCompile("({.*})")
	match := re.FindStringSubmatch(string(body))

	response, err := plexhooks.ParseWebhook([]byte(match[0]))
	if err != nil {
		panic(err)
	}

	h := handler.New()
	h.HandleEvent(response.Account.Title, response.Player.Uuid, response.Event)
}

func index(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("OK"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hook", hook).Methods("POST")
	router.HandleFunc("/", index)

	log.Println("Now serving on 0.0.0.0:8080/hook")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
