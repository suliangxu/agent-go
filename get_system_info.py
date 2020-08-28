import os
import re
import platform
import psutil
import sys
import socket

version = "1.0.0.1"

class linux_info(object):
    def __init__(self):
        self.list = {}

    def getData(self, param, data, value=2):
        data = re.findall(r"{0}.*".format(param), data)[0].replace("  ", " ")
        data = data.replace("\t", " ").split(" ")
        return data[len(data) - value]

    def getResult(self, cpuNumber, cpuType, loginUserName,
                  disk, hostIP, userName,
                  hostName, macAddress,
                  memory, networkCardInfo, netmask,
                  osArch, osName, osVersion, platformType, userIp):
        data = {}
        data["clientType"] = 'pc'
        data["cpuNumber"] = cpuNumber
        data["cpuType"] = cpuType
        data["entriedCode"] = '123456'
        data["disk"] = disk
        data["hostIP"] = hostIP
        data["hostName"] = hostName
        data["loginUserName"] = loginUserName
        data["macAddress"] = macAddress
        data["memory"] = memory
        data["networkCardInfo"] = networkCardInfo
        data["netmask"] = netmask
        data["osArch"] = osArch
        data["osName"] = osName
        data["osVersion"] = osVersion
        data["platformType"] = platformType
        data["userIp"] = userIp
        data["userName"] = userName

        self.list = data

    def exce(self):
        ## cpu
        cpu_data = open("/proc/cpuinfo", "r").read()

        cpuNumber = len(re.findall(r"processor", cpu_data))

        cpuType = re.findall(r"model name.*", cpu_data)[0].split("\t")[1].replace(':', '').strip()

        # loginUserName
        loginUserName = psutil.users()[0][0]

        # userip
        userip = psutil.net_if_addrs()['lo'][0][1]

        # userName
        userName = socket.gethostname()

        disk_list = []
        a = os.popen('lsblk -P').readlines()
        for i in a:
            if('''TYPE="disk"''' in i):
                n = re.match(r'''.*SIZE="(.*)" RO.*\n''', i, re.M | re.I).group(1)
                disk_list.append(n)

        disk = 0
        for j in disk_list:
            if (j[-1] == 'G'):
                disk += float(j[0:-1])
            elif (j[-1] == 'M'):
                disk += (float(j[0:-1]) / 1024.0)
            elif (j[-1] == 'K'):
                pass

        if (disk - int(disk) >= 0.5):
            disk = int(disk) + 1;
        else:
            disk = int(disk)

        # hostname
        hostName = platform.node()

        # hostIP
        hostIP = os.popen("ip addr|grep -E  'eth|en'|grep inet|awk '{print $2}'|head -n 1").read().replace('/24','').replace('\n', '')

        # version
        osVersion = platform.release()

        # platformType
        platformType = platform.system()

        # osname
        osName = platform.dist()[0]

        # osVersion
        osVersion = platform.linux_distribution()[1]

        # osArch
        osArch = platform.architecture()[0]

        ## MEMORY
        mem = psutil.virtual_memory()
        memory = mem.total / (1024*1024*1024*1.0)
        if (memory - int(memory) >= 0.5):
            memory = int(memory) + 1;
        else:
            memory = int(memory)

        # macAddress
        macAddress = os.popen("ip addr|grep -E  'eth|en'|grep 'link/ether'|awk '{print $2}'|head -n 1").read().replace(
            '\n', '')

        # networkCardInfo

        address_list = os.popen(
            "ip addr|grep -E  'eth|en'|grep inet|awk '{print $2}'|awk -F\/ '{print $1}'").readlines()
        mac_list = os.popen("ip addr|grep -E  'eth|en'|grep 'link/ether'|awk '{print $2}'").readlines()
        gateway_list = os.popen("ip addr|grep -E  'eth|en'|grep inet|awk '{print $4}'").readlines()
        network_name = os.popen("ip addr|grep -E 'en|eth'|grep -vE 'inet|lo|link'|awk -F: '{print $2}'").readlines()

        # netmask
	netmask = None
        net = psutil.net_if_addrs()
        for i in net.keys():
            if ('en' in i):
                netmask = net[i][0][2]
                if (netmask == None):
                    continue
                else:
                    break


        networkCardInfo = ''
        network = []

        for num in range(len(address_list)):
            network_list = []
            network_list.append(network_name[num].replace('\n', '').strip())
            network_list.append(address_list[num].replace('\n', ''))
            network_list.append(gateway_list[num].replace('\n', ''))
            network_list.append(mac_list[num].replace('\n', ''))
            network.append(network_list)
            networkCardInfo = str(network)

        self.getResult(cpuNumber,
                       cpuType, loginUserName,
                       disk, hostIP, userName,
                       hostName,
                       macAddress, memory, networkCardInfo, netmask, osArch, osName, osVersion, platformType, userip)


    def Print(self):
        # print(self.list)
        for key, value in self.list.items():
            print("%s: %s" % (key, value))




