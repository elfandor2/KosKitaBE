package service

import (
	"KosKita/features/admin"
	"KosKita/features/booking"
	"KosKita/features/kos"
	"KosKita/features/user"
	"errors"
)

type adminService struct {
	userData    user.UserDataInterface
	kosData     kos.KosDataInterface
	bookingData booking.BookDataInterface
	userService user.UserServiceInterface
}

// dependency injection
func New(repoUser user.UserDataInterface, repoKos kos.KosDataInterface, repoBook booking.BookDataInterface, us user.UserServiceInterface) admin.AdminServiceInterface {
	return &adminService{
		userData:    repoUser,
		kosData:     repoKos,
		bookingData: repoBook,
		userService: us,
	}
}

// GetTotalData implements admin.AdminServiceInterface.
func (as *adminService) GetTotalData(userIdLogin int) (admin.DashboardData, error) {
	user, err := as.userService.GetById(userIdLogin)
	if err != nil {
		return admin.DashboardData{}, err
	}

	if user.Role != "admin" {
		return admin.DashboardData{}, errors.New("anda bukan admin")
	}

	dashboardData := admin.DashboardData{}

	totalUser, err := as.userData.GetTotalUser()
	if err != nil {
		return dashboardData, err
	}
	dashboardData.TotalUser = totalUser

	totalKos, err := as.kosData.GetTotalKos()
	if err != nil {
		return dashboardData, err
	}
	dashboardData.TotalKos = totalKos

	totalBooking, err := as.bookingData.GetTotalBooking()
	if err != nil {
		return dashboardData, err
	}
	dashboardData.TotalBooking = totalBooking

	return dashboardData, nil
}
