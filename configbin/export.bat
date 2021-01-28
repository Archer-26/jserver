@echo off
rd/s/q .\config
md .\config\config_go
md .\config\config_json

excelExport.exe -goPackageName=config_go ^
-tplPath=.\ ^
-saveGoPath=.\config\config_go ^
-saveJsonPath=.\config\config_json ^
-readPath=E:\prj103\trunk\prj_config\configs

xcopy .\config ..\..\planb\root\internal\config\ /s/e/y
pause