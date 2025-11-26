package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/serializer"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var errNoDraft = errors.New("no draft for this user")

func (r *Repository) GetAllAmperageApplications(from, to time.Time, status string) ([]ds.AmperageApplication, error) {
	var amperage_applications []ds.AmperageApplication
	sub := r.db.Where("status != 'deleted' and status != 'draft'")
	if !from.IsZero() {
		sub = sub.Where("created_at > ?", from)
	}
	if !to.IsZero() {
		sub = sub.Where("created_at < ?", to.Add(time.Hour*24))
	}
	if status != "" {
		sub = sub.Where("status = ?", status)
	}
	err := sub.Order("amperage_application_id").Find(&amperage_applications).Error
	if err != nil {
		return nil, err
	}
	return amperage_applications, nil
}

func (r *Repository) GetDevicesAmperageApplications(amperage_application_id int) ([]ds.AmperageApplicationDevices, error) {
	var amperage_applicationDevice []ds.AmperageApplicationDevices
	err := r.db.Where("amperage_application_id = ?", amperage_application_id).Find(&amperage_applicationDevice).Error
	if err != nil {
		return nil, err
	}
	return amperage_applicationDevice, nil
}

func (r *Repository) GetDevicesAmperageApplication(device_id int, amperage_application_id int) (ds.AmperageApplicationDevices, error) {
	var amperage_applicationDevice ds.AmperageApplicationDevices
	err := r.db.Where("device_id = ? and amperage_application_id = ?", device_id, amperage_application_id).First(&amperage_applicationDevice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AmperageApplicationDevices{}, fmt.Errorf("%w: device amperage_application not found", ErrNotFound)
		}
		return ds.AmperageApplicationDevices{}, err
	}
	return amperage_applicationDevice, nil
}

func (r *Repository) GetAmperageApplicationDevices(id int) ([]ds.Device, ds.AmperageApplication, error) {
	amperage_application, err := r.GetSingleAmperageApplication(id)
	if err != nil {
		return []ds.Device{}, ds.AmperageApplication{}, err
	}

	var devices []ds.Device
	sub := r.db.Table("amperage_application_devices").Where("amperage_application_id = ?", amperage_application.Amperage_Application_ID)
	err = r.db.Order("device_id DESC").Where("device_id IN (?)", sub.Select("device_id")).Find(&devices).Error

	if err != nil {
		return []ds.Device{}, ds.AmperageApplication{}, err
	}

	return devices, amperage_application, nil
}

func (r *Repository) CheckCurrentAmperageApplicationDraft(creator_ID uuid.UUID) (ds.AmperageApplication, error) {
	var amperage_application ds.AmperageApplication
	res := r.db.Where("creator_id = ? AND status = ?", creator_ID, "draft").Limit(1).Find(&amperage_application)
	if res.Error != nil {
		return ds.AmperageApplication{}, res.Error
	} else if res.RowsAffected == 0 {
		return ds.AmperageApplication{}, ErrNoDraft
	}
	return amperage_application, nil
}

func (r *Repository) GetAmperageApplicationDraft(creator_ID uuid.UUID) (ds.AmperageApplication, bool, error) {
	amperage_application, err := r.CheckCurrentAmperageApplicationDraft(creator_ID)
	if errors.Is(err, ErrNoDraft) {
		amperage_application = ds.AmperageApplication{
			Status:     "draft",
			Creator_ID: creator_ID,
			Created_At: time.Now(),
		}
		result := r.db.Create(&amperage_application)
		if result.Error != nil {
			return ds.AmperageApplication{}, false, result.Error
		}
		return amperage_application, true, nil
	} else if err != nil {
		return ds.AmperageApplication{}, false, err
	}
	return amperage_application, true, nil
}

func (r *Repository) GetAmperageApplicationCount(creator_ID uuid.UUID) int64 {
	var count int64
	amperage_application, err := r.CheckCurrentAmperageApplicationDraft(creator_ID)
	if err != nil {
		return 0
	}
	err = r.db.Model(&ds.AmperageApplicationDevices{}).Where("amperage_application_id = ?", amperage_application.Amperage_Application_ID).Count(&count).Error
	if err != nil {
		logrus.Println("Error counting records in lists_devices:", err)
	}

	return count
}

func (r *Repository) DeleteCalculation(amperage_application_id int) error {
	return r.db.Exec("UPDATE amperage_applications SET status = 'deleted' WHERE id = ?", amperage_application_id).Error
}

