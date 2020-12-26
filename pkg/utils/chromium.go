package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
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
	cmdList := []string{
		"/usr/bin/chromium-browser & sleep 2",
		fmt.Sprintf(`xdotool type "youtube.com/watch?v=%s"`, videoKey),
		"xdotool key Return",
		"xdotool windowactivate $(%s search --name 'Chromium')",
		"xdotool mousemove 600 600",
		"xdotool click 1",
		"xdotool key Space",
	}
	c.execCommand(cmdList)
}

func (c *Commandor) closeAll() {
	c.execCommand([]string{
		"xdotool windowkill $(xdotool search --name 'Chromium')",
	})
	c.muV.Lock()
	c.videos = map[string]string{}
	c.muV.Unlock()
}

func (c *Commandor) findWin() {}

func (c *Commandor) execCommand(commands []string) {
	command := "export DISPLAY=:0.0 && " + strings.Join(commands, " && ")
	cmd := exec.Command("/bin/sh", "-c", command)
	//cmd := exec.Command(command)
	err := cmd.Run()
	if err != nil {
		log.Printf("error execute %s, %s", command, err.Error())
		return
	}
}
