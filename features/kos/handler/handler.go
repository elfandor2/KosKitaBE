package handler

import (
	"KosKita/features/kos"
	"KosKita/utils/cloudinary"
	"KosKita/utils/middlewares"
	"KosKita/utils/responses"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type KosHandler struct {
	kosService kos.KosServiceInterface
	cld        cloudinary.CloudinaryUploaderInterface
}

func New(ks kos.KosServiceInterface, cloudinaryUploader cloudinary.CloudinaryUploaderInterface) *KosHandler {
	return &KosHandler{
		kosService: ks,
		cld:        cloudinaryUploader,
	}
}

func (handler *KosHandler) CreateKos(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	newKos := KosRequest{}
	errBind := c.Bind(&newKos)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid "+errBind.Error(), nil))
	}

	kosCore := RequestToCore(newKos, uint(userIdLogin))

	kosId, errInsert := handler.kosService.Create(userIdLogin, kosCore)
	if errInsert != nil {
		if errInsert.Error() == "anda bukan owner" {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse(errInsert.Error(), nil))
		}
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	response := KosIdResponse{
		ID: strconv.Itoa(int(kosId)),
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success create kos", response))
}

func (handler *KosHandler) UploadImages(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kosID, err := strconv.Atoi(c.Param("kosid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing kos id", nil))
	}

	newFoto := KosFotoRequest{}
	errBind := c.Bind(&newFoto)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	var imageUrls []string
	photoFields := []string{"main_kos_photo", "front_kos_photo", "back_kos_photo", "front_room_photo", "inside_room_photo"}

	for _, field := range photoFields {
		fileHeader, err := c.FormFile(field)
		if err != nil {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("masukan semua foto", nil))
		}

		imageURL, err := handler.cld.UploadImage(fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error uploading the image", nil))
		}

		imageUrls = append(imageUrls, imageURL)
	}

	kosCore := RequestToCoreFoto(imageUrls, uint(userIdLogin))

	errInsert := handler.kosService.CreateImage(userIdLogin, kosID, kosCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error upload image -> "+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success upload image", nil))
}

func (handler *KosHandler) UpdateImages(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kosID, err := strconv.Atoi(c.Param("kosid"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing kos id", nil))
	}

	newFoto := KosFotoRequest{}
	errBind := c.Bind(&newFoto)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	var imageUrls []string
	photoFields := []string{"main_kos_photo", "front_kos_photo", "back_kos_photo", "front_room_photo", "inside_room_photo"}

	for _, field := range photoFields {
		fileHeader, err := c.FormFile(field)
		if err != nil {
			continue
		}

		imageURL, err := handler.cld.UploadImage(fileHeader)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, responses.WebResponse("error uploading the image", nil))
		}

		imageUrls = append(imageUrls, imageURL)
	}

	kosCore := RequestToCoreFotoPut(imageUrls, uint(userIdLogin))

	errInsert := handler.kosService.CreateImage(userIdLogin, kosID, kosCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error upload image"+errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success upload image", nil))
}

func (handler *KosHandler) UpdateKos(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kosID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing kos id", nil))
	}

	updateKos := KosRequest{}
	errBind := c.Bind(&updateKos)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	kosCore := RequestToCorePut(updateKos, uint(userIdLogin))
	kosCore.ID = uint(kosID)

	errUpdate := handler.kosService.Put(userIdLogin, kosCore)
	if errUpdate != nil {
		if errUpdate.Error() == "anda bukan owner" {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse(errUpdate.Error(), nil))
		}

		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errUpdate.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success update kos", nil))
}

func (handler *KosHandler) CreateRating(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kosId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error parsing kos id", nil))
	}

	newRating := RatingRequest{}
	errBind := c.Bind(&newRating)
	if errBind != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("error bind data. data not valid", nil))
	}

	ratingCore := RequestToCoreRating(newRating, uint(userIdLogin), uint(kosId))

	errInsert := handler.kosService.CreateRating(userIdLogin, kosId, ratingCore)
	if errInsert != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errInsert.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success rating score", nil))
}

func (handler *KosHandler) GetKosByRating(c echo.Context) error {
	kos, err := handler.kosService.GetByRating()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get data", nil))
	}

	var kosResponse []interface{}
	for _, k := range kos {
		kosResponse = append(kosResponse, CoreToGetRating(k))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success get kos", kosResponse))
}

func (handler *KosHandler) DeleteKos(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kosId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("id kos kosong", nil))
	}

	errDelete := handler.kosService.Delete(userIdLogin, kosId)
	if errDelete != nil {
		if errDelete.Error() == "kos id tidak ada" {
			return c.JSON(http.StatusNotFound, responses.WebResponse(errDelete.Error(), nil))
		}

		if errDelete.Error() == "kos ini bukan milik Anda" {
			return c.JSON(http.StatusUnauthorized, responses.WebResponse(errDelete.Error(), nil))
		}

		return c.JSON(http.StatusInternalServerError, responses.WebResponse(errDelete.Error(), nil))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success delete kos", nil))
}

func (handler *KosHandler) GetKosById(c echo.Context) error {
	kosIdStr := c.Param("id")
	if kosIdStr == "" {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("id kos kosong", nil))
	}

	kosId, err := strconv.Atoi(kosIdStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, responses.WebResponse("id kos salah", nil))
	}

	kos, err := handler.kosService.GetById(kosId)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get data", nil))
	}

	var kosResult = CoreToGetDetail(*kos)
	return c.JSON(http.StatusOK, responses.WebResponse("success read data.", kosResult))
}

func (handler *KosHandler) GetKosByUserId(c echo.Context) error {
	userIdLogin := middlewares.ExtractTokenUserId(c)
	if userIdLogin == 0 {
		return c.JSON(http.StatusUnauthorized, responses.WebResponse("Unauthorized user", nil))
	}

	kos, err := handler.kosService.GetByUserId(userIdLogin)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse("error get kos", nil))
	}

	var kosResponse []interface{}
	for _, k := range kos {
		kosResponse = append(kosResponse, CoreToGetUser(k))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success get kos", kosResponse))
}

func (handler *KosHandler) SearchKos(c echo.Context) error {
	expectedParams := map[string]bool{
		"address":  true,
		"category": true,
		"minPrice": true,
		"maxPrice": true,
	}

	for param := range c.QueryParams() {
		if _, ok := expectedParams[param]; !ok {
			return c.JSON(http.StatusBadRequest, responses.WebResponse("invalid search param", nil))
		}
	}

	query := c.QueryParam("address")
	category := c.QueryParam("category")
	minPrice, _ := strconv.Atoi(c.QueryParam("minPrice"))
	maxPrice, _ := strconv.Atoi(c.QueryParam("maxPrice"))

	kos, err := handler.kosService.SearchKos(query, category, minPrice, maxPrice)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, responses.WebResponse(err.Error(), nil))
	}

	var kosResponse []interface{}
	for _, k := range kos {
		kosResponse = append(kosResponse, CoreToGetRating(k))
	}

	return c.JSON(http.StatusOK, responses.WebResponse("success get kos", kosResponse))
}
