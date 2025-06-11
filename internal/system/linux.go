//go:build linux

package system

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/godbus/dbus/v5"
)

type LinuxSystemController struct {
	conn *dbus.Conn
}

func NewSystemController() (*LinuxSystemController, error) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return nil, err
	}

	return &LinuxSystemController{
		conn: conn,
	}, nil
}

func (sc *LinuxSystemController) GetCurrentPlayingTrackTitle() (string, error) {
	var metaDataReply map[string]dbus.Variant

	obj := sc.conn.Object("org.mpris.MediaPlayer2.spotify", "/org/mpris/MediaPlayer2")
	err := obj.Call("org.freedesktop.DBus.Properties.Get", 0, "org.mpris.MediaPlayer2.Player", "Metadata").Store(&metaDataReply)

	if err != nil {
		return "", err
	}

	artists := metaDataReply["xesam:artist"].Value().([]string)
	title := metaDataReply["xesam:title"].Value().(string)

	if title == "" {
		return "", errors.New("Spotify is running but not playing a track")
	}

	return fmt.Sprintf(
		"%s - %s",
		strings.Join(artists, ","),
		title,
	), nil
}

func (sc *LinuxSystemController) OpenURLInBrowser(url string) error {
	return exec.Command("open", url).Start()
}
