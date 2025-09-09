package repository

import (
	"fmt"
	"strings"
)

type Repository struct {
}

func NewRepository() (*Repository, error) {
	return &Repository{}, nil
}

type Device struct {
	ID int
	Title string
	Power float64
	Photo string
	Description string
}

func (r *Repository) GetDevices() ([]Device, error) {
	devices := []Device{
		{
			ID: 1,
			Title: "Ноутбук HUAWEI MateBook D 16 MCLG-X",
			Power: 0.065,
			Photo: "photo_53144771654101б20689_x.jpg",
			Description: "Мощный ноутбук, подходит для компьютерных игр и работы",
		},
		{
			ID: 2, 
			Title: "Холодильник Candy",
			Power: 0.2,
			Photo: "photo_5314477165410120685_x.jpg",
			Description: "Холодильник от китайской фирмы с морозильной камерой",
		},
		{
			ID: 3, 
			Title: "Смартфон Samsung Galaxy A16",
			Power: 0.025,
			Photo: "photo_5314477165410120691_x.jpg",
			Description: "Новая модель телефона Samsung в линейке A с отличной камерой",
		},
		{
			ID: 4,
			Title: "Телевизор Hisense 55E7Q PRO",
			Power: 0.13,
			Photo: "photo_5314477165410120687_x.jpg",
			Description: "Широкоформатный телевизор с большой диагональю и трёхканальным звуком",
		},
		{
			ID: 5, 
			Title: "Телевизор TCL 55P8K",
			Power: 0.15,
			Photo: "photo_5314477165410120686_x.jpg",
			Description: "Флагманский телевизор бренда TCL, отлично подходит как для кино, так и для игр",
		},
		{
			ID: 6, 
			Title: "Электрический духовой шкаф Gorenje",
			Power: 3.5,
			Photo: "photo_5314477165410120693_x.jpg",
			Description: "Духовой шкаф от популярной вьетнамской фирмы с подключением по Wi-fi",
		},
		{
			ID: 7,
			Title: "Смартфон Xiaomi Redmi Note 14",
			Power: 0.033,
			Photo: "photo_5314477165410120683_x.jpg",
			Description: "Бюджетная новинка от всеми любимого китайского бренда",
		},
		{
			ID: 8, 
			Title: "Кофемашина автоматическая Krups",
			Power: 1.45,
			Photo: "photo_5314477165410120690_x.jpg",
			Description: "Кофемашина с интеграцией в умный дом и Алису",
		},
		{
			ID: 9, 
			Title: "Телевизор Haier 50 Smart TV AX Pro",
			Power: 0.16,
			Photo: "photo_5314477165410120688_x.jpg",
			Description: "УЦЕНКА: небольшой скол на задней панели и подставке",
		},
		{
			ID: 10,
			Title: "Наушники True Wireless HUAWEI",
			Power: 0.015,
			Photo: "photo_5314477165410120692_x.jpg",
			Description: "Бюджетные наушники с отличным звуком, хит продаж!",
		},
		{
			ID: 11, 
			Title: "Аэрогриль Tefal Easy Fry & Grill Digital",
			Power: 1.55,
			Photo: "photo_5314477165410120684_x.jpg",
			Description: "Мощный аэрогриль, подойдёт для дачи и дома",
		},
		{
			ID: 12, 
			Title: "Смарт-часы HUAWEI Watch GT 5",
			Power: 0.017,
			Photo: "photo_5314477165410120682_x.jpg",
			Description: "Смарт-часы для бега и повседневной носки",
		},
	}
	if len(devices) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return devices, nil
}

func (r *Repository) GetDevice(id int) (Device, error) {
	devices, err := r.GetDevices()
	if err != nil {
		return Device{}, err
	}

	for _, device := range devices {
		if device.ID == id {
			return device, nil
		}
	}
	return Device{}, fmt.Errorf("Заказ не найден" )
}

func (r *Repository) GetDeviceByTitle(title string) ([]Device, error) {
	devices, err := r.GetDevices()
	if err != nil {
		return []Device{}, err
	}

	var result []Device
	for _, device := range devices {
		if strings.Contains(strings.ToLower(device.Title), strings.ToLower(title)) {
			result = append(result, device)
		}
	}
	return result, nil
}

func (r *Repository) GetCart() ([]Device, error) {
	devices, err := r.GetDevices()
	if err != nil {
		return []Device{}, err
	} 
	
	var result []Device
	for _, device := range devices {
		if device.ID == 4 || device.ID == 7 {
			result = append(result, device)
		}
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return result, nil
}
