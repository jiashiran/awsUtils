#!/usr/bin/env bash

module=$1

echo 'build ' $module

cd ../../pkg/$module

export GOOS=linux
pwd
go build $module.go
mv $module chromedp
chmod +x chromedp
mv chromedp ../../build/chromedp/

cd ../../build/chromedp/
pwd
cp ../../resource/data.csv .

docker build -t registry.cn-beijing.aliyuncs.com/dev/chromedp .

rm -rf data.csv chromedp

docker push registry.cn-beijing.aliyuncs.com/dev/chromedp

#docker run -d -v /var/58/:/var/58/ --log-opt max-size=10m --name=spider