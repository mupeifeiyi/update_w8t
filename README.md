本项目致力于为WatchAlert v3.5.0用户升级v3.6+时，数据库字段有变化，用来刷新告警规则的数据库字段

# 项目背景
作者给[WatchAlert](https://github.com/opsre/WatchAlert)项目发起人提出需求：
1. 告警规则中，存在多条告警等级时，希望每个告警等级都能有对应的持续时间
2. 值班表中希望单天内可以多人值班

# 项目功能
协助WatchAlert项目更新:
1. v3.6.0之前将数据库中alert_rules表中，数据源为prometheus和victoriametrics的所有记录中prometheus_config字段进行update
2. 对新增的值班组功能进行数据库结构更改，现已支持同一天多人值班

更新前后端后，页面表现如下：
<img width="1412" height="497" alt="截屏2025-07-17 22 49 57" src="https://github.com/user-attachments/assets/da0e5523-c863-470c-b754-17cbb4cfffff" />

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
❌ 错误：必须指定 --dsn 参数
用法: ./u8t-linux-amd64 --dsn=<dsn字符串>

参数说明：
--dsn      必填，MySQL连接字符串
-h         显示帮助信息

示例：
./u8t-linux-amd64 --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
```
部署时没有修改任何配置，可以直接复制示例进行刷数据库，重复执行不会影响数据
```shell
./u8t-linux-amd64 --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
✅ 正在使用DSN连接数据库: root:****@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local

2025/07/21 15:54:08 /Users/feiyi/update_w8t/main.go:71
[2.624ms] [rows:55] SELECT * FROM `alert_rules` WHERE datasource_type IN ('prometheus','victoriametrics')
📊 查询到的记录数量： 55

2025/07/21 15:54:08 /Users/feiyi/update_w8t/main.go:89
[4.840ms] [rows:1] UPDATE `alert_rules` SET `prometheus_config`='{"promQL":"max_over_time(reloader_last_reload_successful{namespace=~\".+\"}[5m])","annotations":"Pod ${labels.pod} 中的 config-reloader sidecar 在尝试同步配置时遇到错误","forDuration":600,"rules":[{"forDuration":600,"severity":"P1","expr":"== 0"}]}' WHERE rule_id = 'a-d1lmtnc06bis73ebshj0'
✅ 已更新 PrometheusConfig，ruleId: a-d1lmtnc06bis73ebshj0
...
2025/07/21 16:10:35 /Users/feiyi/update_w8t/main.go:89
[0.506ms] [rows:0] UPDATE `alert_rules` SET `prometheus_config`='{"promQL":"|-","annotations":"Namespace {{ $labels.namespace }} is using {{ $value | humanizePercentage","forDuration":900,"rules":[{"forDuration":900,"severity":"P0","expr":"\u003e 0"}]}' WHERE rule_id = 'a-d1r2mv406bis73ccrn8g'
✅ 已更新 PrometheusConfig，ruleId: a-d1r2mv406bis73ccrn8g
✅ 所有规则更新完成
```
