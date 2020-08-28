localpath=`pwd`
cp agent-go-linux/logs/log_erron.txt $localpath 2>>$localpath/install.log
cp agent-go-linux/logs/log_info.txt $localpath 2>>$localpath/install.log
rm -r -f agent-go-linux
echo "下载agent主体程序"
curl -X POST "http://172.18.23.50/Resource/resource/downloadFile?fileUid=2e8be8452f5e4dd2b83c52c18830a65f" -H "accept: application/json;charset=UTF-8">>agent-go-linux.zip
echo "下载rpm包"
res=`rpm -q centos-release | grep "el6"`
if [ $? -eq 0 ]; then
  curl -X POST "http://172.18.23.50/Resource/resource/downloadFile?fileUid=rpmtar6" -H "accept: application/json;charset=UTF-8">>rpm.tar.gz
  tar -zxvf rpm.tar.gz >>$localpath/install.log 2>&1
  cd rpm_6
  rpm -Uvh  *.rpm  --nodeps  --force >>$localpath/install.log 2>&1
  cd ..
  rm -r -f rpm_6
fi
res=`rpm -q centos-release | grep "el7"`
if [ $? -eq 0 ]; then
  curl -X POST "http://172.18.23.50/Resource/resource/downloadFile?fileUid=rpmtar7" -H "accept: application/json;charset=UTF-8">>rpm.tar.gz
  tar -zxvf rpm.tar.gz >>$localpath/install.log 2>&1
  cd rpm_7
  rpm -Uvh  *.rpm  --nodeps  --force >>$localpath/install.log 2>&1
  cd ..
  rm -r -f rpm_7
fi
res=`rpm -q centos-release | grep "el8"`
if [ $? -eq 0 ]; then
  curl -X POST "http://172.18.23.50/Resource/resource/downloadFile?fileUid=rpmtar8" -H "accept: application/json;charset=UTF-8">>rpm.tar.gz
  tar -zxvf rpm.tar.gz >>$localpath/install.log 2>&1
  cd rpm-8
  rpm -Uvh  *.rpm  --nodeps  --force >>$localpath/install.log 2>&1
  cd ..
  rm -r -f rpm-8 
  rm /usr/bin/python 2>>$localpath/install.log
  ln -s /usr/bin/python3.6 /usr/bin/python
  mkdir -p /usr/local/lib64/python3.6/site-packages/
fi
rm -r -f rpm.tar.gz
echo "解压agent主体程序"
unzip agent-go-linux.zip >>$localpath/install.log 2>&1
rm -r -f agent-go-linux.zip
cp log_erron.txt agent-go-linux/logs 2>>$localpath/install.log
cp log_info.txt agent-go-linux/logs 2>>$localpath/install.log
rm -r -f log_erron.txt
rm -r -f log_info.txt
cd agent-go-linux
path=`pwd`
mkdir $path/configure
touch $path/configure/agent.conf
chmod 777 $path/bin/main
chmod 777 $path/configure/agent.conf
chmod 777 $path/daemon/daemon
res=`python -V 2>&1 | grep "Python"`
if [ $? -eq 0 ]; then
  echo "python已安装，版本为：$res"
else
  echo "python未安装"
  exit 1
fi

echo "开始安装所需的第三方包：psutil" >>$localpath/install.log
tar -zxvf $path/environmentInstall/python/psutil-5.7.2.tar.gz >>$localpath/install.log 2>&1
cd psutil-5.7.2
python setup.py build >>$localpath/install.log 2>&1
python setup.py install >>$localpath/install.log 2>&1
cd ..
rm -rf psutil-5.7.2
#echo "安装完毕!" >>$localpath/install.log

