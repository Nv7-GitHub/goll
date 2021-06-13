package goll

func (p *Program) CleanupFunc() {
	if p.fn != nil && p.block.Term == nil {
		p.block.NewRet(nil)
	}
}

func (p *Program) End() {
	p.CleanupFunc()
}
