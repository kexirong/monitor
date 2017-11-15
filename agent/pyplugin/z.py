import psutil
import time


time.sleep(1)

cpu_value = psutil.cpu_times_percent()
print(cpu_value)
time.sleep(1)
cpu_value = psutil.cpu_times_percent()
print(cpu_value)
time.sleep(1)
cpu_value = psutil.cpu_times_percent()
print(cpu_value)

import threading  
print(time.time())
event = threading.Event()
event.wait(1)
print(time.time())






