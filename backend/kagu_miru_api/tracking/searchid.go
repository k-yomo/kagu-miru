package tracking

import (
	"context"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/k-yomo/kagu-miru/backend/pkg/uuid"
)

type SearchIDManager struct {
	sessionManager *scs.SessionManager
}

func NewSearchIDManager(isDeployed bool) *SearchIDManager {
	sessionManager := scs.New()
	sessionManager.IdleTimeout = 30 * time.Minute
	sessionManager.Cookie.Persist = false
	if isDeployed {
		sessionManager.Cookie.Secure = true
	}
	return &SearchIDManager{
		sessionManager: sessionManager,
	}
}

func (s *SearchIDManager) Middleware() func(next http.Handler) http.Handler {
	return s.sessionManager.LoadAndSave
}

func (s *SearchIDManager) GetSearchID(ctx context.Context) string {
	searchID := s.sessionManager.GetString(ctx, "searchId")
	if searchID != "" {
		return searchID
	}
	searchID = uuid.UUID()
	s.sessionManager.Put(ctx, "searchId", searchID)
	return searchID
}
