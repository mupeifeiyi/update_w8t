package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
	"update_w8t/utils"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	dsnFlag  = flag.String("dsn", "", "MySQL DSN 连接字符串")
	showHelp = flag.Bool("h", false, "显示帮助信息")
)

func usage() {
	progName := os.Args[0]
	fmt.Printf(`用法: %s --dsn=<dsn字符串>

参数说明：
--dsn      必填，MySQL连接字符串
-h         显示帮助信息

示例：
%s --dsn="root:w8t.123@tcp(127.0.0.1:3306)/watchalert?charset=utf8mb4&parseTime=True&loc=Local"
`, progName, progName)
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

	utils.ProcessRuleTemplate(db)
	utils.ProcessAlertRule(db)
	utils.ProcessCalendar(db)
	utils.ProcessAliSLSConfigAlertRule(db)
	utils.ProcessSLSRuleTemplate(db)

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
