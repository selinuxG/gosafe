#!/bin/bash
sleep_time=${1:-'60'} #���ǵ�һ�����յ��Ĳ��������ڶ�ʱ����쳣��ʱ�䣬Ĭ��60�롣
sleep_ping=`expr ${sleep_time} \* 3 ` #���ڼ��Υ������ʱ��,��sleep_time��3����
#----------------------------------------------------------Ϊ�˽����г��򱣳�������ֹ֮ͣǰ����
#kill -9 $(ps aux |grep gosafe |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /var/log/nginx/access.log |awk -F " " '{print $2}')
#kill -9 $(ps aux |grep gocli |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /dhcp/logs |awk -F " " '{print $2}')
kill -9 $(ps aux |grep /var/log/nginx/access.log |awk -F " " '{print $2}')
#----------------����Ŀ¼
mkdir -p /dhcp/logs
touch /dhcp/logs/warnmd5.log
touch /dhcp/history.log
touch /dhcp/logs/ping.log
touch /dhcp/logs/warnnginx.log
touch /dhcp/logs/warnhistory.log
touch /dhcp/tailnginx.log
#-----------------�жϼ����ʷ���������
histenv=`cat /etc/bashrc |grep PROMPT_COMMAND|grep 'read x y'|grep log|wc -l`
if [ $histenv == 0 ]; then
    #echo "IS NULL"
    printf "%s%s%s%s" 'export PROMPT_COMMAND=' "'" 'history -a ; command=$(history 1 | { read x y;echo $y >> /dhcp/history.log; } )' "'" >>/etc/bashrc;
    source /etc/bashrc
else
    echo "NOT NULL"
fi
#--------------------------------------------------------------д��Ŀ¼MD5ֵ
find /etc -type f -exec md5sum {} >>/dhcp/oldetc.md5 \;
#----------------��nginx������־ʵʱ����:����ȴ�1������޷���ջ�����10����ǰ������
nohup tail -f /var/log/nginx/access.log >> /dhcp/tailnginx.log &
sleep 1
> /dhcp/tailnginx.log
#---------------------------------------------------------------����client����
nohup /dhcp/gosafe -run tail -serverip http://103.44.250.69:17171 -path /var/log/secure,/dhcp/logs > /dev/null    &
#---------------------------------------------------------------�����������쳣
nowping_count=0
while true
do
    let nowping_count+=${sleep_time} #ping����+sleep_time
    sleep ${sleep_time}
    md5sum -c /dhcp/oldetc.md5 | grep -vE "ȷ��|OK" >/dhcp/logs/warnmd5.log  #�Ա�MD5ֵ�Ƿ����쳣,�쳣����д�뵽/dhcp/history.log��
    cat /dhcp/history.log >> /dhcp/cmd.log  #���ݵ�ǰ��ʷ���cmd.log��ȥ
    #-------------------------------------------------------------------------
    /dhcp/gosafe -run history -cmdfile /dhcp/history.log -warnfile /dhcp/logs/warnhistory.log   #ͨ��gosafeģʽ���/dhcp/history.log�Ƿ�����쳣����,д�뵽/dhcp/logs/warnhistory.log
    > /dhcp/history.log  #�Ա����Ժ������������´λ��ظ�����
    #--------------------------------------#���webӦ���쳣��¼��/dhcp/logs/warnnginx.log��־��
    cat /dhcp/tailnginx.log |awk -F " " '$9==404 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","��ԴIP��ַ:",$2,"״̬��:",$7,"�������:",$1,"����URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log ||awk -F " " '$9==403 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","��ԴIP��ַ:",$2,"״̬��:",$7,"�������:",$1,"����URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log |grep -E "php|jsp|jspx|phpx|script|alert|onclack" | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  | awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","��ԴIP��ַ:",$2,"״̬��:",$7,"�������:",$1,"����URL:",$4," ",$5}' >> /dhcp/logs/warnnginx.log
    cat /dhcp/tailnginx.log ||awk -F " " '$9==302 {print $0}' | awk -F " " '{print $1" --- "$6" "$7 " --- "$9}'| sort | uniq -c  |awk -F " " '{printf "%s%-16s%s%-4s%s%-3s%s%s%s%s\n","��ԴIP��ַ:",$2,"״̬��:",$7,"�������:",$1,"����URL:",$4," ",$5}' >>/dhcp/logs/warnnginx.log
    > /dhcp/tailnginx.log   #�Ա����Ժ������������´λ��ظ�����
    #--------------------------------------------------------#����Ƿ������ͨ�������Ƶ��ʱ�����������ܵ�3��
    if [[ $nowping_count -eq $sleep_ping ]];then
      #nowping_count=0
      let nowping_count=0
      ping www.baidu.com -c 1 -W 1 >/dev/null
      if [ $? -eq 0 ];then
        printf "%s\n" "ͨ������ping www.baidu.com ���Ϊ��������ͨ��!" >> /dhcp/logs/ping.log
      fi
    fi
done

