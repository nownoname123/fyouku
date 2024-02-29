# fyouku文件
里面放的是项目的前端代码，通过beego框架+html+css+javascript来完成网站页面的布局以及跳转的路由
执行方法：需要安装bee框架，然后执行
```
bee run
```
需要安装beego框架
```
go get github.com/astaxie/beego
```
# fyoukuapi文件
里面存放的是项目的后端代码，技术栈有：
* 用redis缓存热门视频数据以及前几条评论，减少数据库压力，用redis的zset数据结构来对视频按照评论数进行排行
* 由于热门视频排行榜是依靠评论数来决定的，在评论后，通过rabbitMQ来交给别的程序执行更新视频的评论总数以及排名操作，实现了服务的解耦，rabbitMQ部署在远程云服务器：175.178.212.4中
* 用WebSocket 来完成客户端与服务器的双向通信，实现弹幕功能
* 利用ElasticSearch以及ElasticSearch的head插件和ik分词器实现关键字搜索的操作，需要安装对应的插件并且运行
  **安装elasticsearch**
   ```
   https://www.elastic.co/downloads/past-releases
   ```
 **安装head管理工具**
  ```
  https://github.com/mobz/elasticsearch-head
  ```
**安装ik分词器**
```
https://github.com/medcl/elasticsearch-analysis-ik
```
