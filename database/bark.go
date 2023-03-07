package database

import (
	"github.com/gotify/server/v2/model"
	"github.com/jinzhu/gorm"
)

func (d *GormDatabase) CreateBark(bark *model.Bark) error {
	e := new(model.Bark)
	err := d.DB.Where("token = ?", bark.Token).Find(e).Error
	if err == gorm.ErrRecordNotFound {
		return d.DB.Create(bark).Error
	}
	return nil
}

func (d *GormDatabase) GetBarks() ([]*model.Bark, error) {
	var barks []*model.Bark
	err := d.DB.Find(&barks).Error
	return barks, err
}
