package main

import (
	"os/exec"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
)

type meetButton struct {
	sd     *streamdeck.StreamDeck
	Button *buttons.ImageFileButton
	URL    string
}

const (
	iconMeet = "assets/meet.png"
)

func MeetButton(sd *streamdeck.StreamDeck, URL string) meetButton {
	m := meetButton{
		sd:     sd,
		Button: createButton(iconMeet),
		URL:    URL,
	}

	m.Button.SetActionHandler(actionhandlers.NewCustomAction(m.handler))

	return m
}

func (m *meetButton) handler(button streamdeck.Button) {
	cmd := exec.Command("/usr/bin/firefox", "--new-window", m.URL)
	cmd.Start()
}
