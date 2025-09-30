package repository

import (
	"web_backend/internal/app/serializer"
	"web_backend/internal/app/ds"
	"database/sql"
	"errors"
	"fmt"
	"time"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errNoDraft = errors.New("no draft for this user")

func (r *Repository) GetAllApplications(from, to time.Time, status string) ([]ds.Application, error) {
	var applications []ds.Application
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("date_create > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("date_create < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("application_id").Find(&applications).Error
	if err != nil {
		return nil, err
	}
	return applications, nil
}

func (r *Repository) GetDevicesApplcations(application_id int) ([]ds.ApplicationDevices, error) {
	var applicationDevice []ds.ApplicationDevices
	err := r.db.Where("application_id = ?", application_id).Find(&applicationDevice).Error
	if err != nil {
		return nil, err
	}
	return applicationDevice, nil
}

func (r *Repository) GetDevicesApplication(device_id int, application_id int) (ds.ApplicationDevices, error) {
	var applicationDevice ds.ApplicationDevices
	err := r.db.Where("device_id = ? and application_id = ?", device_id, application_id).First(&applicationDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.ApplicationDevices{}, fmt.Errorf("%w: device application not found", ErrNotFound)
		}
		return ds.ApplicationDevices{}, err
	}
	return applicationDevice, nil
}

func (r *Repository) GetApplicationDevices(id int) ([]ds.Device, ds.Application, error) {
	application, err := r.GetSingleApplication(id)
	if err != nil {
		return []ds.Device{}, ds.Application{}, err
	}

	var devices []ds.Device
	sub := r.db.Table("application_devices").Where("application_id = ?", application.Application_ID)
	err = r.db.Order("device_id DESC").Where("device_id IN (?)", sub.Select("device_id")).Find(&devices).Error

	if err != nil {
		return []ds.Device{}, ds.Application{}, err
	}

	return devices, application, nil
}

func (r *Repository) CheckCurrentApplicationDraft(creator_ID uint) (ds.Application, error) {
    // if creatorID == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
    // }
    
	var application ds.Application
	res := r.db.Where("creator_id = ? AND status = ?", creator_ID, "draft").Limit(1).Find(&application)
	if res.Error != nil {
		return ds.Application{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.Application{}, ErrNoDraft
	}
	return application, nil
}

func (r *Repository) GetApplicationDraft(creator_ID uint) (ds.Application, bool, error) {
    // if creatorID == 0 {
    //     return ds.Research{}, false, fmt.Errorf("%w: user not authenticated", ErrNotAllowed)
    // }

	application, err := r.CheckCurrentApplicationDraft(creator_ID)
	if errors.Is(err, ErrNoDraft) {
		application = ds.Application{
			Status:     "draft",
			Creator_ID:  creator_ID,
			Created_At: time.Now(),
		}
		result := r.db.Create(&application)
		if result.Error != nil {
			return ds.Application{}, false, result.Error
		}
		return application, true, nil
	} else if err != nil {
		return ds.Application{}, false, err
	}
	return application, true, nil
}

func (r *Repository) GetApplicationCount(creator_ID uint) int64 {
    if creator_ID == 0 {
        return 0
    }
    
	var count int64
	application, err := r.CheckCurrentApplicationDraft(creator_ID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.ApplicationDevices{}).Where("application_id = ?", application.Application_ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_devices:", err)
	}

	return count
}

func (r *Repository) DeleteCalculation(application_id int) error{
	return r.db.Exec("UPDATE applications SET status = 'deleted' WHERE id = ?", application_id).Error
}

func (r *Repository) GetSingleApplication(id int) (ds.Application, error) {
	if id < 0 {
		return ds.Application{}, errors.New("неверное id, должно быть >= 0")
	}
    
    // userId := r.GetUserID()
    // if userId == 0 {
    //     return ds.Research{}, fmt.Errorf("%w: пользователь не авторизирован", ErrNotAllowed)
    // }
    
	// user, err := r.GetUserByID(userId)
	// if err != nil {
	// 	return ds.Research{}, err
	// }
    
	var application ds.Application
	err := r.db.Where("application_id = ?", id).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Application{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Application{}, err
	} else if application.Status == "deleted"  {
		return ds.Application{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return application, nil
}

func (r *Repository) FormApplication(application_id int, status string) (ds.Application, error) {
	application, err := r.GetSingleApplication(application_id)
	if err != nil {
		return ds.Application{}, err
	}

	// user, err := r.GetUserByID(r.GetUserID())
	// if err != nil{
	// 	return ds.Research{}, fmt.Errorf("%w: пользователь на авторизирован", ErrNotAllowed)
	// }

	// if research.CreatorID != r.userId && !user.IsModerator{
	// 	return ds.Research{}, fmt.Errorf("%w: у вас нет прав чтобы эта заявка имела статус %s", ErrNotAllowed, status)
	// }

	if application.Status != "draft" {
		return ds.Application{}, fmt.Errorf("эта заявка не может быть %s", status)
	}
	
	if status != "deleted"{
		if application.Amperage < 0 {
			return ds.Application{}, errors.New("вы не написали нагрузку системы")
		}
		applicationDevices, _ := r.GetDevicesApplcations(int(application.Application_ID))
		for _, applicationDevices := range applicationDevices{
				if applicationDevices.Notes == ""{
					return ds.Application{}, errors.New("вы не написали заметку" )			
				}
		}
	}	

	err = r.db.Model(&application).Updates(ds.Application{
		Status: status,
		Forming_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},

	}).Error
	if err != nil {
		return ds.Application{}, err
	}

	return application, nil
}

func (r *Repository) EditApplication(id int, applicationJSON serializer.ApplicationJSON) (ds.Application, error) {
	application := ds.Application{}
	if id < 0 {
		return ds.Application{}, errors.New("неправильное id, должно быть >= 0")
	}
	if applicationJSON.Amperage < 0 {  
		return ds.Application{}, errors.New("неправильная нагрузка")
	}
	err := r.db.Where("application_id = ? and status != 'deleted'", id).First(&application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Application{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.Application{}, err
	}
	err = r.db.Model(&application).Updates(serializer.ApplicationFromJSON(applicationJSON)).Error
	if err != nil {
		return ds.Application{}, err
	}
	return application, nil
}

func CalculateDeviceAmperage(power float64) (float64, error) {
	if power < 0 {
		return 0, errors.New("неправильная мощность")
	}
	return float64(power)*1000/220 , nil
}

func (r *Repository) FinishApplication(id int, status string) (ds.Application, error) {
	if status != "completed" && status != "rejected" {
		return ds.Application{}, errors.New("неверный статус")
	}

	user, err := r.GetUserByID(r.GetUserID())
	if err != nil {
		return ds.Application{}, err
	}

	if !user.IsModerator {
		return ds.Application{}, fmt.Errorf("%w: вы не модератор", ErrNotAllowed)
	}

	application, err := r.GetSingleApplication(id)
	if err != nil {
		return ds.Application{}, err
	} else if application.Status != "formed" {
		return ds.Application{}, fmt.Errorf("это исследование не может быть %s", status)
	}

	err = r.db.Model(&application).Updates(ds.Application{
		Status: status,
		Finish_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Moderator_ID: uint(user.User_ID),
	}).Error
	if err != nil {
		return ds.Application{}, err
	}

	if status == "completed" {
		applicationsDevice, err := r.GetDevicesApplcations(int(application.Application_ID))
		if err != nil {
			return ds.Application{}, err
		}
		for _, applicationDevice := range applicationsDevice {
			device, err := r.GetDevice(int(applicationDevice.Device_ID))
			if err != nil {
				return ds.Application{}, err
			}
			device_amperage, err := CalculateDeviceAmperage(device.Dev_Power)
			if err != nil {
				return ds.Application{}, err
			}
			err = r.db.Model(&applicationDevice).Updates(ds.ApplicationDevices{
				Amperage: float64(device_amperage),
			}).Error
			if err != nil {
				return ds.Application{}, err
			}
		}
	}

	return application, nil
}