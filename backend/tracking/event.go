package tracking

import (
	"context"

	"github.com/avct/uasurfer"

	"github.com/k-yomo/kagu-miru/backend/request"

	"github.com/k-yomo/kagu-miru/backend/graph/gqlmodel"
)

type Event struct {
	gqlmodel.Event
	UserID    string `json:"userId"`
	UserAgent string `json:"userAgent"`
	Device    string `json:"devise"`
	IPAddress string `json:"ip"`
}

func NewEvent(ctx context.Context, gqlEvent gqlmodel.Event) *Event {
	event := newDefaultEvent(ctx)
	event.Event = gqlEvent
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
