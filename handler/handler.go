package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
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
	if user != os.Getenv("USER_ID") {
		log.Println("Wrong user:", user)
		return
	}

	if player != os.Getenv("PLAYER_UUID") {
		log.Println("Unkown player:", player)
		return
	}

	if event == "media.play" || event == "media.resume" || event == "media.scrobble" {
		cur := h.Time()
		log.Println("Current UTC hour of the day is", cur.Hour())

		if cur.Hour() < 16 && cur.Hour() > 2 {
			log.Println("Send play (night)")
			h.Requester("plex_play")
		} else {
			log.Println("Send play (day)")
			h.Requester("plex_play_day")
		}
	} else {
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
