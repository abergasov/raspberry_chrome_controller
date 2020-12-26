package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

type Command struct {
	Cmd      string // команда для выполнения
	ActionID string // для какого видео
}

const (
	SPEED_0_75   = "a"
	SPEED_1      = "b"
	SPEED_1_25   = "c"
	SPEED_1_5    = "d"
	SPEED_1_75   = "e"
	SPEED_2      = "f"
	FORWARD      = "g"
	FORWARD_MORE = "h"
	BACK         = "i"
	BACK_MORE    = "j"
	PLAY         = "k"
	CLOSE        = "l"
)

type Commandor struct {
	appWindow string
	muV       *sync.RWMutex
	videos    map[string]string
}

func NewCommandor() *Commandor {
	return &Commandor{
		appWindow: "Firefox",
		videos:    make(map[string]string),
		muV:       &sync.RWMutex{},
	}
}

func (c *Commandor) HandleCommand(command []byte) {
	var cmd Command
	err := json.Unmarshal(command, &cmd)
	if err != nil {
		log.Printf("error parse command %s", string(command))
		return
	}

	switch cmd.Cmd {
	case CLOSE:
		c.closeAll()
	case PLAY:
		c.playWin(cmd.Cmd, cmd.ActionID)
	}
}

func (c *Commandor) playWin(command, videoKey string) {
	winID := ""
	c.muV.Lock()
	c.videos[videoKey] = winID
	c.muV.Unlock()
}

func (c *Commandor) closeAll() {
	command := fmt.Sprintf("export DISPLAY=:0.0 && xdotool windowkill $(xdotool search --name '%s')", c.appWindow)
	c.execCommand(command)
	c.muV.Lock()
	c.videos = map[string]string{}
	c.muV.Unlock()
}

func (c *Commandor) findWin() {}

func (c *Commandor) execCommand(command string) {
	cmd := exec.Command("/bin/sh", "-c", command)
	err := cmd.Run()
	if err != nil {
		log.Printf("error execute %s, %s", command, err.Error())
		return
	}
}
