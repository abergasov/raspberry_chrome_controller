install:
	@echo "-- creating service"
	sudo mkdir -p /home/pi/.config/autostart/
	sudo cp commando.desktop /home/pi/.config/autostart/commando.desktop

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