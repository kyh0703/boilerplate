package txpool

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/kyh0703/template/configs"
)

type Tx struct {
	id     string
	time   time.Time
	ctx    context.Context
	cancel context.CancelFunc
	cfg    *configs.Config
}

func NewTx(inner interface{}, cfg *configs.Config) *Tx {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(cfg.App.TxTimeout))
	id := uuid.NewString()
	return &Tx{
		id:     id,
		time:   time.Now(),
		ctx:    ctx,
		cancel: cancel,
		cfg:    cfg,
	}
}

func (tx *Tx) Close() {
	tx.cancel()
}
