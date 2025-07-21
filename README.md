本项目致力于为WatchAlert v3.5.0用户升级v3.6+时，数据库字段有变化，用来刷新告警规则的数据库字段

# 项目背景
作者给[WatchAlert](https://github.com/opsre/WatchAlert)项目发起人提出需求，告警规则中，存在多条告警等级时，希望每个告警等级都能有对应的持续时间，所以本项目诞生了。

# 项目功能
协助WatchAlert项目更新v3.6.0之前将数据库中alert_rules表中，数据源为prometheus和victoriametrics的所有记录中prometheus_config字段进行update

# 具体实现
将prometheus_config字段中的forDuration和value，添加到数组rules中的每个对象里，更新后如下：
```shell
# 原prometheus_config字段
{"promQL":"round(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\",}[5m])) by (instance) *
 100))","annotations":"节点：${instance}，CPU使用率过高，当前：${value}%，高 CPU 使用率意味着服务器可能
接近处理能力上限，影响性能，导致应用程序响应变慢或服务中断！","forDuration":60,"rules":[{"severity":"P0",
"expr":"\u003e80"},{"severity":"P1","expr":"\u003e75"},{"severity":"P2","expr":"\u003e70"}]}

# 修改后
{"promQL":"round(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\",}[5m])) by (instance) *
 100))","annotations":"节点：${instance}，CPU使用率过高，当前：${value}%，高 CPU 使用率意味着服务器可能
接近处理能力上限，影响性能，导致应用程序响应变慢或服务中断！","forDuration":60,"rules":[{"forDuration":6
0,"severity":"P0","expr":"\u003e80"},{"forDuration":60,"severity":"P1","expr":"\u003e75"},
{"forDuration":60,"severity":"P2","expr":"\u003e70"}]}
```
更新前后端后，页面表现如下：
<img width="1412" height="497" alt="截屏2025-07-17 22 49 57" src="https://github.com/user-attachments/assets/da0e5523-c863-470c-b754-17cbb4cfffff" />

# 使用方法
[release](https://github.com/mupeifeiyi/update_w8t/releases)页面查看和WatchAlert对应版本的二进制文件
下载到部署WatchAlert的服务器中，这里以compose为例
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
部署时没有修改任何配置，可以直接复制示例进行刷数据库，可重复执行
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
