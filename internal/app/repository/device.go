package repository

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"web_backend/internal/app/ds"
	"web_backend/internal/app/minioClient"
	"web_backend/internal/app/serializer"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (r *Repository) GetDevices() ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Order("device_id").Where("is_delete = false").Find(&devices).Error
	if err != nil {
		return nil, err
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("массив пустой")
	}

	return devices, nil
}

func (r *Repository) GetDevice(id int) (*ds.Device, error) {
	device := ds.Device{}
	err := r.db.Order("device_id").Where("device_id = ? and is_delete = ?", id, false).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("%w:  устройство  с id %d", ErrNotFound, id)
		}
		return &ds.Device{}, err
	}
	return &device, nil
}

func (r *Repository) GetDevicesByTitle(title string) ([]ds.Device, error) {
	var devices []ds.Device
	err := r.db.Order("device_id").Where("title ILIKE ? and is_delete = ?", "%"+title+"%", false).Find(&devices).Error
	if err != nil {
		return nil, err
	}
	return devices, nil
}

func (r *Repository) CreateDevice(deviceJSON serializer.DeviceJSON) (ds.Device, error) {
	device := serializer.DeviceFromJSON(deviceJSON)
	if device.Dev_Power <= 0 {
		return ds.Device{}, errors.New("неправильная мощность устройства")
	}
	err := r.db.Create(&device).First(&device).Error
	if err != nil {
		return ds.Device{}, err
	}
	return device, nil
}

func (r *Repository) EditDevice(id int, deviceJSON serializer.DeviceJSON) (ds.Device, error) {
	device := ds.Device{}
	if id < 0 {
		return ds.Device{}, errors.New("id должно быть >= 0")
	}
	err := r.db.Where("device_id = ? and is_delete = ?", id, false).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ds.Device{}, fmt.Errorf("%w: устройство с id %d", ErrNotFound, id)
		}
		return ds.Device{}, err
	}
	if deviceJSON.Dev_Power <= 0 {
		return ds.Device{}, errors.New("неправильная мощность устройства")
	}
	err = r.db.Model(&device).Updates(serializer.DeviceFromJSON(deviceJSON)).Error
	if err != nil {
		return ds.Device{}, err
	}
	return device, nil
}

func (r *Repository) DeleteDevice(id int) error {
	device := ds.Device{}
	if id < 0 {
		return errors.New("id должно быть >= 0")
	}

	err := r.db.Where("device_id = ? and is_delete = ?", id, false).First(&device).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: устройство с id %d", ErrNotFound, id)
		}
		return err
	}
	if device.Photo != "" {
		err = minio.DeleteObject(context.Background(), r.mc, minio.GetImgBucket(), device.Photo)
		if err != nil {
			return err
		}
	}

	err = r.db.Model(&ds.Device{}).Where("device_id = ?", id).Update("is_delete", true).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddToAmperageApplication(amperage_application_id int, device_id int) error {
	var device ds.Device
	if err := r.db.First(&device, device_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: устройство с id %d", ErrNotFound, device_id)
		}
		return err
	}

	var amperage_application ds.AmperageApplication
	if err := r.db.First(&amperage_application, amperage_application_id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("%w: заявка с id %d", ErrNotFound, amperage_application_id)
		}
		return err
	}
	
	amperage_applicationDevice := ds.AmperageApplicationDevices{}
	result := r.db.Where("device_id = ? and amperage_application_id = ?", device_id, amperage_application_id).Find(&amperage_applicationDevice)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected != 0 {
		return fmt.Errorf("%w: устройство %d уже в заявке %d", ErrAlreadyExists, device_id, amperage_application_id)
	}
	return r.db.Create(&ds.AmperageApplicationDevices{
		Device_ID:    uint(device_id),
		Amperage_Application_ID: uint(amperage_application_id),
		Amount: 1,
	}).Error
}

func (r *Repository) GetModeratorAndCreatorLogin(amperage_application ds.AmperageApplication) (string, string, error) {
	var creator ds.Users
	var moderator ds.Users

	err := r.db.Where("user_id = ?", amperage_application.Creator_ID).First(&creator).Error
	if err != nil {
		return "", "", err
	}

	var moderatorLogin string
	if amperage_application.Moderator_ID.Valid{
		err = r.db.Where("user_id = ?", amperage_application.Moderator_ID).First(&moderator).Error
		if err != nil {
			return "", "", err
		}
		moderatorLogin = moderator.Login
	}
	
	return creator.Login, moderatorLogin, nil
}

func (r *Repository) AddPhoto(ctx *gin.Context, device_id int, file *multipart.FileHeader) ( ds.Device, error) {
	device_, err := r.GetDevice(device_id)
	if err != nil {
		return ds.Device{}, err
	}
	
	fileName, err := minio.UploadImage(ctx, r.mc, minio.GetImgBucket(), file, *device_)
	if err != nil {
		return ds.Device{},err
	}

	device, err := r.GetDevice(device_id)
	if err != nil {
		return ds.Device{}, err
	}
	device.Photo = fileName
	err = r.db.Save(&device).Error
	if err != nil {
		return ds.Device{}, err
	}
	return *device, nil
}