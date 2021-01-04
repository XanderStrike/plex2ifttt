package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"time"

	"github.com/gorilla/mux"
	"github.com/kelvins/sunrisesunset"
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
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		panic(err)
	}
	now := time.Now()

	p := sunrisesunset.Parameters{
		Latitude:  -33.960,
		Longitude: -118.351,
		UtcOffset: -8.0,
		Date:      now,
	}

	sunrise, sunset, err := p.GetSunriseSunset()
	if err != nil {
		panic(err)
	}

	// what the actual fuck is going on with timezones
	eight, _ := time.ParseDuration("-8h")
	log.Println("Sunrise: ", sunrise.In(loc).Format("Mon 3:04PM"), "Sunset: ", sunset.In(loc).Format("Mon 3:04PM"), "Now: ", now.Add(eight).Format("Mon 3:04PM"))
	_, _ = w.Write([]byte("OK"))
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/hook", hook).Methods("POST")
	router.HandleFunc("/", index)

	log.Println("Now serving on 0.0.0.0:8080/hook")
	log.Println("Version 0.0.2")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
