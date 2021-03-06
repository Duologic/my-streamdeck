package main

import (
	"bytes"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
)

type micButtons struct {
	sd           *streamdeck.StreamDeck
	VolumeButton *buttons.ImageFileButton
	MuteButton   *buttons.ImageFileButton
	Status       int
	Muted        bool
}

const (
	iconMicHigh   = "assets/microphone-sensitivity-high.png"
	iconMicMedium = "assets/microphone-sensitivity-medium.png"
	iconMicLow    = "assets/microphone-sensitivity-low.png"
	iconMicMuted  = "assets/microphone-sensitivity-muted.png"
)

func MicButtons(sd *streamdeck.StreamDeck) micButtons {
	volumeButton, err := buttons.NewImageFileButton(iconMicMuted)
	if err != nil {
		panic(err)
	}
	volumeButton.RegisterUpdateHandler(func(streamdeck.Button) {})

	muteButton, err := buttons.NewImageFileButton(iconMicMuted)
	if err != nil {
		panic(err)
	}
	muteButton.RegisterUpdateHandler(func(streamdeck.Button) {})

	m := micButtons{
		sd:           sd,
		Status:       0,
		VolumeButton: volumeButton,
		MuteButton:   muteButton,
	}

	volumeButton.SetActionHandler(actionhandlers.NewCustomAction(m.volumeHandler))
	muteButton.SetActionHandler(actionhandlers.NewCustomAction(m.muteHandler))

	m.Update()

	return m
}

func (m *micButtons) volumeHandler(button streamdeck.Button) {
	switch m.Status {
	case 3:
		m.SetVolume(1)
	case 1:
		m.SetVolume(2)
	default:
		m.SetVolume(3)
	}
}

func (m *micButtons) SetVolume(status int) {
	switch status {
	case 1: // Max volume
		capture("100%")
	case 2: // Medium volume
		capture("70%")
	default: // Low volume
		capture("30%")
	}
	m.Update()
}

func (m *micButtons) muteHandler(button streamdeck.Button) {
	capture("toggle")
	m.Update()
}

func (m *micButtons) SetStatus(status int) {
	icon := iconMicLow
	switch status {
	case 1: // Max volume
		icon = iconMicHigh
		m.Status = 1
	case 2: // Medium volume
		icon = iconMicMedium
		m.Status = 2
	default: // Low volume
		icon = iconMicLow
		m.Status = 3
	}
	m.VolumeButton.SetFilePath(icon)
	if m.Muted {
		m.MuteButton.SetFilePath(iconMicMuted)
	} else {
		m.MuteButton.SetFilePath(icon)
	}
}

func (m *micButtons) Update() {
	cmd := exec.Command("/usr/bin/amixer", "sget", "Capture")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	mutere := regexp.MustCompile(`\[off\]`)
	found := mutere.FindString(out.String())
	m.Muted = false
	if found == "[off]" {
		m.Muted = true
	}

	re := regexp.MustCompile(`(\d+)%`)
	percent := re.FindString(out.String())
	value, err := strconv.Atoi(strings.ReplaceAll(percent, "%", ""))
	if err != nil {
		log.Println(err)
	}
	switch {
	case value > 70:
		m.SetStatus(1)
	case value <= 30:
		m.SetStatus(3)
	default:
		m.SetStatus(2)
	}
}

func capture(setting string) {
	cmd := exec.Command("/usr/bin/amixer", "sset", "Capture", setting)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
