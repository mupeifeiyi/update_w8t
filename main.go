package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"update_w8t/models"
)

var (
	dsnFlag  = flag.String("dsn", "", "MySQL DSN 连接字符串")
	showHelp = flag.Bool("h", false, "显示帮助信息")
)

func usage() {
	fmt.Println(`用法: update_w8t --dsn=<dsn字符串>

参数说明：
--dsn      必填，MySQL连接字符串
-h         显示帮助信息

示例：
update_w8t --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if *showHelp {
		flag.Usage()
		return
	}

	if *dsnFlag == "" {
		fmt.Println("❌ 错误：必须指定 --dsn 参数")
		flag.Usage()
		return
	}

	dsn := *dsnFlag
	fmt.Printf("✅ 正在使用DSN连接数据库: %s\n", maskPassword(dsn))

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("❌ 数据库连接失败: " + err.Error())
	}

	// 设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		panic("❌ 获取数据库连接池失败")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	var alertRules []models.AlertRule
	db.Where("datasource_type IN (?)", []string{"prometheus", "victoriametrics"}).
		Find(&alertRules)
	fmt.Println("📊 查询到的记录数量：", len(alertRules))

	for i := range alertRules {
		alertRule := &alertRules[i]

		for i := range alertRule.PrometheusConfig.Rules {
			alertRule.PrometheusConfig.Rules[i].ForDuration = alertRule.PrometheusConfig.ForDuration
		}

		configBytes, err := json.Marshal(alertRule.PrometheusConfig)
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败，ruleId: %s, error: %v\n", alertRule.RuleId, err)
			continue
		}

		err = db.Model(&models.AlertRule{}).
			Where("rule_id = ?", alertRule.RuleId).
			Update("prometheus_config", configBytes).
			Error

		if err != nil {
			fmt.Printf("❌ 更新失败，ruleId: %s, error: %v\n", alertRule.RuleId, err)
		} else {
			fmt.Printf("✅ 已更新 PrometheusConfig，ruleId: %s\n", alertRule.RuleId)
		}
	}

	fmt.Println("✅ 所有规则更新完成")
}

// maskPassword 隐藏 DSN 中的密码部分
func maskPassword(dsn string) string {
	atIndex := strings.Index(dsn, "@")
	if atIndex == -1 {
		return dsn
	}
	beforeAt := dsn[:atIndex]
	afterAt := dsn[atIndex:]

	// 找到用户名和密码部分（形如 user:pass@...）
	colon := strings.LastIndex(beforeAt, ":")
	if colon == -1 {
		return dsn
	}

	return beforeAt[:colon+1] + "****" + afterAt
}
