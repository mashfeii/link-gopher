package application_test

// import (
// 	"io"
// 	"net/http"
// 	"strings"
// 	"testing"
// 	"time"

// 	bot_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/bot"
// 	"github.com/es-debug/backend-academy-2024-go-template/internal/application"
// 	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
// 	"github.com/es-debug/backend-academy-2024-go-template/internal/mocks"
// 	"github.com/stretchr/testify/mock"
// )

// func TestCheckUpdates(t *testing.T) {
// 	mockLinkRepo := mocks.NewLinkRepositoryMock(t)
// 	mockBotClient := mocks.NewClientInterfaceMock(t)
// 	mockGhClient := mocks.NewLinkCheckerMock(t)
// 	mockSoClient := mocks.NewLinkCheckerMock(t)

// 	links := []models.Link{
// 		{ChatID: 123, URL: "https://github.com/golang/go", LastUpdated: time.Now().Add(-time.Hour)},
// 		{ChatID: 123, URL: "https://stackoverflow.com/questions/123", LastUpdated: time.Now().Add(-time.Hour)},
// 	}
// 	mockLinkRepo.On("GetAllActiveLinks", mock.Anything).Return(links, nil)

// 	ghEvent := mocks.NewEventMock(t)
// 	ghEvent.On("GetType").Return("pull_request")
// 	ghEvent.On("GetUser").Return("user")
// 	ghEvent.On("GetTitle").Return("New Pull Request")
// 	ghEvent.On("GetCreatedAt").Return(time.Now().Add(-time.Minute))
// 	ghEvent.On("GetBody").Return("Description of the PR")
// 	ghUpdates := []models.Event{ghEvent}
// 	mockGhClient.On("GetUpdates", "https://github.com/golang/go", time.Now()).Return(ghUpdates, nil).Once()

// 	soEvent := mocks.NewEventMock(t)
// 	soEvent.On("GetType").Return("answer")
// 	soEvent.On("GetUser").Return("stackoverflow_user")
// 	soEvent.On("GetTitle").Return("Question Title")
// 	soEvent.On("GetCreatedAt").Return(time.Now().Add(-time.Minute))
// 	soEvent.On("GetBody").Return("This is an answer to the question.")
// 	soUpdates := []models.Event{soEvent}
// 	mockSoClient.On("GetUpdates", "https://stackoverflow.com/questions/123", time.Now()).Return(soUpdates, nil).Once()
// 	mockSoClient.On("GetQuestionTitle", "https://stackoverflow.com/questions/123").Return("Question Title", nil).Once()

// 	description := mock.Anything

// 	mockBotClient.On("PostUpdates", mock.Anything, bot_client.LinkUpdate{
// 		TgChatId:    &links[0].ChatID,
// 		Url:         &links[0].URL,
// 		Description: &description,
// 	}).Return(&http.Response{
// 		StatusCode: http.StatusOK,
// 		Body:       io.NopCloser(strings.NewReader("")),
// 	}, nil)
// 	mockBotClient.On("PostUpdates", mock.Anything, bot_client.LinkUpdate{
// 		TgChatId:    &links[1].ChatID,
// 		Url:         &links[1].URL,
// 		Description: &description,
// 	}).Return(&http.Response{
// 		StatusCode: http.StatusOK,
// 		Body:       io.NopCloser(strings.NewReader("")),
// 	}, nil)

// 	application.CheckUpdates(mockLinkRepo, mockBotClient, mockGhClient, mockSoClient)

// 	mockLinkRepo.AssertExpectations(t)
// 	mockGhClient.AssertExpectations(t)
// 	mockSoClient.AssertExpectations(t)
// 	mockBotClient.AssertExpectations(t)
// }
