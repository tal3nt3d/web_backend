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

func (r *Repository) GetAmperageApplicationCount() int64 {
   var AmperageApplicationID uint
   var count int64
   creatorID := 1

   err := r.db.Model(&ds.AmperageApplication{}).Where("creator_id = ? AND status = ?", creatorID, "draft").Select("amperage_application_id").First(&AmperageApplicationID).Error
   if err != nil {
      return 0
   }

   err = r.db.Model(&ds.AmperageApplicationDevices{}).Where("amperage_application_id = ?", AmperageApplicationID).Count(&count).Error
   if err != nil {
      logrus.Println("Error counting records in lists_chats:", err)
   }

   return count
}

func (r *Repository) GetActiveAmperageApplicationID() uint {
	var AmperageApplicationID uint
	err := r.db.Model(&ds.AmperageApplication{}).Where("status = ?", "draft").Select("amperage_application_id").First(&AmperageApplicationID).Error
	if err != nil {
		return 0
	}
	return AmperageApplicationID
}

func (r *Repository) GetAmperageApplication(id int) ([]ds.AmperageApplicationDevices, error) {
    var amperageApplicationItems []ds.AmperageApplicationDevices
    err := r.db.Where("amperage_application_id = ?", id).Preload("Device").Find(&amperageApplicationItems).Error
    if err != nil {
        return nil, err
    }

    return amperageApplicationItems, nil
}

func (r *Repository) AddDevice(deviceID uint, creatorID uint) (error) {
    var app ds.AmperageApplication

    err := r.db.Where("creator_id = ? AND status = ?", creatorID, "draft").
        First(&app).Error

    if errors.Is(err, gorm.ErrRecordNotFound) {
        app = ds.AmperageApplication{
            Status:     "draft",
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
    r.db.Model(&ds.AmperageApplicationDevices{}).
        Where("amperage_application_id = ? AND device_id = ?", app.Amperage_Application_ID, deviceID).Preload("Device").
        Count(&count)

    if count == 0 {
		var device ds.Device
		  if err := r.db.First(&device, deviceID).Error; err != nil {
            return err
        }

        appDev := ds.AmperageApplicationDevices{
            Amperage_Application_ID: app.Amperage_Application_ID,
            Device_ID:      deviceID,
			Amperage: 		0,
            Amount:         1,
        }
        if err := r.db.Create(&appDev).Error; err != nil {
            return err
        }
    }

    return nil
}

func (r *Repository) DeleteAmperageApplication(appID uint) error {
	query := `
		UPDATE amperage_applications 
		SET status = 'deleted'
		WHERE amperage_application_id = $1;
	`
	result := r.db.Exec(query, appID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("amperage_application with id %d not found", appID)
	}
	return nil
}

func (r *Repository) IsDraftAmperageApplication(appID int) (bool, error) {
	var app ds.AmperageApplication
	err := r.db.Select("status").Where("amperage_application_id = ?", appID).First(&app).Error
	if err != nil {
		return false, err
	}
	return app.Status == "draft", nil
}
