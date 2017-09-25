## 关于beekeeper

一个Go语言实现的轻量级FCGI框架。

##Go版本要求

go1.7及以上。

## 主要模块

* conf
	*  统一的配置加载模块，可以读取命令行参数与XML配置文件，解析为对象。
	
* grace
	*  提供不中断服务功能，在重启进程期间不会断开已有连接，直到这些连接超时。
	
* log
	*  基于标准库log包实现的rotate日志模块。

* mon
	*  基于标准库expvar包实现的监控接口模块，可以将程序运行期间的一些统计变量（例如内存使用情况、请求数量以及请求耗时等）通过HTTP协议以统一的JSON格式暴露给外部。

* router
	*  一个简单的请求路由模块，根据注册的路由信息将不同路径和方法的请求路由到相应的函数处理。

* db
	*  一个数据库查询helper模块，能够方便地查询单行或多行数据并转化为对象、slice或map。
	
* mailmsg
	*  用于创建带有内嵌资源或附件的MIME格式邮件。

## 如何创建自己的代码

在Go环境已经正确安装和配置的前提下：

1. 将beekeeper框架复制到$GOPATH/src目录下。

2. 运行beekeeper/create.bat（create.sh）脚本，在$GOPATH/src目录下生成目录。

3. 进入生成的目录，使用go build命令编译。

## 配置web服务器

在nginx配置中添加如下转发规则：

```nginx
location /your_program_name/ {
	fastcgi_pass x.x.x.x:yyyy;
	fastcgi_index index.cgi;
	fastcgi_param SCRIPT_FILENAME fcgi$fastcgi_script_name;
	include fastcgi_params;
}
```

注意转发的IP地址和端口应当与配置文件中的IP地址和端口保持一致。
