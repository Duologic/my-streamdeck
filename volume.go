package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	streamdeck "github.com/magicmonkey/go-streamdeck"
	"github.com/magicmonkey/go-streamdeck/actionhandlers"
	"github.com/magicmonkey/go-streamdeck/buttons"
)

type volumeButtons struct {
	sd           *streamdeck.StreamDeck
	VolumeButton *buttons.ImageFileButton
	MuteButton   *buttons.ImageFileButton
	Volume       int
	Mute         bool
}

const (
	iconVolume      = "assets/multimedia-volume-control.png"
	iconVolumeHigh  = "assets/audio-volume-high.png"
	iconVolumeMed   = "assets/audio-volume-medium.png"
	iconVolumeLow   = "assets/audio-volume-low.png"
	iconVolumeMuted = "assets/audio-volume-muted.png"
)

func VolumeButtons(sd *streamdeck.StreamDeck) volumeButtons {
	volumeButton, err := buttons.NewImageFileButton(iconVolume)
	if err != nil {
		panic(err)
	}
	volumeButton.RegisterUpdateHandler(func(streamdeck.Button) {})

	muteButton, err := buttons.NewImageFileButton(iconVolumeMuted)
	if err != nil {
		panic(err)
	}
	muteButton.RegisterUpdateHandler(func(streamdeck.Button) {})

	m := volumeButtons{
		sd:           sd,
		Volume:       0,
		VolumeButton: volumeButton,
		MuteButton:   muteButton,
	}

	volumeButton.SetActionHandler(actionhandlers.NewCustomAction(m.volumeHandler))
	muteButton.SetActionHandler(actionhandlers.NewCustomAction(m.muteHandler))

	m.Update()

	return m
}

func (m *volumeButtons) volumeHandler(button streamdeck.Button) {
	if m.Volume <= 0 {
		m.Volume = 110
	}
	m.Volume = m.Volume - 10
	log.Println(fmt.Sprintf("%d%%", m.Volume))
	volume(fmt.Sprintf("%d%%", m.Volume))
	m.Update()
}

func (m *volumeButtons) muteHandler(button streamdeck.Button) {
	volume("toggle")
	m.Update()
}

func (m *volumeButtons) SetVolume(value int) {
	m.Volume = value
	icon := iconVolumeLow
	switch {
	case value > 70:
		icon = iconVolumeHigh
	case value <= 30:
		icon = iconVolumeLow
	default:
		icon = iconVolumeMed
	}
	//m.VolumeButton.SetFilePath(icon)
	if m.Mute {
		m.MuteButton.SetFilePath(iconVolumeMuted)
	} else {
		m.MuteButton.SetFilePath(icon)
	}
}

func (m *volumeButtons) Update() {
	cmd := exec.Command("/usr/bin/amixer", "sget", "Master")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}

	mutere := regexp.MustCompile(`\[off\]`)
	found := mutere.FindString(out.String())
	m.Mute = false
	if found == "[off]" {
		m.Mute = true
	}

	re := regexp.MustCompile(`(\d+)%`)
	percent := re.FindString(out.String())
	value, err := strconv.Atoi(strings.ReplaceAll(percent, "%", ""))
	if err != nil {
		log.Println(err)
	}
	m.SetVolume(value)
}

func volume(setting string) {
	cmd := exec.Command("/usr/bin/amixer", "sset", "Master", setting)
	if err := cmd.Run(); err != nil {
		log.Println(err)
	}
}
