package ds

type Users struct {
	ID           uint `gorm: "primarykey"`
	FirstName    string
	SecondName   string
	Phone        string
	UserName     string
	UserPassword string
	UserRole     string
}
