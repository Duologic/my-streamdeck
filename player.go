package main

import (
	"bytes"
	"log"
	"os/exec"
	"strings"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
)

type playerButtons struct {
	sd            *streamdeck.StreamDeck
	SpotifyButton *buttons.ImageFileButton
	PlayButton    *buttons.ImageFileButton
	PrevButton    *buttons.ImageFileButton
	NextButton    *buttons.ImageFileButton
	HasPlayer     bool
}

const (
	iconSpotify    = "assets/spotify.png"
	iconMediaPrev  = "assets/media-skip-backward.png"
	iconMediaNext  = "assets/media-skip-forward.png"
	iconMediaPlay  = "assets/media-playback-start.png"
	iconMediaPause = "assets/media-playback-pause.png"
	iconMediaStop  = "assets/media-playback-stop.png"
)

func createButton(initImage string) *buttons.ImageFileButton {
	button, err := buttons.NewImageFileButton(initImage)
	if err != nil {
		panic(err)
	}
	button.RegisterUpdateHandler(func(streamdeck.Button) {})
	return button
}

func PlayerButtons(sd *streamdeck.StreamDeck) playerButtons {
	m := playerButtons{
		sd:            sd,
		SpotifyButton: createButton(iconSpotify),
		PlayButton:    createButton(iconMediaPlay),
		PrevButton:    createButton(iconMediaPrev),
		NextButton:    createButton(iconMediaNext),
	}

	m.SpotifyButton.SetActionHandler(actionhandlers.NewCustomAction(m.spotifyHandler))
	m.PlayButton.SetActionHandler(actionhandlers.NewCustomAction(m.playHandler))
	m.PrevButton.SetActionHandler(actionhandlers.NewCustomAction(m.prevHandler))
	m.NextButton.SetActionHandler(actionhandlers.NewCustomAction(m.nextHandler))

	m.Update()

	return m
}

func (m *playerButtons) spotifyHandler(button streamdeck.Button) {
	cmd := exec.Command("/usr/bin/spotify")
	cmd.Start()
}

func (m *playerButtons) playHandler(button streamdeck.Button) {
	command("play-pause")
	m.Update()
}

func (m *playerButtons) prevHandler(button streamdeck.Button) {
	command("previous")
}

func (m *playerButtons) nextHandler(button streamdeck.Button) {
	command("next")
}

func (m *playerButtons) Update() {
	cmd := exec.Command("/usr/bin/playerctl", "status")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
	if strings.Contains(out.String(), "Playing") {
		m.HasPlayer = true
		m.PlayButton.SetFilePath(iconMediaPause)
		return
	}
	if strings.Contains(out.String(), "Paused") {
		m.HasPlayer = true
		m.PlayButton.SetFilePath(iconMediaPlay)
		return
	}
	m.HasPlayer = false
	m.PlayButton.SetFilePath(iconMediaStop)
}

func command(setting string) {
	cmd := exec.Command("/usr/bin/playerctl", setting)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
