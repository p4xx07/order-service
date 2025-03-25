package user

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"uniqueIndex;type:varchar(100)"`
	CreatedAt time.Time
}
