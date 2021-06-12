FILE_HASH?=$(shell git rev-parse HEAD)
BUILD_TIME?=$(shell date -u '+%Y-%m-%d_%H:%M:%S')

install:
	@echo "-- pull changes"
	git pull origin master
	@echo "-- creating service"
	sudo mkdir -p /home/pi/.config/autostart/
	sudo cp commando.desktop /home/pi/.config/autostart/commando.desktop
	sudo cp configs/conf.yaml /etc/commando.yaml
	if pgrep commando; then pkill commando; fi
	sudo cp bin/commando /usr/bin/commando
	sudo chmod +x /usr/bin/commando
	commando

cron_job:
	@echo "-- creating cron job"
	crontab -l | { cat; echo "*/15 * * * * cd /home/pi/raspberry_chrome_controller && make install"; } | crontab -

build_rasp:
	env GOOS=linux GOARCH=arm GOARM=5 go build -ldflags "-X main.buildHash=${FILE_HASH} -X main.buildTime=${BUILD_TIME}" -o ./bin/commando ./cmd

build:
	@echo "-- pull changes"
	git pull origin master
	@echo "-- building puppet binary"
	go build -ldflags "-X main.buildHash=${FILE_HASH} -X main.buildTime=${BUILD_TIME}" -o ./bin/commando ./cmd

run:
	@echo "-- pull changes"
	git pull origin master
	@echo "-- building puppet binary"
	go build -ldflags "-X main.buildHash=${FILE_HASH} -X main.buildTime=${BUILD_TIME}" -o ./bin/commando ./cmd
	./bin/commando

proto:
	@echo "-- generate proto"
	protoc ./api/* --go_out=plugins=grpc:./pkg/grpc