gsframe
=======
这是一个游戏服务器的基础框架。（golang）

介绍
--------
* 它的特点就是把游戏服务器中需要的每个功能点都分成一个个原子，让你在实际使用的时候，按需求对其进行组合。

QQ群号：441240897
--------
文档
--------
* [HTTPMoudle](https://github.com/buguang01/gsframe/blob/master/module/README_HTTP.md)
* [LogicMoudle](https://github.com/buguang01/gsframe/blob/master/module/README_LOGIC.md)
* [SqlDataModule](https://github.com/buguang01/gsframe/blob/master/module/README_SQL.md)
* [WebSocket](https://github.com/buguang01/gsframe/blob/master/module/README_WEBSOCKET.md)

使用的第三方库
--------
* mysql : go get -u github.com/go-sql-driver/mysql
* redis : go get -u github.com/garyburd/redigo
* 打印颜色：go get -u github.com/gookit/color

借用库
-------
* utils: go get -u github.com/typa01/go-utils

开发进度
--------
* 已完成的子功能：
*   loglogic        日志管理器
*   HTTPModule      HTTP的收消息模块
*   WebSocket       WebSocket收发消息模块
*   LogicModule     业务逻辑模块，用来管理业务协程，可以让业务逻辑在指定KEY的协程上运行
*   SqlDataModule   数据库处理模块，可以让DB操作在指定KEY的协程上运行，还可以设置延时运行
*   event           收到的消息基础类型、module用到的一些信道数据结构
*   model           mysql的模块、Redis的模块
*   threads         协程管理
*   util            通用基础模块（String、StringBuilder、TimeConvert、(WorkerID)SnowFlakeID、BaseData）

