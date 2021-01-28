#!/bin/bash
rm -rf ./config
mkdir ./config
mkdir ./config/config_go
mkdir ./config/config_json
./excelExport_mac -goPackageName=config_go \
-tplPath ./ \
-saveGoPath ./config/config_go \
-saveJsonPath ./config/config_json \
-readPath /Users/hr/Desktop/svn/prj103/prj_config/configs
go fmt ./config/config_go/*
rm -rf ../root/internal/config/config_go/*
rm -rf ../root/internal/config/config_json/*
cp -rpf ./config/ ../root/internal/config/
rm -rf ./config