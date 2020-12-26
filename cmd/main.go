package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"raspberry_chrome_controller/pkg/config"
	"raspberry_chrome_controller/pkg/logger"
	"raspberry_chrome_controller/pkg/utils"
	"time"

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

	ticker := time.NewTicker(150 * time.Millisecond)

	var str []byte
	client := &http.Client{}
	req, _ := http.NewRequest("POST", appConf.FullPath, bytes.NewBuffer(str))
	req.Header.Set("Token", appConf.KeyToken)
	for range ticker.C {
		piCmd.curlConnect(client, req)
	}
}

func (p *PiCommander) curlConnect(client *http.Client, req *http.Request) {
	response, err := client.Do(req)
	if err != nil {
		log.Println("error in curl request", err.Error())
		return
	}
	defer response.Body.Close()
	body, errB := ioutil.ReadAll(response.Body)
	if errB != nil {
		log.Println("error in parse curl request", errB.Error())
		return
	}
	var pR struct {
		Ok  bool           `json:"ok"`
		Cmd *utils.Command `json:"cmd"`
	}
	err = json.Unmarshal(body, &pR)
	if err != nil {
		log.Println("error in parse curl request", err.Error())
		return
	}
	if pR.Cmd != nil {
		p.commando.HandleCommand(pR.Cmd)
	}
}
