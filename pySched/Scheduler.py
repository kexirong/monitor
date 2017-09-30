import os
import time
import Queue
import threading
import select
import socket
import json

VAL_QUEUE=Queue.Queue()
PATH="./pythscript/"


PLUGINMAP={}

class Plugin(object):
    __slots__ = {"name",'step','instan'}
    def __init__(self):
        self.name = None
        self.step = None
        self.instan = None


    

def instance_add(plugin):
    if not isinstance(plugin,Plugin):
        return False
    if not isinstance(plugin.step,int) and plugin.step <= 0 :
        return  False
    if not isinstance(plugin.name,str) and plugin.name  in  [None,'']:
        return  False
    if not hasattr(plugin.instan,"getvalue"):
        return  False
    if PLUGINMAP.get(plugin.step):
        PLUGINMAP[plugin.step].append(plugin)
    else:
        PLUGINMAP[plugin.step]=[plugin]
    return True
        
 
def plugin_run(instan):
    try:
        value=getattr(instan.instan,"getvalue")()
        VAL_QUEUE.put(value)
        print(value)
    except AttributeError,e:
        print(dir(instan))
        
    except Exception,e:
        print(instan.name,e)
       
 
 
class Cron(object):
    _instance = None
    
    def __new__(cls, *args, **kw):
        if not cls._instance:
            cls._instance = super(Cron, cls).__new__(cls, *args, **kw)  
        return cls._instance  
        
    def __init__(self):
        self.timerec={k:0 for k in PLUGINMAP}
    
    def cron(self):
        while True:
           # timesec=int(time.time())
            
            for i in PLUGINMAP: # key is int type,so...
                timewant=int(time.time())-self.timerec[i]
                print("############:",i,timewant)
                
                if timewant < i:
                    print("sleep...")
                    time.sleep(i-timewant)
                    
                self.timerec[i]=int(time.time())
                for j in PLUGINMAP[i]:
                    t =threading.Thread(target=plugin_run,args=(j,))
                    t.setDaemon(1)
                    t.start()
                    
    def up_timerec(self):
        if set(PLUGINMAP)==set(self.timerec):
            return
        comp=set(PLUGINMAP)-set(self.timerec)
        for i in comp:
            self.timerec[i]=0
        
    

'''
import socket
import os
import time

sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM )
path = './agent.sock'
if os.path.exists(path):
    os.unlink(path)


sock.bind(path)
sock.listen(1)
conn,addr = sock.accept()
data=""
while True:
     msg = conn.recv(1)
     if msg  != "\n" :
         data+=msg
     else:   
       print(data)
       data=''


'''
            
            
class  AFUNIX_TCP(object):
    def __init__(self):
        self.sock=socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
        self.sock.setblocking(0)
        self.epoll=select.epoll()
    
    def conn(self,path='./agent.sock'):
        while True:
            if os.path.exists(path):
                try:
                    self.sock.connect(path)
                except socket.error:
                    self.__init__()
                    continue
                break
            else:
                print("waiting server...")
                time.sleep(5)

    def send(self,msg):

        self.sock.send(msg)
        self.sock.send("\n")
         
    def recv(self):
        return self.sock.recv(1024)
        
        
    def colse(self):       
        self.sock.close()

   
    def transfer(self):
        
        while True:
            self.conn()
            self.epoll.register(self.sock.fileno(),select.EPOLLIN|select.EPOLLOUT)
  
            while True:
                events = self.epoll.poll(15)

                if not events :
                    break

                print("events:",events)
                time.sleep(2)
                for fileno, event in events:
                    if event&select.EPOLLIN:
                        rec=self.recv()
                        do(parse(rec))
                        
                    elif event&select.EPOLLOUT:
                        
                        while not VAL_QUEUE.empty():
                            msg=VAL_QUEUE.get()
                            self.send(msg)

                if event==select.EPOLLHUP :
                    self.epoll.unregister(fileno)
                    self.close()
                    break
                        
                
       

def do(cmd):
    if cmd is None:
        return
    if cmd["op"]=="addplugin":
        if not cmd.get('plugins'):
            return
            
        for i in cmd['plugins']:
            plugin=loadplugin(i)
            if plugin:
                instance_add(plugin)
        Cron().up_timerec()   
        
    return 'done' 
        
        
    
def parse(recv):
    try:
        cmd=json.loads(recv)
    except ValueError:
        print(recv)
        cmd=None
    return cmd
    
def loadplugin(name):
    try:
        plugin=Plugin()
        plugin.instan=__import__(name)
        plugin.name=plugin.instan.name
        plugin.step=plugin.instan.step
    except Exception,e:
        print(e)
        plugin=None
        
    return plugin
    
    
    

def mian():
    dirlist=os.listdir(PATH)
    for i in dirlist:
        if not i.endswith(".py"):
            continue
        plugin=loadplugin(i[0:-3])
        if plugin:
            instance_add(plugin)
   
    print('start AFUNIX_TCP client ')  
    t=threading.Thread(target=AFUNIX_TCP().transfer) 
    t.setDaemon(1)
    t.start()
    print('run cron... ')
    
    Cron().cron()
    
        

if __name__ == '__main__':
    mian()