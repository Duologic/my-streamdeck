package main

import (
	"os"
	"os/signal"
	"sync"
	"time"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/buttons"
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

func createButton(initImage string) *buttons.ImageFileButton {
	button, err := buttons.NewImageFileButton(initImage)
	if err != nil {
		panic(err)
	}
	button.RegisterUpdateHandler(func(streamdeck.Button) {})
	return button
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

	mic := MicButtons(sd)
	volume := VolumeButtons(sd)
	player := PlayerButtons(sd)

	go func() {
		ticker := time.NewTicker(1000 * time.Millisecond)
		defer ticker.Stop()

		for range ticker.C {
			mic.Update()
			volume.Update()
			player.Update()
		}
	}()

	sd.AddButton(0, mic.MuteButton)
	sd.AddButton(1, mic.VolumeButton)
	sd.AddButton(2, volume.VolumeButton)
	sd.AddButton(3, volume.MuteButton)
	sd.AddButton(5, player.PrevButton)
	sd.AddButton(6, player.PlayButton)
	sd.AddButton(7, player.NextButton)

	URL := os.Getenv("TEAM_MEET")
	if URL != "" {
		fwj := MeetButton(sd, URL)
		sd.AddButton(4, fwj.Button)
	}

	WaitForCtrlC()
}
