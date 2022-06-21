# nexttrace-backend

NextTrace BackEnd

## Get Started

修改 Token 配置文件 `nexttrace-backend\ipgeo\basic.go` ，填入 `token`  信息，然后编译运行即可。

下面是一个 `systemd service` 的运行模版

```
[Unit]
Description=Nexttrace Backend
After=network.target

[Service]
ExecStart=/root/nexttrace-backend/nexttrace-backend
ExecReload=/bin/kill -HUP $MAINPID
KillMode=process
Restart=on-failure

[Install]
WantedBy=multi-user.target
```
