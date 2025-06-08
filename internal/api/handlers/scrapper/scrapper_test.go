package scrapper_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/es-debug/backend-academy-2024-go-template/internal/api/handlers/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostTgChatId(t *testing.T) {
	mockUserService := new(mocks.UserServiceMock)
	mockLinksService := new(mocks.LinksServiceMock)
	api := scrapper.NewAPI(mockUserService, mockLinksService)

	mockUserService.On("RegisterUser", mock.Anything, int64(123)).Return(nil)

	req, _ := http.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rr := httptest.NewRecorder()
	api.PostTgChatId(rr, req, 123)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "123\n", rr.Body.String())
	mockUserService.AssertExpectations(t)
}

func TestPostTgChatId_UserAlreadyExists(t *testing.T) {
	mockUserService := new(mocks.UserServiceMock)
	mockLinksService := new(mocks.LinksServiceMock)
	api := scrapper.NewAPI(mockUserService, mockLinksService)

	mockUserService.On("RegisterUser", mock.Anything, int64(123)).Return(&errors.ErrUserAlreadyExists{})

	req, _ := http.NewRequest(http.MethodPost, "/tg-chat/123", http.NoBody)
	rr := httptest.NewRecorder()
	api.PostTgChatId(rr, req, 123)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "already exists")
	mockUserService.AssertExpectations(t)
}

func TestDeleteTgChatId(t *testing.T) {
	mockUserService := new(mocks.UserServiceMock)
	mockLinksService := new(mocks.LinksServiceMock)
	api := scrapper.NewAPI(mockUserService, mockLinksService)

	mockUserService.On("DeleteUser", mock.Anything, int64(123)).Return(nil)

	req, _ := http.NewRequest(http.MethodDelete, "/tg-chat/123/links", http.NoBody)
	rr := httptest.NewRecorder()
	api.DeleteTgChatId(rr, req, 123)

	assert.Equal(t, http.StatusOK, rr.Code)
	mockLinksService.AssertExpectations(t)
}

func TestDeleteTgChatId_UserNotFound(t *testing.T) {
	mockUserService := new(mocks.UserServiceMock)
	mockLinksService := new(mocks.LinksServiceMock)
	api := scrapper.NewAPI(mockUserService, mockLinksService)

	mockUserService.On("DeleteUser", mock.Anything, int64(123)).Return(&errors.ErrUserNotFound{})

	req, _ := http.NewRequest(http.MethodDelete, "/tg-chat/123/links", http.NoBody)
	rr := httptest.NewRecorder()
	api.DeleteTgChatId(rr, req, 123)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.Contains(t, rr.Body.String(), "not found")
	mockUserService.AssertExpectations(t)
}
