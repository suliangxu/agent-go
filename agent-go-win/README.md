## Agent-go



### 拷贝项目


~~~
把主体文件agent_go_win拷贝到需要运行的pc上
~~~


### 环境安装与程序运行

**要求python版本为3.7.x，且含有第三方模块 psutil 和 wmi**

程序根目录\environmentInstall\python下存放有配置环境相关的安装包

#Python的安装
~~~
使用python-3.7.7-amd64.exe安装python3.7.7安装时勾选Add Python 3.7 to PATH 把python加入环境变量中

~~~
#Python库的安装
~~~
安装方法1:
离线安装（要求要将Python添加到环境变量里）：
1、查看python路径
where python

2、将 environmentInstall/python 下的文件夹（psutil-5.7.0 和 WMI-1.5.1）放到Python安装目录下的Lib\site_packages 中

3、打开命令行窗口，分别进入psutil-5.7.0和WMI-1.5.1文件夹中，输入命令：
python setup.py install

4、等待安装成功即可。

5、安装psutil时出现  raise distutils.errors.DistutilsPlatformError(err)
              distutils.errors.DistutilsPlatformError: 
              Microsoft Visual C++ 14.0 is required. Get it with "
              Microsoft Visual C++ Build Tools": https://visualstudio.microsoft.com/downloads/
原因是本机缺少 Microsoft Visual C++ 和.NET Framework 4.6包
可以按照以下步骤安装或选择方法二安装psutil
  1）安装.NET Framework 4.6包
  2）再安装Microsoft Visual C++ 14.0软件
  
6、安装wmi时出现 Couldn't find index page for 'pywin32'
原因是本机没有pywin32 使用目录下 pywin32-221.win-amd64-py3.7.exe 安装pywin32后再次安装wmi

安装方法2:
使用pip包管理工具（需要联网且速度较慢）：
pip install psutil
pip install wmi

~~~


#### 1、修改配置文件

可在 **程序主体根目录/configure** 文件夹的 **agent.conf** 文件修改配置

#### 2、运行

**（1）运行监控程序**

~~~
程序根目录/daemon/daemon.exe
~~~



**（2）运行agent**

~~~
程序根目录/bin/main.exe 空格 程序根目录/configure/agent.conf
~~~

前面为项目bin文件夹下**main.exe可执行文件的目录**，后面所带的参数为 **上面修改的配置文件。



可通过 **程序根目录/logs/log.txt** 路径查看日志文件。
