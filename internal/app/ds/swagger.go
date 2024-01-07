package ds

type ShipToAppReq struct {
	ShipId uint `json:"id_ship"`
}

type AppStatus struct {
	Id     uint   `json:"id"`
	Status string `json:"status"`
}

type delete_flight struct {
	IdShip uint `json: id_ship`
	IdApp  uint `json: id_app`
}

type NewShip struct {
	Title       string
	Rocket      string
	Type        string
	Description string
	Image_url   string
}
