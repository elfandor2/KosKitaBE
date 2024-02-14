package handler

import (
	"KosKita/features/booking"
	"KosKita/utils/middlewares"
	"KosKita/utils/responses"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type BookHandler struct {
	bookService booking.BookServiceInterface
}

func New(bs booking.BookServiceInterface) *BookHandler {
	return &BookHandler{
		bookService: bs,
	}
}

func (handler *BookHandler) CreateBook(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	newBook := BookRequest{}
	errBind := c.Bind(&newBook)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data booking not valid", nil))
	}

	bookCore := RequestToCoreBook(newBook, uint(userIdLogin))
	payment, errInsert := handler.bookService.Create(userIdLogin, bookCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	// result := CoreToResponseBook(payment)
	result := BookingResponse{}
	if payment != nil {
		result = CoreToResponseBook(payment)
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success booking kos", result))
}

func (handler *BookHandler) CancelBooking(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	bookingId := c.Param("id")
	fmt.Println(bookingId)

	updateBookingStatus := CancelBookingRequest{}
	errBind := c.Bind(&updateBookingStatus)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	bookingCore := CancelRequestToCoreBooking(updateBookingStatus)
	errCancel := handler.bookService.CancelBooking(userIdLogin, bookingId, bookingCore)
	if errCancel != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error cancel booking "+errCancel.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success cancel booking", nil))
}

func (handler *BookHandler) GetBooking(c echo.Context) error {
	idJWT := middlewares.ExtractTokenUserId(c)
	if idJWT == 0 {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("unauthorized or jwt expired", nil))
	}

	results, err := handler.bookService.GetBooking(uint(idJWT))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("Error booking. "+err.Error(), nil))
	}
	var bookingResults []BookingHistoryResponse
	for _, result := range results {
		bookingResults = append(bookingResults, CoreToResponseBookHistory(&result))
	}
	return c.JSON(http.StatusOK, responses.WebResponse("Success get booking.", bookingResults))
}

func (handler *BookHandler) WebhoocksNotification(c echo.Context) error {
	var webhoocksReq = WebhoocksRequest{}
	errBind := c.Bind(&webhoocksReq)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	bookingCore := WebhoocksRequestToCore(webhoocksReq)
	err := handler.bookService.WebhoocksData(bookingCore)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update status booking", nil))
}