class win_info(object):

    def __init__(self):
        self.data = {}
        self.list = []

    def get_code(self, str_name):

        return str_name.decode('gbk').encode('utf-8')

    def getResult(self, cpuNumber, cpuType,
                  disk, hostIP,
                  hostName, macAddress,
                  memory, networkCardInfo,
                  osArch, osName, osVersion, platformType, loginUserName, userIp, userName):
        data = {}
        data["clientType"] = 'pc'
        data["cpuNumber"] = cpuNumber
        data["cpuType"] = cpuType
        data["entriedCode"] = '123456'
        data["disk"] = disk
        data["hostIP"] = hostIP
        data["hostName"] = hostName
        data["loginUserName"] = loginUserName
        data["macAddress"] = macAddress
        data["memory"] = memory
        data["networkCardInfo"] = networkCardInfo
        data["osArch"] = osArch
        data["osName"] = osName
        data["osVersion"] = osVersion
        data["platformType"] = platformType
        data["userName"] = userName
        data["userIp"] = userIp

        self.list.append(data)

    def replace_str(self, str_name):
        if isinstance(str_name, int):
            new_str = str_name
        elif isinstance(str_name, str):
            new_str = int(str_name)
        else:

            new_str = str_name.replace('\r', '').replace('\n', '').replace(' ', '').replace('{', '') \
                .replace('}', '')
        return new_str

    def main(self):
        # cpu

        cpu_process = psutil.cpu_count()

        cpuNumber = self.replace_str(cpu_process)
        cpuType = wmi.WMI().Win32_Processor()[0].Name

        # disk
        a = os.popen('wmic DiskDrive get Size /value').readlines()
        disk = 0

        for i in a:
            if ('Size=' in i):
                n = i.replace('Size=', '').strip('\n')
                disk += int(n)

        disk /= (1024 * 1024 * 1024 * 1.0)
        if (disk - int(disk) >= 0.5):
            disk = int(disk) + 1;
        else:
            disk = int(disk)

        # memory

        mem = psutil.virtual_memory()
        memory = mem.total / (1024*1024*1024*1.0)
        if (memory - int(memory) >= 0.5):
            memory = int(memory) + 1;
        else:
            memory = int(memory)


        # hostname
        hostName = socket.gethostname()

        # hostIP
        hostIP = socket.gethostbyname(hostName)

        # loginUserName
        loginUserName = psutil.users()[0][0]

        # userName
        userName = psutil.users()[0][1]

        # userip
        user = psutil.net_if_addrs()
        for i in user.keys():
            if("Loopback" in i):
                continue
            userip = user[i][1][1]


        # hostinfo
        net_data = os.popen('wmic nicconfig where \"IPEnabled=\'TRUE\'\" get ipaddress,macaddress').read()
        c = wmi.WMI()
        for interface in c.Win32_NetworkAdapterConfiguration(IPEnabled=1):
            macAddress = interface.MACAddress

            for ip in interface.IPAddress:

                if re.match(r"^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$", ip):
                    hostIP = ip

            break
        network_info_list = []
        network_info = []
        num = 0
        for interface in c.Win32_NetworkAdapterConfiguration(IPEnabled=1):
            network_info = []
            network_info.append('net{0}'.format(num))
            macAddress = interface.MACAddress

            for ip in interface.IPAddress:

                if re.match(r"^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$", ip):
                    network_info.append(ip)
            for sub in interface.IPSubnet:
                if re.match(r"^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$", sub):
                    network_info.append(sub)

            num += 1

            network_info.append('')
            network_info.append(macAddress)

            network_info_list.append(network_info)

        osArch = platform.machine()
        osName = platform.release()
        osVersion = platform.version()
        platformType = platform.system()

        hostName = platform.node()
        networkCardInfo = network_info_list
        self.getResult(cpuNumber, cpuType, disk, hostIP, hostName, macAddress, memory, networkCardInfo, osArch, osName,
                       osVersion, platformType, loginUserName, userip, userName)


    def Print(self):
        # print(self.list)
        for key, value in self.list[0].items():
            print("%s: %s" % (key, value))


if ("win" in sys.platform):
    import wmi
    mywin_info = win_info()
    mywin_info.main()
    mywin_info.Print()

else:
    mylinux_info = linux_info()
    mylinux_info.exce()
    mylinux_info.Print()