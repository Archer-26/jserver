# PlanB
###### 启动命令 及 参数配置:
>#### -app 进程名
>#### -id  进程id
>#### -log 日志级别
>#### -ini 进程配置路径

#### example:
>#####gate1
>>./planb -app gate -id 1 -log 0 -ini ./ini/servers.ini


######目录结构描述:
- configbin`导出excel配置文件,执行export脚本`
- ini
    - servers.ini `进程启动基础配置 进程名、配置路径、日志路径、监听端口等`
- root
  - cmd
    - excelexport `导表工具`
    - exec        `项目启动入口`  
  - servers `项目服务单元`
    - clients
    - game
    - gate
    - internal
      - inner_message   `服务器内部消息`
      - message         `前后端消息`
      - models          `gorm表结构`
      - mysql
    - login
    - webconsole
  - pkg `公共基础库`
    - actor
    - abtime
    - container
    - ev
    - iniconfig
    - log
    - network
    - tools
  - internal    `本项目使用内部公共代码`
    - coder
    - common
    - config `配置文件`
      - config_go
      - config_json
    - system
  - go.mod
- build_excel.sh `内网服务器更新编译脚本`

