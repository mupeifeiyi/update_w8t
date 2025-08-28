本项目致力于为WatchAlert v3.5.0用户升级v3.6+时，数据库字段有变化，用来刷新告警规则的数据库字段，可在w8t v3.6.0的所有版本中重复执行，不会影响现有数据

# 项目背景
作者给[WatchAlert](https://github.com/opsre/WatchAlert)项目发起人提出需求：
1. 告警规则中，存在多条告警等级时，希望每个告警等级都能有对应的持续时间
2. 值班表中希望单天内可以多人值班

# 项目功能
协助WatchAlert项目更新:
1. 每个告警等级支持对应各自的持续时间（v3.6.0-beat.1）
<img width="1412" height="497" alt="截屏2025-07-17 22 49 57" src="https://github.com/user-attachments/assets/da0e5523-c863-470c-b754-17cbb4cfffff" />
2. 支持同一天多人值班（v3.6.0-beat.2）
<img width="1473" height="624" alt="bfc4160a-0ed5-4bcd-bfb7-f2920017852f" src="https://github.com/user-attachments/assets/5c3015ba-dc75-4b79-8df9-676b18e03bfe" />
3. 阿里云SLS数据源的告警规则支持配置多个logstore（v3.6.0-beat.11）
# 使用方法
以compose为例的升级步骤：
1. 使用本项目刷数据库结构
2. docker compose down
3. 更新compose文件镜像tag
4. docker compose up -d

## 本项目刷数据库
[release](https://github.com/mupeifeiyi/update_w8t/releases)页面查看和WatchAlert对应版本的二进制文件
下载到部署WatchAlert的服务器中
```shell
$ chmod +x u8t-linux-amd64
$ ./u8t-linux-amd64 -h
用法: ./u8t-linux-amd64 --dsn=<dsn字符串>

参数说明：
--dsn      必填，MySQL连接字符串
-h         显示帮助信息

示例：
./u8t-linux-amd64 --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
```
部署时若没有对数据库信息做任何修改，可以直接复制示例进行刷数据库，重复执行不会影响数据
```shell
./u8t-linux-amd64 --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
✅ 正在使用DSN连接数据库: root:****@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local
📣 开始刷metrics告警规则模版数据结构

2025/08/29 00:46:08 /Users/feiyi/update_w8t/utils/utils.go:54
[4.574ms] [rows:62] SELECT * FROM `rule_templates` WHERE datasource_type IN ('prometheus','victoriametrics')
📊 查询到的记录数量： 62
...
✅ 所有metrics告警规则模版更新完成
📣 开始刷metrics告警规则数据结构

2025/08/29 00:46:09 /Users/feiyi/update_w8t/utils/utils.go:16
[1.132ms] [rows:4] SELECT * FROM `alert_rules` WHERE datasource_type IN ('prometheus','victoriametrics')
📊 查询到的记录数量： 4
...
2025/08/29 00:46:09 /Users/feiyi/update_w8t/utils/utils.go:38
[7.252ms] [rows:1] UPDATE `alert_rules` SET `prometheus_config`='{"promQL":"round(max((node_filesystem_size_bytes{fstype=~\"ext.?|xfs\",}-node_filesystem_free_bytes{fstype=~\"ext.?|xfs\",}) *100/(node_filesystem_avail_bytes {fstype=~\"ext.?|xfs\",}+(node_filesystem_size_bytes{fstype=~\"ext.?|xfs\",}-node_filesystem_free_bytes{fstype=~\"ext.?|xfs\",})))by(ecs_cname,instance))","annotations":"节点：${instance}，磁盘使用率过高，当前：${value}%，磁盘空间不足会导致文件无法写入、新日志无法记录，甚至可能ion":60,"rules":[{"forDuration":60,"severity":"P0","expr":"\u003e85"}]}' WHERE rule_id = 'a-d1sgbf5p1r5s73e0men0'
✅ 所有metrics告警规则表更新完成
📣 开始刷值班表数据结构
...
📊 查询到的记录数量： 0
✅ 所有值班表更新完成
📣 开始刷阿里云SLS配置数据结构

2025/08/29 00:46:09 /Users/feiyi/update_w8t/utils/utils.go:140
[0.987ms] [rows:1] SELECT * FROM `alert_rules` WHERE datasource_type = 'AliCloudSLS'
📊 查询到的记录数量： 1

2025/08/29 00:46:09 /Users/feiyi/update_w8t/utils/utils.go:167
[4.511ms] [rows:1] UPDATE `alert_rules` SET `ali_cloud_sls_config`='{"project":"12","logstore":["12"],"logQL":"123","logScope":1}' WHERE rule_id = 'a-d2o5sem2uivc739qrolg'
✅ 所有阿里云SLS规则配置数据结更新完成
```
