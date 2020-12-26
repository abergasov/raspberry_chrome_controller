install:
	@echo "-- creating service"
	sudo mkdir -p /etc/systemd/system
	sudo cp commando.service /etc/systemd/system/commando.service
	@echo "-- enable service"
	sudo service commando start && sudo systemctl enable commando

build:
	@echo "-- building puppet binary"
	go build -ldflags "-X main.buildHash=${FILE_HASH} -X main.buildTime=${BUILD_TIME}" -o ./bin/commando ./cmd