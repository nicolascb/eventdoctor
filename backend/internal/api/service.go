package api

import (
	"context"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

type Service interface {
	SaveSpec(ctx context.Context, spec eventdoctor.EventDoctorSpec) error
	ListProducers(ctx context.Context) ([]response.ProducerView, error)
	ListConsumers(ctx context.Context) (*response.ConsumerView, error)
	ListEvents(ctx context.Context) ([]response.TopicEventsView, error)
	Overview(ctx context.Context) (*response.OverviewResponse, error)
	GetTopicView(ctx context.Context, topicName string) (*response.TopicView, error)
	GetServiceView(ctx context.Context, serviceName string) (*response.ServiceView, error)
}
