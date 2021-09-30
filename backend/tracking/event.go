package tracking

import (
	"context"
	"time"

	"github.com/avct/uasurfer"

	"github.com/k-yomo/kagu-miru/backend/request"

	"github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
)

type Event struct {
	// passed from client
	ID        string                 `json:"id"`
	Action    string                 `json:"action"`
	CreatedAt time.Time              `json:"created_at"`
	Params    map[string]interface{} `json:"params"`

	// fill on backend
	UserID    string `json:"user_id"`
	UserAgent string `json:"user_agent"`
	Device    string `json:"devise"`
	IPAddress string `json:"ip"`
}

func NewEvent(ctx context.Context, gqlEvent gqlmodel.Event) *Event {
	event := newDefaultEvent(ctx)
	event.ID = gqlEvent.ID.String()
	event.Action = gqlEvent.Action.String()
	event.CreatedAt = gqlEvent.CreatedAt
	event.Params = gqlEvent.Params
	return event
}

func newDefaultEvent(ctx context.Context) *Event {
	req, ok := request.GetRequestFromCtx(ctx)
	if !ok {
		return &Event{}
	}

	event := &Event{
		UserAgent: req.UserAgent(),
		Device:    uasurfer.Parse(req.UserAgent()).DeviceType.StringTrimPrefix(),
		IPAddress: request.RealClientIP(req),
	}
	return event
}
