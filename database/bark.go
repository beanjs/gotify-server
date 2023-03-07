package database

import "github.com/gotify/server/v2/model"

func (d *GormDatabase) CreateBark(bark *model.Bark) error {
	return d.DB.Create(bark).Error
}

func (d *GormDatabase) GetBarks() ([]*model.Bark, error) {
	var barks []*model.Bark
	err := d.DB.Find(&barks).Error
	return barks, err
}
