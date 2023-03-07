package database

import (
	"github.com/gotify/server/v2/model"
	"github.com/jinzhu/gorm"
	"github.com/lithammer/shortuuid/v3"
)

func (d *GormDatabase) FindBarkByToken(token string) (*model.Bark, error) {
	bark := new(model.Bark)
	err := d.DB.Where("token = ?", token).Find(bark).Error
	if err == gorm.ErrRecordNotFound {
		bark.Key = shortuuid.New()
		bark.Token = token

		return bark, d.DB.Create(bark).Error
	}
	return bark, nil
}

func (d *GormDatabase) DeleteByKey(key string) error {
	return d.DB.Where("key = ?", key).Delete(&model.Bark{}).Error
}

func (d *GormDatabase) GetBarks() ([]*model.Bark, error) {
	var barks []*model.Bark
	err := d.DB.Find(&barks).Error
	return barks, err
}
