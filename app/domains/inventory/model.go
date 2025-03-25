package inventory

type Inventory struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	ProductID uint `gorm:"index;unique"`
	Stock     int  `gorm:"type:int;not null"`
}
