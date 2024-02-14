package admin

type DashboardData struct {
	TotalUser int 
	TotalBooking int 
	TotalKos int
	TotalBookingPerMonth int
}


// interface untuk Service Layer
type AdminServiceInterface interface {
	GetTotalData(userIdLogin int)(DashboardData, error)
}
