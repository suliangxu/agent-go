## Agent-go

### 下载insatll.sh
安装了wget的linux系统直接通过
~~~
wget -O install.sh http://172.18.22.20:8888/Resource/../Resource/resource/downloadFile?fileUid=c6de5ab8253d4b1299804bc724506651
~~~
从服务器下载install.sh
未安装wget通过curl进行安装

~~~
curl -X POST "http://172.18.22.20:8888/Resource/../Resource/resource/downloadFile?fileUid=c6de5ab8253d4b1299804bc724506651" -H "accept: application/json;charset=UTF-8" -H "Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiJ7XCJ1c2VySWRcIjpcIjNhOTM1MDYwMTE1ODQ1ZGFiMWYxN2Y2YzlmZjBhMzU1XCIsXCJuYW1lXCI6XCLolKHlrpfpnJZcIixcInVzZXJOYW1lXCI6XCJjYWl6b25nbGluXCIsXCJsb2dpbk9yZ0lkXCI6XCIwOGE2ODY5Y2FjNGI0MGYyYjQ1YWZkZDU3NGE3OWQxMVwiLFwiY2xpZW50VHlwZVwiOlwicGNcIixcInVpZFwiOlwiZTcwYWQ1MWM0MzZhNDI2N2I4YTZmNjhhMzI2NjhjYTlcIn0iLCJpc3MiOiJKU1dMIiwiaWF0IjoxNTkwMzc4OTI0LCJqdGkiOiJlNzBhZDUxYzQzNmE0MjY3YjhhNmY2OGEzMjY2OGNhOSJ9.aiAhEYGB4gsmg22sTIupZ1eKaD6ChMEOTknmDrel2tg"
~~~
后面的toke失效时需要重新获取新的链接



### 环境安装与程序运行

运行install.sh从服务器下载agent程序主体和unzip安装文件。
install.sh自动安装unzip命令对agent程序主体进行解压、配置项目的环境和配置文件，最后运行项目。

~~~sh
sh install.sh
~~~



可通过 **logs/log.txt** 路径查看日志文件。



### 附录：

**程序运行**

若你需要改变服务端口号、脚本接口、脚本执行间隔等参数，请参考以下步骤。

#### 1、修改配置文件

可在 agent-go-linux/configure** 文件夹的 **agent.conf** 文件修改配置

可在 **/home/你的用户/monitor/agent-go-linux/configure** 文件夹的 **agent.conf** 文件修改配置



#### 2、运行

**（1）运行监控程序**

~~~sh
nohup /home/你的用户/monitor/agent-go-linux/daemon/daemon &
~~~



**（2）运行agent**

~~~
nohup /home/你的用户/monitor/agent-go-linux/bin/main /home/你的用户/monitor/agent-go-linux/configure/agent.conf &
~~~

前面为项目bin文件夹下**main.sh可执行文件的目录**，后面所带的参数为 **上面修改的配置文件** 的目录。
