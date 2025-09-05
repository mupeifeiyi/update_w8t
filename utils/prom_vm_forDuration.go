package utils

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"update_w8t/models"
)

// 修改prom规则支持不同告警级别拥有不同的持续时间
func ProcessAlertRule(db *gorm.DB) {
	fmt.Println("📣 开始刷metrics告警规则数据结构")

	var alertRules []models.AlertRule
	db.Where("datasource_type IN (?)", []string{"prometheus", "victoriametrics"}).
		Find(&alertRules)
	fmt.Println("📊 查询到的记录数量：", len(alertRules))

	for i := range alertRules {
		alertRule := &alertRules[i]

		if alertRule.PrometheusConfig.ForDuration <= 0 {
			continue
		}

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
		}
	}

	fmt.Println("✅ 所有metrics告警规则表更新完成")
}

func ProcessRuleTemplate(db *gorm.DB) {
	fmt.Println("📣 开始刷metrics告警规则模版数据结构")

	var ruleTemplates []models.RuleTemplate
	db.Where("datasource_type IN (?)", []string{"prometheus", "victoriametrics"}).
		Find(&ruleTemplates)
	fmt.Println("📊 查询到的记录数量：", len(ruleTemplates))

	for i := range ruleTemplates {
		ruleTemplate := &ruleTemplates[i]

		if ruleTemplate.PrometheusConfig.ForDuration <= 0 {
			continue
		}

		for i := range ruleTemplate.PrometheusConfig.Rules {
			ruleTemplate.PrometheusConfig.Rules[i].ForDuration = ruleTemplate.PrometheusConfig.ForDuration
		}

		configBytes, err := json.Marshal(ruleTemplate.PrometheusConfig)
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败， error: %v\n", err)
			continue
		}

		err = db.Model(&models.RuleTemplate{}).
			Where("rule_name = ?", ruleTemplate.RuleName).
			Update("prometheus_config", configBytes).
			Error

		if err != nil {
			fmt.Printf("❌ 更新失败，error: %v\n", err)
		}
	}

	fmt.Println("✅ 所有metrics告警规则模版更新完成")
}
