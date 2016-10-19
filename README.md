##关于beekeeper

一个GO语言实现的轻量级fcgi框架。

##GO版本要求

go1.7及以上。

##主要模块

* conf
	*  统一的配置加载模块，可以读取命令行参数与xml配置文件，解析为对象。
	
* grace
	*  提供不中断服务功能，在重启进程期间不会断开已有连接，直到这些连接超时。
	
* log
	*  基于标准库log包实现的rotate日志模块。

* mon
	*  基于标准库expvar包实现的监控接口模块，可以将程序运行期间的一些统计变量通过http协议以统一的格式暴露给外部。

* router
	*  一个简单的请求路由模块，根据注册的路由信息将不同路径和方法的请求路由到相应的函数处理。

##如何创建自己的代码

在确保正确安装GO环境的前提下，运行beekeeper/create.bat（create.sh）脚本，在$GOPATH/src目录下生成目录。

进入生成的目录，使用go build命令编译。

##配置web服务器

在nginx配置中添加如下转发规则：

```javascript
location /your_program_name/ {
	fastcgi_pass x.x.x.x:yyyy;
	fastcgi_index index.cgi;
	fastcgi_param SCRIPT_FILENAME fcgi$fastcgi_script_name;
	include fastcgi_params;
}
```

注意转发的ip地址和端口应当与配置文件中的ip地址和端口保持一致。