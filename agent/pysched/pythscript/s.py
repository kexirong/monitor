import time,json
name="cpu"
step=5
def getvalue():
    ret={"HostName":"labsr202",
    "TimeStamp":time.time(),
    "Plugin":"cpu",
    "Instance":"0",
    "Type":"percent",
    "Value":[20,80],
    "VlTags":"idle|user",
    "Message":"test"}
    return json.dumps(ret)