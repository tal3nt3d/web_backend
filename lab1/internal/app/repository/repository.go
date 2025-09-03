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

type Order struct {
	ID int
	Title string
	Price int
	Photo string
}

func (r *Repository) GetOrders() ([]Order, error) {
	orders := []Order{
		{
			ID: 1,
			Title: "Ноутбук HUAWEI MateBook D 16 MCLG-X",
			Price: 64999,
			Photo: "photo_5314477165410120689_x.jpg",
		},
		{
			ID: 2, 
			Title: "Холодильник Candy",
			Price: 37999,
			Photo: "photo_5314477165410120685_x.jpg",
		},
		{
			ID: 3, 
			Title: "Смартфон Samsung Galaxy A16",
			Price: 19999,
			Photo: "photo_5314477165410120691_x.jpg",
		},
		{
			ID: 4,
			Title: "Телевизор Hisense 55E7Q PRO",
			Price: 47999,
			Photo: "photo_5314477165410120687_x.jpg",
		},
		{
			ID: 5, 
			Title: "Телевизор TCL 55P8K",
			Price: 46999,
			Photo: "photo_5314477165410120686_x.jpg",
		},
		{
			ID: 6, 
			Title: "Электрический духовой шкаф Gorenje",
			Price: 45999,
			Photo: "photo_5314477165410120693_x.jpg",
		},
		{
			ID: 7,
			Title: "Смартфон Xiaomi Redmi Note 14",
			Price: 21999,
			Photo: "photo_5314477165410120683_x.jpg",
		},
		{
			ID: 8, 
			Title: "Кофемашина автоматическая Krups",
			Price: 39999,
			Photo: "photo_5314477165410120690_x.jpg",
		},
		{
			ID: 9, 
			Title: "Телевизор Haier 50 Smart TV AX Pro",
			Price: 44999,
			Photo: "photo_5314477165410120688_x.jpg",
		},
		{
			ID: 10,
			Title: "Наушники True Wireless HUAWEI",
			Price: 2999,
			Photo: "photo_5314477165410120692_x.jpg",
		},
		{
			ID: 11, 
			Title: "Аэрогриль Tefal Easy Fry & Grill Digital",
			Price: 16999,
			Photo: "photo_5314477165410120684_x.jpg",
		},
		{
			ID: 12, 
			Title: "Смарт-часы HUAWEI Watch GT 5",
			Price: 15999,
			Photo: "photo_5314477165410120682_x.jpg",
		},
	}
	if len(orders) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return orders, nil
}

func (r *Repository) GetOrder(id int) (Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return Order{}, err
	}

	for _, order := range orders {
		if order.ID == id {
			return order, nil
		}
	}
	return Order{}, fmt.Errorf("Заказ не найден" )
}

func (r *Repository) GetOrderByTitle(title string) ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	}

	var result []Order
	for _, order := range orders {
		if strings.Contains(strings.ToLower(order.Title), strings.ToLower(title)) {
			result = append(result, order)
		}
	}
	return result, nil
}

func (r *Repository) GetCart() ([]Order, error) {
	orders, err := r.GetOrders()
	if err != nil {
		return []Order{}, err
	} 
	
	var result []Order
	for _, order := range orders {
		if order.ID == 4 || order.ID == 7 {
			result = append(result, order)
		}
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("Массив пустой")
	}

	return result, nil
}
