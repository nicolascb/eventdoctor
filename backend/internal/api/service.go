package api

import (
	"context"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

type Service interface {
	SaveSpec(ctx context.Context, spec eventdoctor.EventDoctorSpec) error
	ListProducers(ctx context.Context, page, pageSize int, search string) (*response.ProducersListView, error)
	GetProducerDetail(ctx context.Context, serviceID, topicID int64) (*response.ProducerDetailView, error)
	ListConsumers(ctx context.Context, page, pageSize int, search string) (*response.ConsumerView, error)
	ListUndocumentedConsumers(ctx context.Context) (*response.UndocumentedConsumerView, error)
	ListEvents(ctx context.Context, page, pageSize int, search string) (*response.EventsListView, error)
	Overview(ctx context.Context) (*response.OverviewResponse, error)
	GetTopicView(ctx context.Context, topicName string) (*response.TopicView, error)
	GetServiceView(ctx context.Context, serviceName string) (*response.ServiceView, error)

	ListTopics(ctx context.Context, page, pageSize int) (*response.TopicListView, error)
	GetEvent(ctx context.Context, id int64) (*response.EventView, error)
}
