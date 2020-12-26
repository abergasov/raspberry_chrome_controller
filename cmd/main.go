package main

import (
	"log"
	"net/http"
	"net/url"
	"raspberry_chrome_controller/pkg/config"
	"raspberry_chrome_controller/pkg/logger"
	"raspberry_chrome_controller/pkg/utils"

	"github.com/gorilla/websocket"

	"go.uber.org/zap"
)

type PiCommander struct {
	commando *utils.Commandor
}

func main() {
	logger.NewLogger()
	appConf := config.InitConf("/configs/conf.yaml")
	logger.Info("Starting app", zap.String("token", appConf.KeyToken), zap.String("host", appConf.HostURL))

	piCmd := PiCommander{
		commando: utils.NewCommandor(),
	}
	scheme := "ws"
	if appConf.UseSSL {
		scheme = "wss"
	}
	u := url.URL{Scheme: scheme, Host: appConf.HostURL, Path: appConf.Path}
	for {
		log.Printf("connecting to %s/%s", appConf.HostURL, appConf.Path)
		err := piCmd.connect(&u, appConf.KeyToken)
		if err != nil {
			log.Println("error in socket connection", err.Error())
		}
	}
}

func (p *PiCommander) connect(u *url.URL, token string) error {
	c, _, err := websocket.DefaultDialer.Dial(u.String(), http.Header{"Token": []string{token}})
	if err != nil {
		return err
	}
	defer c.Close()

	for {
		_, message, errR := c.ReadMessage()
		if errR != nil {
			return errR
		}
		log.Printf("recv: %s", message)
		p.commando.HandleCommand(message)
	}
}
