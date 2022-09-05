# 工具由来
​	公司安排护网一周，面临的情况是集中检测200+服务器记录，客户现场不具备日志审计等安全设备，只能手撸检测功能并通过Golin工具集中部署、运行。

​	其中共检测：系统登陆记录、异常执行命令记录、重要文件MD5对比记录、与外网互通记录、WEB出现敏感状态码以及敏感请求记录。

# 使用方法

## 使用参数

-run server -port 17171		//此模式是运行server模式，指定端口为17171端口运行

-run tail -path -path /var/log/secure,/etc/passwd		//此模式是client模式，指定检测的目录或文件

-run history -cmdfile history.log -warnfile warnhistory.log		//此参数是基于历史命令文件监测是否存在可疑命令写入到新文件中在基于tail模式发生给server

-run info -apiserver http://103.44.250.69:17171/api?page=1&pageSize=100000		//此模式是获取API数据写入到xlsx

## 服务端

直接运行gosafe(默认运行为server端，端口为17171)，可通过 -port 指定启动端口。

## 客户端

​	wget -P /root/ [http://103.44.250.69:11111/gocli.sh;chmod](http://103.44.250.69:11111/gocli.sh;chmod) +x /root/gocli.sh;nohup sh /root/gocli.sh  10  > gocli.log 2>&1 & (其中需要调整server地址)

# 记录方式

1. 会输出文本格式的报警记录。
1. 提供API报警接口查询，未避免数据量太庞大，接口调整为分页。
1. 可通过运行info模式将报警记录记录到xlsx文件中，方便过滤生成图表等。

