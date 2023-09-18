package ds

type Ship struct {
	ID          uint `gorm:"primarykey"`
	Name        string
	Image       string
	Description string
}
