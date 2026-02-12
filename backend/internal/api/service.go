package api

import (
	"context"

	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

type Service interface {
	SaveSpec(ctx context.Context, spec eventdoctor.EventDoctorSpec) error
	ListProducers(ctx context.Context) ([]eventdoctor.Producer, error)
	ListConsumers(ctx context.Context) ([]eventdoctor.Consumer, error)
	ListEvents(ctx context.Context) ([]eventdoctor.Event, error)
}
