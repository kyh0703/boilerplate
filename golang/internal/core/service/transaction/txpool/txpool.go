package txpool

import (
	"fmt"
	"sync"

	"github.com/gofrs/uuid"
)

type TxPool struct {
	txMap      map[uuid.UUID]*Tx
	txCountMap map[string]int
	wg         sync.WaitGroup
	mutex      sync.RWMutex
}

func NewTxPool() *TxPool {
	return &TxPool{
		txMap:      make(map[uuid.UUID]*Tx),
		txCountMap: make(map[string]int),
	}
}

func (p *TxPool) PrintTx() {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	for _, tx := range p.txMap {
		fmt.Println(tx)
	}
}

func (p *TxPool) BeginTx(inner *interface{}) {
	// create tx data.
	p.mutex.Lock()
	defer p.mutex.Unlock()
	// compute tx type and count.
	uuid, _ := uuid.NewV4()
	p.wg.Add(1)
	go func() {
		defer p.wg.Done()
		defer p.EndTx(uuid)
	}()
}

func (p *TxPool) EndTx(id uuid.UUID) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	tx, ok := p.txMap[id]
	if !ok {
		return
	}
	tx.Close()
	delete(p.txMap, id)
}

func (p *TxPool) EndAllTx() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, tx := range p.txMap {
		tx.Close()
	}
	p.wg.Wait()
	for k := range p.txMap {
		delete(p.txMap, k)
	}
	for t := range p.txCountMap {
		delete(p.txCountMap, t)
	}
}
