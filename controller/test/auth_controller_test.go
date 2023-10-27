package controller

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"github.com/IrvanWijayaSardam/SelfBank/controller"
	"github.com/IrvanWijayaSardam/SelfBank/dto"
	"github.com/IrvanWijayaSardam/SelfBank/entity"
	"github.com/IrvanWijayaSardam/SelfBank/service/mocks"
)

func TestAuthController_Login(t *testing.T) {
	e := echo.New()
	t.Run("Success Login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email": "test@gmail.com", "password": "password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authService := mocks.NewAuthService(t)
		jwtService := mocks.NewMockJWTService(t)

		controller := controller.NewAuthController(authService, jwtService)

		authService.On("VerifyCredential", "test@gmail.com", "password123").Return(
			entity.User{
				ID:            1,
				Email:         "aminivan@gmail.com",
				Password:      "<PASSWORD>",
				Namadepan:     "aminivan",
				Telephone:     "08123456789",
				Jk:            "Laki-Laki",
				IdRole:        10,
				AccountNumber: 123456789,
			}).Once()
		jwtService.On(
			"GenerateToken",
			strconv.FormatUint(1, 10),
			"aminivan",
			"aminivan@gmail.com",
			"08123456789",
			"Laki-Laki",
			uint64(10),
			"123456789").Return("Valid Token", nil)

		err := controller.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
	t.Run("Failed Login", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email": "test@gmail.com", "password": "password123"}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authService := mocks.NewAuthService(t)
		jwtService := mocks.NewMockJWTService(t)

		controller := controller.NewAuthController(authService, jwtService)

		authService.On("VerifyCredential", "test@gmail.com", "password123").Return(nil).Once()

		err := controller.Login(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
	t.Run("Bind Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/login", strings.NewReader(`{"email": "test@gmail.com", "password": "password123"}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		authService := mocks.NewAuthService(t)
		jwtService := mocks.NewMockJWTService(t)

		controller := controller.NewAuthController(authService, jwtService)

		err := controller.Login(c)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})
}

func TestAuthController_Register(t *testing.T) {

	var registerData = dto.RegisterDTO{
		Namadepan:    "Zelvia",
		Namabelakang: "Olga Maharani",
		Email:        "zeolga@gmail.com",
		Username:     "zeolga",
		Password:     "zeolga",
		Telephone:    "6281340691423",
		Jk:           "F",
		Status:       1,
		IdRole:       2,
	}

	var dataUser = entity.User{
		ID:            1,
		Email:         "aminivan@gmail.com",
		Password:      "<PASSWORD>",
		Namadepan:     "aminivan",
		Telephone:     "08123456789",
		Jk:            "Laki-Laki",
		IdRole:        2,
		AccountNumber: 123456789,
	}

	authService := mocks.NewAuthService(t)
	jwtService := mocks.NewMockJWTService(t)

	e := echo.New()
	t.Run("Success Register", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{
			"nama_depan": "Zelvia",
			"nama_belakang": "Olga Maharani",
			"email": "zeolga@gmail.com",
			"username": "zeolga",
			"password": "zeolga",
			"telp": "6281340691423",
			"jk": "F",
			"idrole": 2,
			"status": 1
		}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller := controller.NewAuthController(authService, jwtService)

		authService.On("IsDuplicateEmail", "zeolga@gmail.com").Return(true).Once()
		authService.On("CreateUser", registerData).Return(dataUser).Once()
		jwtService.On(
			"GenerateToken",
			"1",
			"aminivan",
			"aminivan@gmail.com",
			"08123456789",
			"Laki-Laki",
			uint64(2),
			"123456789").Return("Valid Token", nil)

		err := controller.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("Bind Error", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{
			"nama_depan": "Zelvia",
			"nama_belakang": "Olga Maharani",
			"email": "zeolga@gmail.com",
			"username": "zeolga",
			"password": "zeolga",
			"telp": "6281340691423",
			"jk": "F",
			"idrole": 2,
			"status": 1
		}`))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller := controller.NewAuthController(authService, jwtService)

		err := controller.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

	t.Run("Duplicate Email", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/api/auth/register", strings.NewReader(`{
			"nama_depan": "Zelvia",
			"nama_belakang": "Olga Maharani",
			"email": "zeolga@gmail.com",
			"username": "zeolga",
			"password": "zeolga",
			"telp": "6281340691423",
			"jk": "F",
			"idrole": 2,
			"status": 1
		}`))
		req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		controller := controller.NewAuthController(authService, jwtService)

		authService.On("IsDuplicateEmail", "zeolga@gmail.com").Return(false).Once()

		err := controller.Register(c)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusConflict, rec.Code)

		authService.AssertExpectations(t)
		jwtService.AssertExpectations(t)
	})

}
