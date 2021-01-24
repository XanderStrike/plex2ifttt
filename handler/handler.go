package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kelvins/sunrisesunset"
)

type Handler struct {
	Requester func(string)
	Time      func() time.Time
}

func New() Handler {
	return Handler{
		Requester: request,
		Time:      time.Now,
	}
}

func (h Handler) HandleEvent(user, player, event string) {
	if !strings.Contains(os.Getenv("USER_ID"), user) {
		log.Println("Wrong user:", user)
		return
	}

	if !strings.Contains(os.Getenv("PLAYER_UUID"), player) {
		log.Println("Unkown player:", player)
		return
	}

	if event == "media.play" || event == "media.resume" || event == "media.scrobble" {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			panic(err)
		}
		eight, _ := time.ParseDuration("-8h")
		now := h.Time().In(loc).Add(eight) // this is some bullshit but it doesn't work otherwise???

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
		log.Println("Sunrise: ", sunrise.In(loc).Format("Mon 3:04PM"), "Sunset: ", sunset.In(loc).Format("Mon 3:04PM"), "Now: ", now.Format("Mon 3:04PM"))
		if sunrise.In(loc).Hour() < now.In(loc).Hour() && now.In(loc).Hour() < sunset.In(loc).Hour() {
			log.Println("Send play (day)")
			h.Requester("plex_play_day")
		} else {
			log.Println("Send play (night)")
			h.Requester("plex_play")
		}
	} else if event == "media.pause" || event == "media.stop" {
		log.Println("Send pause")
		h.Requester("plex_pause")
	}
}

func request(event string) {
	key := os.Getenv("IFTTT_KEY")
	_, err := http.Get(fmt.Sprintf("https://maker.ifttt.com/trigger/%s/with/key/%s", event, key))
	if err != nil {
		panic(err)
	}
}
