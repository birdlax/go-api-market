package utils

import (
	"log"
	"time"

	"backend/domain"
	"gorm.io/gorm"
)

func StartUserCountLogger(db *gorm.DB) {
	go func() {
		for {
			var count int64
			if err := db.Model(&domain.User{}).Count(&count).Error; err != nil {
				log.Println("Error counting users:", err)
			} else {
				log.Println("Current user count:", count)
			}
			time.Sleep(10 * time.Second)
		}
	}()
}
