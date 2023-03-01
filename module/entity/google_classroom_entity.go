package entity

import (
	"google.golang.org/api/classroom/v1"
	"google.golang.org/api/option"
)

type GoogleClassroom interface {
	NewGoogleClassroomService(opts ...option.ClientOption) (*classroom.Service, error)
}
