package handler

import (
	"KosKita/features/admin"
	"KosKita/utils/middlewares"
	"KosKita/utils/responses"
	"net/http"

	"github.com/labstack/echo/v4"
)

type AdminHandler struct {
	adminService admin.AdminServiceInterface
}

func New(as admin.AdminServiceInterface) *AdminHandler {
	return &AdminHandler{
		adminService: as,
	}
}

func (handler *AdminHandler) GetAllData(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	dashboardData, errGet := handler.adminService.GetTotalData(userIdLogin)
	if errGet != nil {
		if errGet.Error() == "anda bukan admin" {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse(errGet.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errGet.Error(), nil))
	}

	dashboardResult := CoreToResponseDashboard(&dashboardData)

	return c.JSON(http.StatusOK, responses.WebResponse("success get dashboard data", dashboardResult))
}