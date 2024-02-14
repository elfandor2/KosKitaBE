package handler

import "KosKita/features/admin"

type DashboardData struct {
	TotalUser           int `json:"total_user"`
	TotalBooking        int `json:"total_booking"`
	TotalKos            int `json:"total_kos"`
}

func CoreToResponseDashboard(core *admin.DashboardData) DashboardData {
	return DashboardData{
		TotalUser:           core.TotalUser,
		TotalBooking:        core.TotalBooking,
		TotalKos:            core.TotalKos,
	}
}