echo "# 写配置文件" >>$localpath/install.log
echo "当前路径是：$path" >>$localpath/install.log
echo "# agent版本" >> $path/configure/agent.conf
echo AgentVersion=2.1.2 >> $path/configure/agent.conf
echo "# 是否打印调试信息 true | false" >> $path/configure/agent.conf
echo Debug=true >> $path/configure/agent.conf
echo "# 设置日志路径" >> $path/configure/agent.conf
echo LogPath= >> $path/configure/agent.conf
echo "# httpServer的端口号" >> $path/configure/agent.conf
echo "# HttpServerUrl（默认路径）http://127.0.0.1" >> $path/configure/agent.conf
echo HttpServerPort=9090 >> $path/configure/agent.conf
echo "# 是否批量运行脚本 true | false" >> $path/configure/agent.conf
echo RunScripts=true >> $path/configure/agent.conf
echo "# 批量运行的脚本路径" >> $path/configure/agent.conf
echo RunScriptsSrc=http://172.18.23.50/TopMonitor/osMonitorItem/open-api/fetchTask >> $path/configure/agent.conf
echo RunScriptsDst=http://172.18.23.50/TopMonitor/osMonitorItem/open-api/excResult >> $path/configure/agent.conf
echo "# 执行脚本：获取脚本接口、 返回结果的接口、 执行间隔时间" >> $path/configure/agent.conf
echo "# 间隔时间单位为s(秒)、m(分钟)、h(小时)，只允许小写" >> $path/configure/agent.conf
echo GetInfoScriptSrc=http://172.18.23.50/TopCMDB/autoDiscover/getSystemInfoScript?orgId=18 >> $path/configure/agent.conf
echo GetInfoScriptDst=http://172.18.23.50/TopCMDB/autoDiscover/pushSystemInfo?orgId=e12b254e8abc43ba991c57ec9c1ea793>> $path/configure/agent.conf
echo RunScriptIntervalTime=30m >> $path/configure/agent.conf
echo GetInfoScriptIntervalTime=30m >> $path/configure/agent.conf
echo UninstallDst=http://172.18.23.50/UAC/agentList/reciveAgentUninstall >> $path/configure/agent.conf

tar -xvf $path/environmentInstall/python/sysstat-12.3.3.tar.gz >>$localpath/install.log 2>&1
cd sysstat-12.3.3/
./configure >>$localpath/install.log 2>&1
make >>$localpath/install.log 2>&1
make install >>$localpath/install.log 2>&1
cd ..
rm -r -f sysstat-12.3.3
#rm -r -f environmentInstall
echo "# 运行agent" >>$localpath/install.log
nohup $path/daemon/daemon >/dev/null 2>&1 &
nohup $path/bin/main $path/configure/agent.conf >/dev/null 2>&1 &
res=`rpm -q centos-release | grep "el8"`
if [ $? -ne 0 ]; then
  echo "将agent设为开机自启" >>$localpath/install.log
  touch $path/agent-go-linux
  chmod +x agent-go-linux
  echo "#add for chkconfig" >> $path/agent-go-linux
  echo "#chkconfig: 2345 70 30" >> $path/agent-go-linux
  echo "#description: agent-go-linux" >> $path/agent-go-linux
  echo "#processname: agent-daemon " >> $path/agent-go-linux
  echo "nohup $path/daemon/daemon &" >> $path/agent-go-linux
  echo "nohup $path/bin/main $path/configure/agent.conf &" >> $path/agent-go-linux
  mv $path/agent-go-linux /etc/init.d/agent-go-linux
  chkconfig --add agent-go-linux
  chmod +x /etc/rc.d/rc.local
  res=`cat /etc/rc.d/rc.local | grep /etc/init.d/agent-go-linux`
  if [ $? -ne 0 ]; then
    echo "/etc/init.d/agent-go-linux" >> /etc/rc.d/rc.local
  fi
fi
res=`rpm -q centos-release | grep "el8"`
if [ $? -ne 0 ]; then
  echo "将agent设为开机自启" >>$localpath/install.log
  touch $path/agent-go-linux
  chmod +x agent-go-linux
  echo "#add for chkconfig" >> $path/agent-go-linux
  echo "#chkconfig: 2345 70 30" >> $path/agent-go-linux
  echo "#description: agent-go-linux" >> $path/agent-go-linux
  echo "#processname: agent-daemon " >> $path/agent-go-linux
  echo "nohup $path/daemon/daemon &" >> $path/agent-go-linux
  echo "nohup $path/bin/main $path/configure/agent.conf &" >> $path/agent-go-linux
  mv $path/agent-go-linux /etc/init.d/agent-go-linux
  chkconfig --add agent-go-linux
  chmod +x /etc/rc.d/rc.local
  res=`cat /etc/rc.d/rc.local | grep /etc/init.d/agent-go-linux`
  if [ $? -ne 0 ]; then
    echo "/etc/init.d/agent-go-linux" >> /etc/rc.d/rc.local
  fi
fi
echo "防止运行命令被下载时重复的字符影响" >>$localpath/install.log
echo "agent installed!"
# rm -r -f install.sh
