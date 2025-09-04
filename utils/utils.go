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

func ProcessCalendar(db *gorm.DB) {
	fmt.Println("📣 开始刷值班表数据结构")

	var dutys []models.DutySchedule
	db.Model(&models.DutySchedule{}).Find(&dutys)

	fmt.Println("📊 查询到的记录数量：", len(dutys))

	for i := range dutys {
		duty := &dutys[i]
		if duty.UserId == "" && duty.Username == "" {
			continue
		}
		duty.Users = []models.DutyUser{
			{
				UserId:   duty.UserId,
				Username: duty.Username,
			},
		}

		bytes, err := json.Marshal(duty.Users)
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败，error: %v\n", err)
			continue
		}

		err = db.Model(&models.DutySchedule{}).
			Where("duty_id = ? and time = ?", duty.DutyId, duty.Time).
			Update("users", bytes).
			Error

		if err != nil {
			fmt.Printf("❌ 更新失败 error: %v\n", err)
		}
	}

	fmt.Println("✅ 所有值班表更新完成")
}

// 新的告警规则中，logstore的类型是数组
type NewSLSConfig struct {
	Project  string   `json:"project"`
	Logstore []string `json:"logstore"` // 新格式为数组
	LogQL    string   `json:"logQL"`
	LogScope int      `json:"logScope"`
}

// 修改SLSConfig中logstore数据格式，支持多个logstore，要修改所有规则，包括非SLS数据源的，否则前端无法正常展示
func ProcessAliSLSConfigAlertRule(db *gorm.DB) {
	fmt.Println("📣 开始刷告警规则中logstore的数据结构")

	var alertRules []models.AlertRule
	db.Find(&alertRules)
	fmt.Println("📊 查询到的记录数量：", len(alertRules))

	for i := range alertRules {
		alertRule := &alertRules[i]

		// 1. 存储旧配置
		oldConfig := alertRule.AliCloudSLSConfig

		// 2. 转换为新格式
		newConfig := NewSLSConfig{
			Project:  oldConfig.Project,
			Logstore: []string{oldConfig.Logstore}, // 字符串 → 数组
			LogQL:    oldConfig.LogQL,
			LogScope: oldConfig.LogScope,
		}

		// 3. 序列化新配置
		configBytes, err := json.Marshal(newConfig)
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败，ruleId: %s, error: %v\n", alertRule.RuleId, err)
			continue
		}

		// 4. 更新数据库
		if err := db.Model(&models.AlertRule{}).
			Where("rule_id = ?", alertRule.RuleId).
			Update("ali_cloud_sls_config", configBytes).
			Error; err != nil {
			fmt.Printf("❌ 数据库更新失败，ruleId: %s, error: %v\n", alertRule.RuleId, err)
		}
	}

	fmt.Println("✅ 所有告警规则中logstore的数据结构更新完成")
}

// 规则模版中的SLSConfig字段也同样要刷
func ProcessSLSRuleTemplate(db *gorm.DB) {
	fmt.Println("📣 开始刷告警规则模版中logstore的数据结构")

	var ruleTemplates []models.RuleTemplate
	db.Find(&ruleTemplates)
	fmt.Println("📊 查询到的记录数量：", len(ruleTemplates))

	for i := range ruleTemplates {
		ruleTemplates := &ruleTemplates[i]

		// 1. 存储旧配置
		oldConfig := ruleTemplates.AliCloudSLSConfig

		// 2. 转换为新格式
		newConfig := NewSLSConfig{
			Project:  oldConfig.Project,
			Logstore: []string{oldConfig.Logstore}, // 字符串 → 数组
			LogQL:    oldConfig.LogQL,
			LogScope: oldConfig.LogScope,
		}

		// 3. 序列化新配置
		configBytes, err := json.Marshal(newConfig)
		if err != nil {
			fmt.Printf("❌ JSON 序列化失败，ruleName: %s, error: %v\n", ruleTemplates.RuleName, err)
			continue
		}

		// 4. 更新数据库
		if err := db.Model(&models.RuleTemplate{}).
			Where("rule_name = ?", ruleTemplates.RuleName).
			Update("ali_cloud_sls_config", configBytes).
			Error; err != nil {
			fmt.Printf("❌ 数据库更新失败，ruleId: %s, error: %v\n", ruleTemplates.RuleName, err)
		}
	}

	fmt.Println("✅ 所有告警规则模版中的logstore的数据结构更新完成")
}
