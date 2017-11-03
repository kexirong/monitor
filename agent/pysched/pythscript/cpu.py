#*-* coding:utf-8 *-*

import time
import json


NAME = "cpu"

STEP = 30

def getvalue():
    '''
       return type  must is list,because while  need  peer return
       for exampel  [cpu0,cpu1]
    '''
    try:
        import psutil
    except ImportError:
        return ["import psutil Error"]
    #cpu_value=psutil.cpu_times_percent(percpu=True)
    cpu_value = psutil.cpu_times_percent()
    value = [cpu_value.user, cpu_value.nice, cpu_value.system, cpu_value.idle,]

    ret = {
        "HostName":"labsr202",
        "TimeStamp":time.time(),
        "Plugin":"cpu",
        "Instance":"",
        "Type":"percent",
        "Value":value,
        "VlTags":"user|nice|system|idle",
        #  "Message":"",
        }

    return [json.dumps(ret),]


if __name__ == "__main__":
    for i in range(100):
        print(i, ":", getvalue())
        time.sleep(1)
