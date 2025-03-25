package user

import "time"

type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"type:varchar(100)"`
	Email     string `gorm:"uniqueIndex;type:varchar(100)"`
	Password  string `gorm:"type:varchar(255)"`
	CreatedAt time.Time
}
