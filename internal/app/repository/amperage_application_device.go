package repository

import (
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/ds"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func (r *Repository) DeleteDeviceFromAmperageApplication(amperage_application_id int, device_id int) (ds.AmperageApplication, error) {
	// userId := r.userId
    // if userId == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }
    
	var amperage_application ds.AmperageApplication
	err := r.db.Where("amperage_application_id = ?", amperage_application_id).First(&amperage_application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AmperageApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, amperage_application_id)
		}
		return ds.AmperageApplication{}, err
	}
    
	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: Вы не создатель этого исследования", ErrNotAllowed)
	// }
    
	err = r.db.Where("device_id = ? and amperage_application_id = ?", device_id, amperage_application_id).Delete(&ds.AmperageApplicationDevices{}).Error
	if err != nil {
		return ds.AmperageApplication{}, err
	}
	return amperage_application, nil
}

func (r *Repository) EditDeviceFromAmperageApplication(amperage_application_id int, device_id int, amperage_applicationDeviceJSON serializer.AmperageApplicationDeviceJSON) (ds.AmperageApplicationDevices, error) {
	var amperage_applicationsDevice ds.AmperageApplicationDevices
	err := r.db.Model(&amperage_applicationsDevice).Where("device_id = ? and amperage_application_id = ?", device_id, amperage_application_id).Updates(serializer.AmperageApplicationDeviceFromJSON(amperage_applicationDeviceJSON)).First(&amperage_applicationsDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AmperageApplicationDevices{}, fmt.Errorf("%w: устройства в заявке", ErrNotFound)
		}
		return ds.AmperageApplicationDevices{}, err
	}
	return amperage_applicationsDevice, nil
}