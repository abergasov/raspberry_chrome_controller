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
	SPEED_LESS   = "a"
	SPEED_UP     = "b"
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
	playing   bool
}

func NewCommandor() *Commandor {
	return &Commandor{
		appWindow: "Chromium",
		playing:   false,
		muV:       &sync.RWMutex{},
	}
}

func (c *Commandor) HandleCommand(cmd *Command) {
	log.Println(fmt.Sprintf("handle command action: %s 4 %s", cmd.Cmd, cmd.ActionID))
	switch cmd.Cmd {
	case SPEED_UP:
		c.execComboKey("shift", "greater")
	case SPEED_LESS:
		c.execComboKey("shift", "less")
	case BACK:
		c.execComboKey("shift", "Left")
	case BACK_MORE:
		c.execComboKey("shift", "Left")
		c.execComboKey("shift", "Left")
		c.execComboKey("shift", "Left")
	case FORWARD:
		c.execComboKey("shift", "Right")
	case FORWARD_MORE:
		c.execComboKey("shift", "Right")
		c.execComboKey("shift", "Right")
		c.execComboKey("shift", "Right")
	case CLOSE:
		c.closeAll()
	case PLAY:
		c.playWin(cmd.ActionID)
	}
}

func (c *Commandor) execKey(key string) {
	c.execCommand([]string{
		//"xdotool windowactivate $(%s search --name 'Chromium')",
		"xdotool mousemove 500 500",
		fmt.Sprintf("xdotool key %s", key),
	})
}

func (c *Commandor) execComboKey(key1, key2 string) {
	c.execCommand([]string{
		//"xdotool windowactivate $(%s search --name 'Chromium')",
		"xdotool mousemove 500 500",
		fmt.Sprintf("xdotool keydown %s keydown %s", key1, key2),
		"sleep 0.1",
		fmt.Sprintf("xdotool keyup %s keyup %s", key1, key2),
	})
}

func (c *Commandor) playWin(videoKey string) {
	if c.playing {
		c.execCommand([]string{
			"xdotool mousemove 500 500",
			"xdotool click 1",
		})
		return
	}
	c.closeAll()
	cmdList := []string{
		"/usr/bin/chromium-browser & sleep 2",
		fmt.Sprintf(`xdotool type "youtube.com/watch?v=%s"`, videoKey),
		"xdotool key Return",
		//"xdotool windowactivate $(%s search --name 'Chromium')",
		"sleep 4",
		"xdotool mousemove 500 500",
		"xdotool click 1",
		"xdotool key space",
	}
	c.execCommand(cmdList)
	c.muV.Lock()
	c.playing = true
	c.muV.Unlock()
}

func (c *Commandor) closeAll() {
	c.execCommand([]string{
		"xdotool windowkill $(xdotool search --name 'Chromium')",
	})
	c.muV.Lock()
	c.playing = false
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
