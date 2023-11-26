package ds

type Cosmodroms struct {
	ID      uint `gorm:"primarykey"`
	Title   string
	City    string
	Country string
}
