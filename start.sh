#!/usr/bin/env bash
export DB_URL="cailianpress_dba:xxxxxx@tcp(47.93.71.5:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
export WRITER_DB_URL="cailianpress_dba:xxxxxx@tcp(47.93.71.5:3306)/cls_1508?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai"
export APP_PORT="1323"
export REDIS_ADDR=""
export REDIS_PWD=""
go run main.go
