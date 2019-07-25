bige
=======
这是一个游戏服务器的基础框架。（golang）

介绍
--------
* 它的特点就是把游戏服务器中需要的每个功能点都分成一个个原子，让你在实际使用的时候，按需求对其进行组合。

QQ群号：441240897
--------
文档
--------
* [HTTPMoudle](https://github.com/buguang01/bige/blob/master/module/README.md)
* [LogicMoudle](https://github.com/buguang01/bige/blob/master/module/README.md)
* [SqlDataModule](https://github.com/buguang01/bige/blob/master/module/README.md)
* [WebSocketModule](https://github.com/buguang01/bige/blob/master/module/README.md)
* [MemoryModule](https://github.com/buguang01/bige/blob/master/module/README.md)
* [NsqdModule](https://github.com/buguang01/bige/blob/master/module/README.md)

使用的第三方库
--------
* mysql : go get -u github.com/go-sql-driver/mysql
* redis : go get -u github.com/garyburd/redigo
* 打印颜色：go get -u github.com/gookit/color
* Nsq   : go get -u github.com/nsqio/go-nsq


开发进度
--------
* 已完成的子功能：
*   NsqdModule      与Nsq中间件通信的模块
*   HTTPModule      HTTP的收消息模块
*   WebSocket       WebSocket收发消息模块
*   LogicModule     业务逻辑模块，用来管理业务协程，可以让业务逻辑在指定KEY的协程上运行
*   SqlDataModule   数据库处理模块，可以让DB操作在指定KEY的协程上运行，还可以设置延时运行
*   MemoryModule   内存数据管理器，可以用来管理，数据什么空闲多少时间后，进行卸载
*   event           收到的消息基础类型、module用到的一些信道数据结构
*   model           mysql的模块、Redis的模块
*   threads         协程管理

