package discord

import (
	"fmt"
	"github.com/nattawitc/rich-go/client"
	"github.com/wieku/danser-go/app/settings"
	"log"
	"time"
)

const appId = "658093518396588032"

var queue chan func()
var connected bool

var startTime = time.Now()
var endTime = time.Now()
var mapString string

var lastSentActivity = "Idle"

func Connect() {
	if !settings.General.DiscordPresenceOn {
		return
	}

	err := client.Login(appId)
	if err != nil {
		log.Println("Can't login to Discord RPC")
		return
	}

	connected = true

	queue = make(chan func(), 100)

	go func() {
		for {
			f, keepOpen := <-queue

			if f != nil {
				f()
			}

			if !keepOpen {
				break
			}
		}
	}()
}

func SetDuration(duration int64) {
	startTime = time.Now()
	endTime = time.Now().Add(time.Duration(duration) * time.Millisecond)

	sendActivity(lastSentActivity)
}

func SetMap(artist, title, version string) {
	mapString = fmt.Sprintf("%s - %s [%s]", artist, title, version)

	sendActivity(lastSentActivity)
}

func sendActivity(state string) {
	if !connected {
		return
	}

	queue <- func() {
		lastSentActivity = state
		err := client.SetActivity(client.Activity{
			State:      state,
			Details:    mapString,
			LargeImage: "danser-logo",
			Timestamps: &client.Timestamps{
				Start: &startTime,
				End:   &endTime,
			},
		})
		if err != nil {
			log.Println("Can't send activity")
		}
	}
}

func UpdateKnockout(alive, players int) {
	sendActivity(fmt.Sprintf("Watching knockout (%d of %d alive)", alive, players))
}

func UpdatePlay(name string) {
	state := "Clicking circles"
	if name != "" {
		state = fmt.Sprintf("Watching %s", name)
	}

	sendActivity(state)
}

func UpdateDance(tag, divides int) {
	statusText := "Watching "
	if tag > 1 {
		statusText += fmt.Sprintf("TAG%d ", tag)
	}

	if divides > 2 {
		statusText += "mandala"
	} else if divides == 2 {
		statusText += "mirror collage"
	} else {
		statusText += "cursor dance"
	}

	sendActivity(statusText)
}

func ClearActivity() {
	if !connected {
		return
	}

	queue <- func() {
		lastSentActivity = "Idle"

		err := client.ClearActivity()
		if err != nil {
			log.Println("Can't clear activity")
		}
	}
}

func Disconnect() {
	if !connected {
		return
	}

	err := client.ClearActivity()
	if err != nil {
		log.Println("Can't clear activity")
	}

	connected = false

	close(queue)

	client.Logout()
}
