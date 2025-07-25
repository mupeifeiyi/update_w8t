package models

type RuleTemplate struct {
	Type                 string              `json:"type"`
	RuleGroupName        string              `json:"ruleGroupName"`
	RuleName             string              `json:"ruleName"  gorm:"type:varchar(255);not null"`
	DatasourceType       string              `json:"datasourceType"`
	EvalInterval         int64               `json:"evalInterval"`
	ForDuration          int64               `json:"forDuration"`
	RepeatNoticeInterval int64               `json:"repeatNoticeInterval"`
	Description          string              `json:"description"`
	EffectiveTime        EffectiveTime       `json:"effectiveTime" gorm:"effectiveTime;serializer:json"`
	PrometheusConfig     PrometheusConfig    `json:"prometheusConfig" gorm:"prometheusConfig;serializer:json"`
	AliCloudSLSConfig    AliCloudSLSConfig   `json:"alicloudSLSConfig" gorm:"alicloudSLSConfig;serializer:json"`
	LokiConfig           LokiConfig          `json:"lokiConfig" gorm:"lokiConfig;serializer:json"`
	JaegerConfig         JaegerConfig        `json:"jaegerConfig" gorm:"JaegerConfig;serializer:json"`
	KubernetesConfig     KubernetesConfig    `json:"kubernetesConfig" gorm:"kubernetesConfig;serializer:json"`
	ElasticSearchConfig  ElasticSearchConfig `json:"elasticSearchConfig" gorm:"elasticSearchConfig;serializer:json"`
}
