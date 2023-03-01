package utility

import (
	"os"

	"github.com/Yothgewalt/aufruf-jaeger-bot/module/controller"
	"github.com/Yothgewalt/aufruf-jaeger-bot/module/repository"
	"github.com/rs/zerolog"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

const CREDENTIALS_FILENAME string = "credentials.json"

func NewClassroomService(logger zerolog.Logger) *classroom.Service {
	b, err := os.ReadFile(CREDENTIALS_FILENAME)
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to read credentials file.")
	}

	googleClient := repository.NewGoogleClientRepository(logger)
	httpClient := googleClient.GetClient(b, classroom.ClassroomCoursesReadonlyScope)

	googleService := controller.NewGoogleClassroomController(logger)
	classroomService, err := googleService.NewGoogleClassroomService(option.WithHTTPClient(httpClient))
	if err != nil {
		logger.Fatal().Err(err).Msg("unable to create classroom client.")
	}

	return classroomService
}
