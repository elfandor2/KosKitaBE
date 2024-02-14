package handler

import (
	"KosKita/features/order"
	"KosKita/utils/middlewares"
	"KosKita/utils/responses"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
)

// _midtransController "BE-REPO-20/features/midtrans/controller"

type OrderHandler struct {
	orderService order.OrderServiceInterface
}

func NewOrder(service order.OrderServiceInterface) *OrderHandler {
	return &OrderHandler{
		orderService: service,
	}
}

func (handler *OrderHandler) CreateOrder(c echo.Context) error {
	idJWT := middlewares.ExtractTokenUserId(c)
	if idJWT == 0 {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("unauthorized or jwt expired", nil))
	}

	newOrder := OrderRequest{}
	errBind := c.Bind(&newOrder)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data order not valid", nil))
	}

	orderCore := RequestToCoreOrder(newOrder)
	payment, errInsert := handler.orderService.PostOrder(uint(idJWT), orderCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error insert order "+errInsert.Error(), nil))
	}

	results := CoreToResponse(*payment)

	return c.JSON(http.StatusOK, responses.WebResponse("Success get order.", results))
}

func (handler *OrderHandler) GetOrders(c echo.Context) error {
	idJWT := middlewares.ExtractTokenUserId(c)
	if idJWT == 0 {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("unauthorized or jwt expired", nil))
	}

	results, err := handler.orderService.GetOrders(uint(idJWT))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error order. "+err.Error(), nil))
	}

	var bookingResults []OrderHistoryResponse
	for _, result := range results {
		bookingResults = append(bookingResults, CoreToResponseOrderHistory(&result))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("Success get order.", bookingResults))
}

func (handler *OrderHandler) CancelOrderById(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	orderId := c.Param("order_id")

	updateOrderStatus := CancelOrderRequest{}
	errBind := c.Bind(&updateOrderStatus)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}
	fmt.Println("order id::", orderId)
	orderCore := CancelRequestToCoreOrder(updateOrderStatus)
	errCancel := handler.orderService.CancelOrder(userIdLogin, orderId, orderCore)
	if errCancel != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error cancel order "+errCancel.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success cancel order", nil))
}

func (handler *OrderHandler) WebhoocksNotification(c echo.Context) error {

	var webhoocksReq = WebhoocksRequest{}
	errBind := c.Bind(&webhoocksReq)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	orderCore := WebhoocksRequestToCore(webhoocksReq)
	err := handler.orderService.WebhoocksService(orderCore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error Notif "+err.Error(), nil))
	}

	log.Println("transaction success")
	return c.JSON(http.StatusOK, responses.WebResponse("transaction success", nil))
}
