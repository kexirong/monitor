#!/bin/env python
import os
import time
import Queue
import threading
import select
import sys
import socket
import json

VAL_QUEUE=Queue.Queue()
PATH= "./pysched/pythscript"




PLUGINMAP={}

F=open (r'SchecdOut.log','w')

print >>F ,PATH

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
        print >>F, value
    except AttributeError,e:



        print >>F, dir(instan)
        
    except Exception,e:
        print >>F, instan.name,e
       
 
 
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
                print >>F, "############:",i,timewant
                
                if timewant < i:
                    print >> F ,"sleep..."
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
        
    


            
            
class  AFUNIX_TCP(object):
    def __init__(self):
        self.sock=socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        self.sock.setsockopt(socket.SOL_SOCKET, socket.SO_KEEPALIVE, 1)
        self.sock.setblocking(0)
        self.epoll=select.epoll()
    
    def conn(self,path='pysched/agent.sock'):
        while True:
            if os.path.exists(path):
                try:
                    self.sock.connect(path)
                except socket.error:
                    self.__init__()
                    continue
                break
            else:
                print >> F,"waiting server..."
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

                print >>F,"events:",events

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
        print >>F,recv
        cmd=None
    return cmd
    
def loadplugin(name):
    try:
        plugin=Plugin()
        
        plugin.instan=__import__(name)
        plugin.name=plugin.instan.name
        plugin.step=plugin.instan.step
    except Exception,e:
        print >>F,e
        plugin=None
        
    return plugin
    

def mian():
    sys.path.append(PATH) 
    dirlist=os.listdir(PATH)
    print >>F,"cur....dir",sys.path,dirlist
    for i in dirlist:
        if not i.endswith(".py"):
            continue
        plugin=loadplugin(i[0:-3])
        if plugin:
            instance_add(plugin)
   
    print >>F, 'start AFUNIX_TCP client '
    t=threading.Thread(target=AFUNIX_TCP().transfer) 
    t.setDaemon(1)
    t.start()
    print >>F ,'run cron... '
    
    Cron().cron()
    
        

if __name__ == '__main__':
    mian()