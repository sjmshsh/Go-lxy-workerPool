package internal

import "github.com/sjmshsh/model"

type Option func(pool *model.Pool)

func WithBlock(block bool) Option {
	return func(p *model.Pool) {
		p.Block = block
	}
}

func WithPreAllocWorkers(preAlloc bool) Option {
	return func(p *model.Pool) {
		p.PreAlloc = preAlloc
	}
}
