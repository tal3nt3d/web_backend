package repository

import (
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/ds"
	"errors"
	"fmt"
	"gorm.io/gorm"
)

func (r *Repository) DeleteDeviceFromApplication(application_id int, device_id int) (ds.Application, error) {
	// userId := r.userId
    // if userId == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }
    
	var application ds.Application
	err := r.db.Where("application_id = ?", application_id).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Application{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, application_id)
		}
		return ds.Application{}, err
	}
    
	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: Вы не создатель этого исследования", ErrNotAllowed)
	// }
    
	err = r.db.Where("device_id = ? and application_id = ?", device_id, application_id).Delete(&ds.ApplicationDevices{}).Error
	if err != nil {
		return ds.Application{}, err
	}
	return application, nil
}

func (r *Repository) EditDeviceFromApplication(application_id int, device_id int, applicationDeviceJSON serializer.ApplicationDeviceJSON) (ds.ApplicationDevices, error) {
	var applicationsDevice ds.ApplicationDevices
	err := r.db.Model(&applicationsDevice).Where("device_id = ? and application_id = ?", device_id, application_id).Updates(serializer.ApplicationDeviceFromJSON(applicationDeviceJSON)).First(&applicationsDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ApplicationDevices{}, fmt.Errorf("%w: планеты в исследовании", ErrNotFound)
		}
		return ds.ApplicationDevices{}, err
	}
	return applicationsDevice, nil
}