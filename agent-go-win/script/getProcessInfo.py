import psutil
import operator
import copy
import sys
import json

# get process
# pid = int(sys.argv[1])
result = []

def get_info(pid):
    try:
        result_map = {}
        p = psutil.Process(pid)
        p_name = p.name()
        t = p.username()
	t2 = t.rfind('\\')
	p_username = t[t2+1:]
        cpu_percent = p.cpu_percent()
        mem_percent = p.memory_percent()

        result_map["process"] = p_name
        result_map["user"] = p_username
        result_map["cpuUsage"] = cpu_percent
        result_map["memoryUsage"] = mem_percent

        result.append(result_map)
    except:
        pass


def process_info(orderType, orderDesc, size):
    result2 = copy.deepcopy(result)

    if orderType==1 and orderDesc==1 :
        list2 = sorted(result2, key=operator.itemgetter('cpuUsage'))
    elif orderType==1 and orderDesc==2 :
        list2 = sorted(result2, key=operator.itemgetter('cpuUsage'), reverse=True)
    elif orderType == 2 and orderDesc == 1:
        list2 = sorted(result2, key=operator.itemgetter('memoryUsage'))
    elif orderType == 2 and orderDesc == 2:
        list2 = sorted(result2, key=operator.itemgetter('memoryUsage'), reverse=True)

    for i in list2:
        i["cpuUsage"] = "{:.2f}".format(i["cpuUsage"])
        i["memoryUsage"] = "{:.2f}".format(i["memoryUsage"])
    str_dic = json.dumps(list2[:size])
    print(str_dic)


ot = int(sys.argv[1])
od = int(sys.argv[2])
size = int(sys.argv[3])

for i in psutil.pids():
    get_info(i)

process_info(ot, od, size)