#!/usr/bin/env bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build main.go
git add .
git commit -m "$1"
git push
docker build -t registry.cn-beijing.aliyuncs.com/lanjing/cailianpress-cron:$1 .
docker push registry.cn-beijing.aliyuncs.com/lanjing/cailianpress-cron:$1
rm main
