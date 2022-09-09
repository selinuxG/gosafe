#!/bin/bash
sleep_time=${1:-'60'} #这是第一个接收到的参数，用于定时检测异常的时间，默认60秒。
sleep_ping=`expr ${sleep_time} \* 3 ` #用于检测违规外联时常,是sleep_time的3倍。
#----------------------------------------------------------为了将所有程序保持最新先停止之前程序
#kill -9 $(ps aux |grep gosafe |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /var/log/nginx/access.log |awk -F " " '{print $2}')
#kill -9 $(ps aux |grep gocli |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /dhcp/logs |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /var/log/nginx/access.log |awk -F " " '{print $2}')
#----------------创建目录
mkdir -p /dhcp/logs
touch /dhcp/logs/warnmd5.log
touch /dhcp/history.log
touch /dhcp/logs/ping.log
touch /dhcp/logs/warnnginx.log
touch /dhcp/logs/warnhistory.log
touch /dhcp/tailnginx.log
#-----------------判断监测历史命令环境变量
histenv=`cat /etc/bashrc |grep PROMPT_COMMAND|grep 'read x y'|grep log|wc -l`
if [ $histenv == 0 ]; then
    #echo "IS NULL"
    printf "%s%s%s%s" 'export PROMPT_COMMAND=' "'" 'history -a ; command=$(history 1 | { read x y;echo $y >> /dhcp/history.log; } )' "'" >>/etc/bashrc;
    source /etc/bashrc
else
    echo "NOT NULL"
fi
#--------------------------------------------------------------写入目录MD5值
find /etc -type f -exec md5sum {} >>/dhcp/oldetc.md5 \;
#----------------将nginx访问日志实时备份:必须等待1秒否则无法清空会留存10条以前的数据
nohup tail -f /var/log/nginx/access.log >> /dhcp/tailnginx.log &
sleep 1
> /dhcp/tailnginx.log
#---------------------------------------------------------------启动client程序
nohup /dhcp/gosafe -run tail -serverip http://103.44.250.69:17171 -path /var/log/secure,/dhcp/logs > /dev/null    &
#---------------------------------------------------------------开启监测过滤异常
nowping_count=0
while true
do
    let nowping_count+=${sleep_time} #ping计数+sleep_time
    sleep ${sleep_time}
    md5sum -c /dhcp/oldetc.md5 | grep -vE "确定|OK" >/dhcp/logs/warnmd5.log  #对比MD5值是否有异常,异常内容写入到/dhcp/history.log中
    cat /dhcp/history.log >> /dhcp/cmd.log  #备份当前历史命令到cmd.log中去
    #-------------------------------------------------------------------------
    /dhcp/gosafe -run history -cmdfile /dhcp/history.log -warnfile /dhcp/logs/warnhistory.log   #通过gosafe模式监测/dhcp/history.log是否存在异常命令,写入到/dhcp/logs/warnhistory.log
    > /dhcp/history.log  #对比完以后清空命令否则下次会重复过滤
    #--------------------------------------#检查web应用异常记录到/dhcp/logs/warnnginx.log日志中
    cat /dhcp/tailnginx.log |awk -F " " '$9==404 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","来源IP地址:",$2,"状态码:",$7,"请求次数:",$1,"请求URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log ||awk -F " " '$9==403 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","来源IP地址:",$2,"状态码:",$7,"请求次数:",$1,"请求URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log |grep -E "php|jsp|jspx|phpx|script|alert|onclack" | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","来源IP地址:",$2,"状态码:",$7,"请求次数:",$1,"请求URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log ||awk -F " " '$9==302 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  |awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","来源IP地址:",$2,"状态码:",$7,"请求次数:",$1,"请求URL:",$4," ",$5}' >>/dhcp/logs/warnnginx.log
    > /dhcp/tailnginx.log   #对比完以后清空命令否则下次会重复过滤
    #--------------------------------------------------------#检查是否可以联通外网检测频率时间是其他功能的3倍
    if [[ $nowping_count -eq $sleep_ping ]];then
      #nowping_count=0
      let nowping_count=0
      ping www.baidu.com -c 1 -W 1 >/dev/null
      if [ $? -eq 0 ];then
        printf "%s\n" "通过尝试ping www.baidu.com 结果为可与外网通信!" >> /dhcp/logs/ping.log
      fi
    fi
done

