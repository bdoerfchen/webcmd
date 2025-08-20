package shellexecuter

import (
	"context"

	"github.com/bdoerfchen/webcmd/src/model/process"
)

type shellPool struct {
	pool chan *process.Process
}

func NewPool(size int, template process.Template) *shellPool {
	// Create pool
	pool := &shellPool{
		pool: make(chan *process.Process, size),
	}

	// Start worker that constantly fills the pool
	go filler(pool, &template)

	return pool
}

// Worker function that inserts new processes into the pool, stopping only when the capacity is reached
func filler(pool *shellPool, template *process.Template) {

	for {
		// Create process with connected input and output buffers
		proc, err := process.Prepare(template)
		if err != nil {
			panic("failed to start process: " + err.Error())
		}

		// Start process
		proc.Proc.Start()

		// Insert into queue
		pool.pool <- proc
	}
}

// Number of processes that can be prepared
func (p *shellPool) Capacity() int {
	return cap(p.pool)
}

// Number of processes available to take
func (p *shellPool) Available() int {
	return len(p.pool)
}

// Take a process from the pool. Returns error if the context is closed before.
func (p *shellPool) Take(ctx context.Context) (*process.Process, error) {
	select {
	case proc := <-p.pool:
		return proc, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
