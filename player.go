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
	m.Update()
	if !m.HasPlayer {
		cmd := exec.Command("/usr/bin/spotify")
		cmd.Start()
		m.Update()
	}
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
		if err := m.PlayButton.SetFilePath(iconMediaPause); err != nil {
			log.Println(err)
		}
		return
	}
	if strings.Contains(out.String(), "Paused") {
		m.HasPlayer = true
		if err := m.PlayButton.SetFilePath(iconMediaPlay); err != nil {
			log.Println(err)
		}

		return
	}
	if strings.Contains(out.String(), "Stopped") {
		m.HasPlayer = true
		if err := m.PlayButton.SetFilePath(iconMediaStop); err != nil {
			log.Println(err)
		}

		return
	}
	m.HasPlayer = false
	if err := m.PlayButton.SetFilePath(iconSpotify); err != nil {
		log.Println(err)
	}

}

func command(setting string) {
	cmd := exec.Command("/usr/bin/playerctl", setting)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
