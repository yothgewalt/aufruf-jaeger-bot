package controller

import (
	"context"

	"github.com/Yothgewalt/aufruf-jaeger-bot/module/entity"
	"github.com/rs/zerolog"
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

type googleClassroom struct {
	logger zerolog.Logger
}

// Initial a function for export functions that contains in the interface in this function return.
func NewGoogleClassroomController(loggerSource zerolog.Logger) entity.GoogleClassroom {
	return &googleClassroom{
		logger: loggerSource,
	}
}

// Export classroom services from google authorization
func (gcr *googleClassroom) NewGoogleClassroomService(opts ...option.ClientOption) (*classroom.Service, error) {
	googleService, err := classroom.NewService(context.Background(), opts...)
	if err != nil {
		return nil, err
	}

	return googleService, nil
}
