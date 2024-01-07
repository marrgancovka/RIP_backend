package role

type Role string

const (
	Buyer   Role = "client"  // 0
	Manager Role = "manager" // 1
	Admin   Role = "admin"   // 2
)
