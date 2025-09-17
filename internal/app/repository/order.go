package repository

import (
	"errors"
	"fmt"
	"time"
	"web_backend/internal/app/ds"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func (r *Repository) GetDevices() ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Find(&devices).Error
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return devices, nil
}

func (r *Repository) GetDevice(id int) (ds.Device, error) {
	device := ds.Device{}
	err := r.db.Where("device_id = ?", id).Find(&device).Error
	if err != nil {
		return ds.Device{}, err
	}
	return device, nil
}

func (r *Repository) GetDevicesByTitle(title string) ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Where("title ILIKE ?", "%"+title+"%").Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) GetCartCount() int64 {
   var ApplicationID uint
   var count int64
   creatorID := 1

   err := r.db.Model(&ds.Application{}).Where("creator_id = ? AND status = ?", creatorID, "черновик").Select("application_id").First(&ApplicationID).Error
   if err != nil {
      return 0
   }

   err = r.db.Model(&ds.ApplicationDevices{}).Where("application_id = ?", ApplicationID).Count(&count).Error
   if err != nil {
      logrus.Println("Error counting records in lists_chats:", err)
   }

   return count
}

func (r *Repository) GetActiveApplicationID() uint {
	var ApplicationID uint
	err := r.db.Model(&ds.Application{}).Where("status = ?", "черновик").Select("application_id").First(&ApplicationID).Error
	if err != nil {
		return 0
	}
	return ApplicationID
}

func (r *Repository) GetCart(id int) ([]ds.ApplicationDevices, error) {
    var cartItems []ds.ApplicationDevices
    err := r.db.Where("application_id = ?", id).Preload("Device").Find(&cartItems).Error
    if err != nil {
        return nil, err
    }

    return cartItems, nil
}

func (r *Repository) MarkApplicationDeleted(appID uint) error {
    return r.db.Model(&ds.Application{}).
        Where("application_id = ?", appID).
        Update("status", "удалён").Error
}

func (r *Repository) AddDevice(deviceID uint, creatorID uint) (error) {
    var app ds.Application

    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "черновик").
        First(&app).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        app = ds.Application{
            Status:     "черновик",
            Created_At: time.Now(),
            Creator_ID: creatorID,
			Moderator_ID: 2,
        }
        if err := r.db.Create(&app).Error; err != nil {
            return err
        }
    } else if err != nil {
        return err
    }

    var count int64
    r.db.Model(&ds.ApplicationDevices{}).
        Where("application_id = ? AND device_id = ?", app.Application_ID, deviceID).Preload("Device").
        Count(&count)

    if count == 0 {
		var device ds.Device
		  if err := r.db.First(&device, deviceID).Error; err != nil {
            return err
        }

		Dev_Power := device.Dev_Power
        appDev := ds.ApplicationDevices{
            Application_ID: app.Application_ID,
            Device_ID:      deviceID,
			Amperage: 		Dev_Power*1000/220,
            Amount:         1,
        }
        if err := r.db.Create(&appDev).Error; err != nil {
            return err
        }
    }

    return nil
}

func (r *Repository) DeleteApplication(appID uint) error {
	query := `
		UPDATE applications 
		SET status = 'удалён'
		WHERE application_id = $1;
	`
	result := r.db.Exec(query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("application with id %d not found", appID)
	}
	return nil
}


func (r *Repository) IsDraftApplication(appID int) (bool, error) {
	var app ds.Application
	err := r.db.Select("status").Where("application_id = ?", appID).First(&app).Error
	if err != nil {
		return false, err
	}
	return app.Status == "черновик", nil
}
