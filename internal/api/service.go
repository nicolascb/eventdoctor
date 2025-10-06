package api

import (
	"context"

	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

type Service interface {
	SaveSpec(ctx context.Context, spec eventdoctor.EventDoctorSpec) error
	WebDocs(ctx context.Context) (eventdoctor.WebData, error)
}
