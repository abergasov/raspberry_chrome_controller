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
	appWindow   string
	muV         *sync.RWMutex
	videos      map[string]string
	xdotool     string
	browserPath string
}

func NewCommandor() *Commandor {
	return &Commandor{
		appWindow:   "Chromium",
		videos:      make(map[string]string),
		muV:         &sync.RWMutex{},
		xdotool:     "/usr/bin/xdotool",
		browserPath: "/usr/bin/chromium-browser",
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
	c.execCommand([]string{
		fmt.Sprintf(`export DISPLAY=:0.0 && %s & sleep 2 && %s type "youtube.com/watch?v=%s" && %s key Return`, c.browserPath, c.xdotool, videoKey, c.xdotool),
		fmt.Sprintf(`%s windowactivate $(%s search --name '%s')`, c.xdotool, c.xdotool, c.appWindow),
	})
}

func (c *Commandor) closeAll() {
	c.execCommand([]string{
		fmt.Sprintf("export DISPLAY=:0.0 && xdotool windowkill $(xdotool search --name '%s')", c.appWindow),
	})
	c.muV.Lock()
	c.videos = map[string]string{}
	c.muV.Unlock()
}

func (c *Commandor) findWin() {}

func (c *Commandor) execCommand(commands []string) {
	for _, command := range commands {
		cmd := exec.Command("/bin/sh", "-c", command)
		//cmd := exec.Command(command)
		err := cmd.Run()
		if err != nil {
			log.Printf("error execute %s, %s", command, err.Error())
			return
		}
	}
}
