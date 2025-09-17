package ds

type Users struct {
	User_ID     uint   `gorm:"primary_key" json:"id"`
	Login       string `gorm:"type:varchar(20);unique;not null" json:"login"`
	Password    string `gorm:"type:varchar(20);not null" json:"-"`
	IsModerator bool   `gorm:"type:boolean;default:false" json:"is_moderator"`
}