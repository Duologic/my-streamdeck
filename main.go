package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
)

func WaitForCtrlC() {
	var end_waiter sync.WaitGroup
	end_waiter.Add(1)
	var signal_channel chan os.Signal
	signal_channel = make(chan os.Signal, 1)
	signal.Notify(signal_channel, os.Interrupt)
	go func() {
		<-signal_channel
		end_waiter.Done()
		os.Exit(0)
	}()
	end_waiter.Wait()
}

func main() {
	rawsd, err := streamdeck.Open()
	if err != nil {
		panic(err)
	}
	rawsd.ClearButtons()
	rawsd.Close()

	sd, err := streamdeck.New()
	if err != nil {
		panic(err)
	}
	sd.SetBrightness(50)

	mute := MicButtons(sd)
	volume := VolumeButtons(sd)
	player := PlayerButtons(sd)

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			mute.Update()
			volume.Update()
			player.Update()
		}
	}()

	sd.AddButton(0, mute.MuteButton)
	sd.AddButton(1, mute.VolumeButton)
	sd.AddButton(2, volume.VolumeButton)
	sd.AddButton(3, volume.MuteButton)
	sd.AddButton(4, player.SpotifyButton)
	sd.AddButton(5, player.PrevButton)
	sd.AddButton(6, player.PlayButton)
	sd.AddButton(7, player.NextButton)

	WaitForCtrlC()
}
