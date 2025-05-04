package repository

import (
	"context"
	"database/sql"
	"fmt"
	"sync"
	"time"

	"github.com/kyh0703/template/internal/core/domain/model"
	"github.com/kyh0703/template/internal/core/domain/repository"
)

type memoryOAuthRepository struct {
	mu     sync.RWMutex
	states map[string]model.OauthState
}

func NewMemoryOAuthRepository() repository.OAuthRepository {
	return &memoryOAuthRepository{
		states: make(map[string]model.OauthState),
	}
}

func (r *memoryOAuthRepository) CreateState(ctx context.Context, arg model.CreateOAuthStateParams) (model.OauthState, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	state := model.OauthState{
		ID:          int64(len(r.states) + 1), // 단순한 ID 생성
		State:       arg.State,
		RedirectUrl: arg.RedirectUrl,
		ExpiresAt:   arg.ExpiresAt,
		CreateAt:    sql.NullString{String: time.Now().Format(time.RFC3339), Valid: true},
	}

	r.states[arg.State] = state
	return state, nil
}

func (r *memoryOAuthRepository) FindByState(ctx context.Context, state string) (model.OauthState, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if savedState, ok := r.states[state]; ok {
		return savedState, nil
	}
	return model.OauthState{}, fmt.Errorf("state not found")
}

func (r *memoryOAuthRepository) DeleteState(ctx context.Context, state string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.states[state]; !ok {
		return fmt.Errorf("state not found")
	}

	delete(r.states, state)
	return nil
}
