package utils

import (
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
		appWindow: "Chromium",
		videos:    make(map[string]string),
		muV:       &sync.RWMutex{},
	}
}

func (c *Commandor) HandleCommand(cmd *Command) {
	log.Println(fmt.Sprintf("handle command action: %s 4 %s", cmd.Cmd, cmd.ActionID))
	switch cmd.Cmd {
	case CLOSE:
		c.closeAll()
	case PLAY:
		c.playWin(cmd.ActionID)
	}
}

func (c *Commandor) playWin(videoKey string) {
	c.closeAll()
	cmd := fmt.Sprintf(`env DISPLAY=:0 && chromium-browser & sleep 2 && xdotool type "youtube.com/watch?v=%s" && xdotool key Return`, videoKey)
	c.execCommand(cmd)
	cmd = fmt.Sprintf(`xdotool windowactivate $(xdotool search --name '%s')`, c.appWindow)
}

func (c *Commandor) closeAll() {
	command := fmt.Sprintf("env DISPLAY=:0 && xdotool windowkill $(xdotool search --name '%s')", c.appWindow)
	c.execCommand(command)
	c.muV.Lock()
	c.videos = map[string]string{}
	c.muV.Unlock()
}

func (c *Commandor) findWin() {}

func (c *Commandor) execCommand(command string) {
	cmd := exec.Command("/usr/bin/bash", "-c", command)
	//cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		log.Printf("error execute %s, %s", command, err.Error())
		return
	}
}
