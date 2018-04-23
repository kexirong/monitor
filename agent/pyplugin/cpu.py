#*-* coding:utf-8 *-*
import  time
import  platform
import json


NAME = "cpu"

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
    
    cpu_value = psutil.cpu_times_percent()
    value = [cpu_value.user, cpu_value.nice, cpu_value.system, cpu_value.idle,]

    ret = {
        "hostname":hostname,
        "timestamp":time.time(),
        "plugin":"cpu",
        "instance":"",
        "type":"percent",
        "value":value,
        "vltags":"user|nice|system|idle",
        #  "Message":"",
        }
    return json.dumps([ret,])

if __name__ == "__main__":
        print( getvalue())
        
