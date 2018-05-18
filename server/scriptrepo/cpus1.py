#!/usr/bin/env python

#*-* coding:utf-8 *-*

import time
import platform
import json

NAME = "cpus"

interval = 30

def getvalue():
    '''
       return type  must is list,because while  need  peer return
       for exampel  [cpu0,cpu1]
    '''
    hostname = platform.node()
    try:
        import psutil
    except ImportError:
        return ["import psutil Error"]

    cpu_value = psutil.cpu_times_percent(percpu=True)
    #cpu_value = psutil.cpu_times_percent()
    value0 = [cpu_value[0].user, cpu_value[0].nice, cpu_value[0].system, cpu_value[0].idle,]
    value1 = [cpu_value[1].user, cpu_value[1].nice, cpu_value[1].system, cpu_value[1].idle,]

    ret0 = {
        "hostname":hostname,
        "timestamp":time.time(),
        "plugin":"cpus1",
        "instance":"0",
        "type":"percent",
        "value":value0,
        "vltags":"user|nice|system|idle",
        #  "Message":"",
        }
    ret1 = {
        "hostname":hostname,
        "timestamp":time.time(),
        "plugin":"cpus1",
        "instance":"1",
        "type":"percent",
        "value":value1,
        "vltags":"user|nice|system|idle",
        #  "Message":"",
        }

    return json.dumps([ret0, ret1])


if __name__ == "__main__":

        print( getvalue())
