package ds

type ShipToAppReq struct {
	ShipId uint `json:"id_ship"`
}

type AppStatus struct {
	Id     uint   `json:"id"`
	Status string `json:"status"`
}
