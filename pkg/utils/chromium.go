package utils

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"strings"
	"sync"
	"time"
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
	FULL_SCREEN  = "f"
	VOLUME_UP    = "v"
	VOLUME_DOWN  = "x"
	VOLUME_MUTE  = "m"
)

type Commandor struct {
	appWindow string
	muV       *sync.RWMutex
	playing   string
}

func NewCommandor() *Commandor {
	return &Commandor{
		appWindow: "Chromium",
		playing:   "",
		muV:       &sync.RWMutex{},
	}
}

func (c *Commandor) HandleCommand(cmd *Command) {
	log.Println(fmt.Sprintf("handle command action: %s 4 %s", cmd.Cmd, cmd.ActionID))
	switch cmd.Cmd {
	case VOLUME_MUTE:
		c.execKey("m")
		return
	case VOLUME_DOWN:
		c.execKey("Down")
		c.execKey("Down")
		return
	case VOLUME_UP:
		c.execKey("Up")
		c.execKey("Up")
		return
	case SPEED_UP:
		c.execComboKey("shift", "greater")
		return
	case SPEED_LESS:
		c.execComboKey("shift", "less")
		return
	case BACK:
		c.execComboKey("shift", "Left")
		return
	case BACK_MORE:
		c.execComboKey("shift", "Left")
		c.execComboKey("shift", "Left")
		c.execComboKey("shift", "Left")
		return
	case FORWARD:
		c.execComboKey("shift", "Right")
		return
	case FULL_SCREEN:
		c.execKey("f")
		return
	case FORWARD_MORE:
		c.execComboKey("shift", "Right")
		c.execComboKey("shift", "Right")
		c.execComboKey("shift", "Right")
		return
	case CLOSE:
		c.closeAll()
		return
	case PLAY:
		c.playWin(cmd.ActionID)
		return
	}
	_, err := strconv.Atoi(cmd.Cmd)
	if err == nil {
		c.execKey(cmd.Cmd)
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
	videoKey = strings.ReplaceAll(videoKey, "%", "_")
	if c.playing == videoKey {
		c.execCommand([]string{
			"xdotool mousemove 500 500",
			"xdotool click 1",
		})
		return
	}
	if c.playing != "" {
		c.muV.Lock()
		c.playing = videoKey
		c.muV.Unlock()
		c.execComboKey("ctrl", "t")
		time.Sleep(1 * time.Second)
		c.execCommand([]string{fmt.Sprintf(`xdotool type "youtube.com/watch?v=%s"`, videoKey)})
		time.Sleep(1 * time.Second)
		c.execComboKey("ctrl", "1")
		time.Sleep(1 * time.Second)
		c.execComboKey("ctrl", "w")
		c.execCommand([]string{
			"sleep 2",
			"xdotool mousemove 500 500",
			"xdotool click 1",
			"xdotool key space",
			"xdotool key f",
		})
		return
	}
	//c.closeAll()
	c.execCommand([]string{
		"/usr/bin/chromium-browser > /dev/null & disown",
	})
	c.muV.Lock()
	c.playing = videoKey
	c.muV.Unlock()
	time.Sleep(3 * time.Second)
	c.execCommand([]string{
		fmt.Sprintf(`xdotool type "youtube.com/watch?v=%s"`, videoKey),
		"xdotool key Return",
		//"xdotool windowactivate $(%s search --name 'Chromium')",
		"sleep 4",
		"xdotool mousemove 500 500",
		"xdotool click 1",
		"xdotool key space",
		"xdotool key f",
	})
}

func (c *Commandor) closeAll() {
	c.execCommand([]string{
		"xdotool windowkill $(xdotool search --name 'Chromium')",
	})
	c.muV.Lock()
	c.playing = ""
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
