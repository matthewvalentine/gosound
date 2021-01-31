package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
)

func main() {
	f, err := os.Open("/home/mokee/Dropbox/Library/Music/Rachmaninoff/Earl Wild plays Liebesleid - YouTube.mp3")
	check(err)
	defer f.Close()

	streamer, format, err := mp3.Decode(f)
	check(err)
	defer streamer.Close()

	check(speaker.Init(format.SampleRate, format.SampleRate.N(5 * time.Minute)))

	done := make(chan struct{})
	ctrl := &beep.Ctrl{Streamer: beep.Seq(
		beep.ResampleRatio(4, 20, streamer),
		beep.Callback(func() { close(done) }),
	)}
	volume := &effects.Volume{
		Streamer: ctrl,
		Base:     2,
		Volume:   0,
		Silent:   false,
	}
	speaker.Play(volume)
	for {
		select {
		case <-done:
			fmt.Println("Done!")
			return
		case <-time.After(10 * time.Second):
			fmt.Println("Pausing!")
			speaker.Lock()
			ctrl.Paused = true
			speaker.Unlock()
		}
	}


}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
