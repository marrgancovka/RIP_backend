package ds

type Users struct {
	ID           uint   `gorm:"primarykey"`
	FirstName    string `gorm:"column:firstname"`
	SecondName   string `gorm:"column:secondname"`
	Phone        string `gorm:"column:phone"`
	UserName     string `gorm:"column:username"`
	UserPassword string `gorm:"column:userpassword"`
	UserRole     string `gorm:"column:userrole"`
}
