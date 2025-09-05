package utils

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"update_w8t/models"
)

// 新的告警规则中，logstore的类型是数组
type NewSLSConfig struct {
	Project  string   `json:"project"`
	Logstore []string `json:"logstore"` // 新格式为数组
	LogQL    string   `json:"logQL"`
	LogScope int      `json:"logScope"`
}

// 尝试解析为旧格式
var oldConfig struct {
	Project  string `json:"project"`
	Logstore string `json:"logstore"` // 旧格式是字符串
	LogQL    string `json:"logql"`
	LogScope int    `json:"logScope"`
}

// 临时结构体，用于查询时绕过 JSON 字段反序列化错误
type AlertRule struct {
	models.AlertRule                  // 嵌入原始结构体，继承所有字段
	AliCloudSLSConfig json.RawMessage `gorm:"column:ali_cloud_sls_config" json:"alicloudSLSConfig"`
}

// 修改SLSConfig中logstore数据格式，支持多个logstore，要修改所有规则，包括非SLS数据源的，否则前端无法正常展示
func ProcessAliSLSConfigAlertRule(db *gorm.DB) {
	fmt.Println("📣 开始刷告警规则中logstore的数据结构")

	var alertRules []AlertRule
	db.Find(&alertRules)

	fmt.Println("📊 查询到的记录数量：", len(alertRules))

	for i := range alertRules {
		rule := &alertRules[i]

		if err := json.Unmarshal(rule.AliCloudSLSConfig, &oldConfig); err == nil {
			// 成功解析为旧格式，转换为新格式
			newConfig := NewSLSConfig{
				Project:  oldConfig.Project,
				Logstore: []string{oldConfig.Logstore},
				LogQL:    oldConfig.LogQL,
				LogScope: oldConfig.LogScope,
			}
			// 序列化新配置
			configBytes, err := json.Marshal(newConfig)
			if err != nil {
				fmt.Printf("❌ JSON 序列化失败，ruleId: %s, error: %v\n", rule.RuleId, err)
				continue
			}
			// 更新数据库
			if err := db.Model(&models.AlertRule{}).
				Where("rule_id = ?", rule.RuleId).
				Update("ali_cloud_sls_config", configBytes).
				Error; err != nil {
				fmt.Printf("❌ 数据库更新失败，ruleId: %s, error: %v\n", rule.RuleId, err)
			}
		} else {
			// 尝试解析为新格式
			var newConfig NewSLSConfig
			if err := json.Unmarshal(rule.AliCloudSLSConfig, &newConfig); err == nil {
				// 已是新格式，跳过处理
				fmt.Printf("ℹ️ 规则 %s 已是新格式，跳过转换\n", rule.RuleId)
			} else {
				// 无法解析为新旧格式，记录错误
				fmt.Printf("❌ 无法解析配置，ruleId: %s, error: %v\n", rule.RuleId, err)
			}
		}
	}
	fmt.Println("✅ 所有告警规则中logstore的数据结构更新完成")
}

type RuleTemplate struct {
	models.RuleTemplate                 // 嵌入原始结构体，继承所有字段
	AliCloudSLSConfig   json.RawMessage `gorm:"column:ali_cloud_sls_config" json:"alicloudSLSConfig"`
}

// 规则模版中的SLSConfig字段也同样要刷
func ProcessSLSRuleTemplate(db *gorm.DB) {
	fmt.Println("📣 开始刷告警规则模版中logstore的数据结构")

	var ruleTemplates []RuleTemplate
	db.Find(&ruleTemplates)
	fmt.Println("📊 查询到的记录数量：", len(ruleTemplates))

	for i := range ruleTemplates {
		ruleTemplates := &ruleTemplates[i]

		if err := json.Unmarshal(ruleTemplates.AliCloudSLSConfig, &oldConfig); err == nil {
			// 成功解析为旧格式，转换为新格式
			newConfig := NewSLSConfig{
				Project:  oldConfig.Project,
				Logstore: []string{oldConfig.Logstore},
				LogQL:    oldConfig.LogQL,
				LogScope: oldConfig.LogScope,
			}
			// 序列化新配置
			configBytes, err := json.Marshal(newConfig)
			if err != nil {
				fmt.Printf("❌ JSON 序列化失败，ruleName: %s, error: %v\n", ruleTemplates.RuleName, err)
				continue
			}
			// 更新数据库
			if err := db.Model(&models.RuleTemplate{}).
				Where("rule_name = ?", ruleTemplates.RuleName).
				Update("ali_cloud_sls_config", configBytes).
				Error; err != nil {
				fmt.Printf("❌ 数据库更新失败，ruleId: %s, error: %v\n", ruleTemplates.RuleName, err)
			}
		} else {
			// 尝试解析为新格式
			var newConfig NewSLSConfig
			if err := json.Unmarshal(ruleTemplates.AliCloudSLSConfig, &newConfig); err == nil {
				// 已是新格式，跳过处理
				fmt.Printf("ℹ️ 规则 %s 已是新格式，跳过转换\n", ruleTemplates.RuleName)
			} else {
				// 无法解析为新旧格式，记录错误
				fmt.Printf("❌ 无法解析配置，ruleName: %s, error: %v\n", ruleTemplates.RuleName, err)
			}
		}
	}

	fmt.Println("✅ 所有告警规则模版中的logstore的数据结构更新完成")
}
