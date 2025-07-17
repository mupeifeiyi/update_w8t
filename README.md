本项目致力于为WatchAlert v3.5.0用户升级v3.6+时，数据库字段有变化，用来刷新告警规则的数据库字段

# 项目背景
本人给[WatchAlert](https://github.com/opsre/WatchAlert)项目发起人提出需求，告警规则中，存在多条告警等级时，希望每个告警等级都能有对应的持续时间，所以本项目诞生了。

# 项目功能
协助WatchAlert项目更新v3.6.0之前将数据库中alert_rules表中，数据源为prometheus和victoriametrics的所有记录中prometheus_config字段进行update

# 具体实现
将prometheus_config字段中的forDuration和值，添加到数组rules中的每个对象里，更新后如下：
```shell
# 原prometheus_config字段
{"promQL":"round(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\",}[5m])) by (instance) * 100))","annotations":"节点：${instance}，CPU使用率过高，当前：${value}%，高 CPU 使用率意味着服务器可能接近处理能力上限，影响性能，导致应用程序响应变慢或服务中断！","forDuration":60,"rules":[{"severity":"P0","expr":"\u003e80"},{"severity":"P1","expr":"\u003e75"},{"severity":"P2","expr":"\u003e70"}]}

# 修改后
{"promQL":"round(100 - (avg(irate(node_cpu_seconds_total{mode=\"idle\",}[5m])) by (instance) * 100))","annotations":"节点：${instance}，CPU使用率过高，当前：${value}%，高 CPU 使用率意味着服务器可能接近处理能力上限，影响性能，导致应用程序响应变慢或服务中断！","forDuration":60,"rules":[{"forDuration":60,"severity":"P0","expr":"\u003e80"},{"forDuration":60,"severity":"P1","expr":"\u003e75"},{"forDuration":60,"severity":"P2","expr":"\u003e70"}]}
```
更新前后端后，页面表现如下：
<img width="1412" height="497" alt="截屏2025-07-17 22 49 57" src="https://github.com/user-attachments/assets/da0e5523-c863-470c-b754-17cbb4cfffff" />
