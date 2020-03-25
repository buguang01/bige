bige
=======
这是一个游戏服务器的基础框架。（golang）

介绍
--------
* 它的特点就是把游戏服务器中需要的每个功能点都分成一个个原子，让你在实际使用的时候，按需求对其进行组合。
* 例子工程：https://github.com/buguang01/gsdemo

交流QQ群号：441240897
--------
开发进度
--------
* 现在是2.0版本
* 新的模块都放在了modules
* AutoTask为自动任务模块
* DataBase为DB任务模块
* Logic为逻辑处理模块
* Nsqd为与nsq进行通信的模块
* SocketCli与Socket模块是二个相互通信用的模块一个是服务一个是客户端
* Web是走http的通信模块
* WebSocket与socket模块类似