package ds

import "time"

type Flights struct {
	Id_Ship            uint
	Id_Application     uint
	Id_Cosmodrom_Begin uint
	Id_cosmodrom_End   uint
	Date_Flight        time.Time
}
