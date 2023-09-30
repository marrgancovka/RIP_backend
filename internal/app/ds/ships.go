package ds

type Ship struct {
	ID          uint `gorm:"primarykey"`
	Title       string
	Rocket      string
	Type        string
	Description string
	Image_url   string
	Is_delete   bool
}
