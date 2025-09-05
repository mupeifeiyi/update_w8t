package utils

import (
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"update_w8t/models"
)

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
