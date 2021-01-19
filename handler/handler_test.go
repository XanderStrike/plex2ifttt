package handler_test

import (
	"os"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/xanderstrike/plexlights/handler"
)

var h handler.Handler

var _ = Describe("Handler", func() {
	Context("Incorrect player or user ids", func() {
		BeforeEach(func() {
			h = handler.Handler{
				Requester: func(msg string) {
					Fail("Don't send requests!")
				},
				Time: func() time.Time {
					Fail("Don't get the time!")
					return time.Now()
				},
			}
		})

		It("Rejects incorrect users", func() {
			os.Setenv("USER_ID", "user")
			os.Setenv("PLAYER_UUID", "player")

			h.HandleEvent("wronguser", "player", "event")
		})

		It("Rejects incorrect players", func() {
			os.Setenv("USER_ID", "user")
			os.Setenv("PLAYER_UUID", "player")

			h.HandleEvent("user", "wrong", "event")
		})
	})

	Context("During the day", func() {
		BeforeEach(func() {
			h = handler.Handler{
				Requester: func(msg string) {
					Fail("Overwrite me!")
				},
				Time: func() time.Time {
					str := "2014-11-12T8:00:26.371Z"
					t, err := time.Parse(time.RFC3339, str)
					Expect(err).NotTo(HaveOccurred())
					return t
				},
			}
		})

		// It("sends the plex_play_day event in a play", func() {
		// 	h.Requester = func(msg string) {
		// 		Expect(msg).To(Equal("plex_play_day"))
		// 	}
		// 	os.Setenv("USER_ID", "user")
		// 	os.Setenv("PLAYER_UUID", "player")

		// 	h.HandleEvent("user", "player", "media.play")
		// 	Expect(true).To(BeTrue())
		// })

		It("sends the plex_pause event in a pause", func() {
			h.Requester = func(msg string) {
				Expect(msg).To(Equal("plex_pause"))
			}
			os.Setenv("USER_ID", "other,user")
			os.Setenv("PLAYER_UUID", "player")

			h.HandleEvent("user", "player,other", "media.pause")
			Expect(true).To(BeTrue())
		})
	})
	Context("During the night", func() {
		BeforeEach(func() {
			h = handler.Handler{
				Requester: func(msg string) {
					Fail("Overwrite me!")
				},
				Time: func() time.Time {
					str := "2014-11-12T10:00:26.371Z"
					t, err := time.Parse(time.RFC3339, str)
					Expect(err).NotTo(HaveOccurred())
					return t
				},
			}
		})

		It("sends the plex_play event in a play", func() {
			h.Requester = func(msg string) {
				Expect(msg).To(Equal("plex_play"))
			}
			os.Setenv("USER_ID", "user")
			os.Setenv("PLAYER_UUID", "player")

			h.HandleEvent("user", "player", "media.play")
			Expect(true).To(BeTrue())
		})

		It("sends the plex_pause event in a pause", func() {
			h.Requester = func(msg string) {
				Expect(msg).To(Equal("plex_pause"))
			}
			os.Setenv("USER_ID", "user")
			os.Setenv("PLAYER_UUID", "player")

			h.HandleEvent("user", "player", "media.pause")
			Expect(true).To(BeTrue())
		})
	})
})