func (r *Repository) GetSingleAmperageApplication(id int) (ds.AmperageApplication, error) {
	if id < 0 {
		return ds.AmperageApplication{}, errors.New("неверное id, должно быть >= 0")
	}

	var amperage_application ds.AmperageApplication
	err := r.db.Where("amperage_application_id = ?", id).First(&amperage_application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AmperageApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.AmperageApplication{}, err
	} else if amperage_application.Status == "deleted" {
		return ds.AmperageApplication{}, fmt.Errorf("%w: заявка удалена", ErrNotAllowed)
	}
	return amperage_application, nil
}

func (r *Repository) FormAmperageApplication(amperage_application_id int, status string) (ds.AmperageApplication, error) {
	amperage_application, err := r.GetSingleAmperageApplication(amperage_application_id)
	if err != nil {
		return ds.AmperageApplication{}, err
	}

	if amperage_application.Status != "draft" {
		return ds.AmperageApplication{}, fmt.Errorf("эта заявка не может быть %s", status)
	}

	if status != "deleted" {
		if amperage_application.Amperage < 0 {
			return ds.AmperageApplication{}, errors.New("вы не написали нагрузку системы")
		}
	}

	err = r.db.Model(&amperage_application).Updates(ds.AmperageApplication{
		Status: status,
		Forming_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.AmperageApplication{}, err
	}

	return amperage_application, nil
}

func (r *Repository) EditAmperageApplication(id int, amperage_applicationJSON serializer.AmperageApplicationJSON) (ds.AmperageApplication, error) {
	amperage_application := ds.AmperageApplication{}
	if id < 0 {
		return ds.AmperageApplication{}, errors.New("неправильное id, должно быть >= 0")
	}
	if amperage_applicationJSON.Amperage < 0 {
		return ds.AmperageApplication{}, errors.New("неправильная нагрузка")
	}
	err := r.db.Where("amperage_application_id = ? and status != 'deleted'", id).First(&amperage_application).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.AmperageApplication{}, fmt.Errorf("%w: заявка с id %d", ErrNotFound, id)
		}
		return ds.AmperageApplication{}, err
	}
	err = r.db.Model(&amperage_application).Updates(serializer.AmperageApplicationFromJSON(amperage_applicationJSON)).Error
	if err != nil {
		return ds.AmperageApplication{}, err
	}
	return amperage_application, nil
}

func CalculateDeviceAmperage(power float64, amount float64) (float64, error) {
	if power < 0 {
		return 0, errors.New("неправильная мощность")
	}
	return float64(power) * 1000 * amount / 220, nil
}

func (r *Repository) FinishAmperageApplication(id int, status string, currentUserID uuid.UUID) (ds.AmperageApplication, error) {
	if status != "completed" && status != "rejected" {
		return ds.AmperageApplication{}, errors.New("неверный статус")
	}

	amperage_application, err := r.GetSingleAmperageApplication(id)
	if err != nil {
		return ds.AmperageApplication{}, err
	} else if amperage_application.Status != "formed" {
		return ds.AmperageApplication{}, fmt.Errorf("этот расчёт не может быть %s", status)
	}

	err = r.db.Model(&amperage_application).Updates(ds.AmperageApplication{
		Status: status,
		Finish_Date: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		Moderator_ID: uuid.NullUUID{
			UUID:  currentUserID,
			Valid: true,
		},
	}).Error
	if err != nil {
		return ds.AmperageApplication{}, err
	}

	if status == "completed" {
		amperage_applicationsDevice, err := r.GetDevicesAmperageApplications(int(amperage_application.Amperage_Application_ID))
		if err != nil {
			return ds.AmperageApplication{}, err
		}
		for _, amperage_applicationDevice := range amperage_applicationsDevice {
			device, err := r.GetDevice(int(amperage_applicationDevice.Device_ID))
			if err != nil {
				return ds.AmperageApplication{}, err
			}
			device_amperage, err := CalculateDeviceAmperage(device.Dev_Power, float64(amperage_applicationDevice.Amount))
			if err != nil {
				return ds.AmperageApplication{}, err
			}
			err = r.db.Model(&amperage_applicationDevice).Updates(ds.AmperageApplicationDevices{
				Amperage: float64(device_amperage),
			}).Error
			if err != nil {
				return ds.AmperageApplication{}, err
			}
		}
	}

	return amperage_application, nil
}
