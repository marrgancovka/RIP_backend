package ds

type Users struct {
	ID           uint   `gorm:"primarykey"`
	UserName     string `gorm:"column:username"`
	UserPassword string `gorm:"column:userpassword"`
	UserRole     string `gorm:"column:userrole"`
}
