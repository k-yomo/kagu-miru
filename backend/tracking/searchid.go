package tracking

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	"github.com/oklog/ulid/v2"

	"github.com/alexedwards/scs/v2"
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
	searchID = uuid()
	s.sessionManager.Put(ctx, "searchId", searchID)
	return searchID
}

func uuid() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
