package main

import (
	"encoding/json"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"update_w8t/models"
)

func main() {
	dsn := "root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		panic("无法连接数据库")
	}

	// 可选：启用 Debug 模式输出 SQL 语句
	db = db.Debug()

	sqlDB, err := db.DB()
	if err != nil {
		panic("无法获取 DB 连接池")
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	var alertRules []models.AlertRule
	db.Where("datasource_type = ?", "prometheus").Find(&alertRules)
	fmt.Println("查询到的记录数量：", len(alertRules))

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

		// ✅ 更新当前 alertRule 的 prometheus_config 字段，只更新当前 alertRule 对应的记录
		err = db.Model(&models.AlertRule{}).
			Where("rule_id = ?", alertRule.RuleId). // 🎯 只更新当前这条数据
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
