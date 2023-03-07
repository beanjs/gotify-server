package model

type Bark struct {
	ID    uint   `gorm:"primary_key;unique_index;AUTO_INCREMENT"`
	Key   string `gorm:"type:varchar(180);unique_index"`
	Token string `gorm:"type:varchar(180);unique_index"`
}
