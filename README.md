Instalation
```shell
cd /home/pi && git clone https://github.com/abergasov/raspberry_chrome_controller.git && cd raspberry_chrome_controller
# create configs/conf.yaml
make install

# create autorefresh server
make cron_job
```
## Chromium extenions
```shell
https://chrome.google.com/webstore/detail/youtube-ad-auto-skipper/lokpenepehfdekijkebhpnpcjjpngpnd
https://chrome.google.com/webstore/detail/auto-quality-for-youtube/iaddfgegjgjelgkanamleadckkpnjpjc
```

## software
```shell
sudo apt install mpv 
sudo apt update && sudo apt install snapd && sudo reboot
sudo snap install core && sudo snap install youtube-dl && sudo ln -s /snap/bin/youtube-dl /bin/youtube-dl
```