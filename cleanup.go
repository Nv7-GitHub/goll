package goll

func (p *Program) CleanupFunc() {
	if p.Fn != nil && p.Block.Term == nil {
		p.Block.NewRet(nil)
	}
}

func (p *Program) End() {
	p.CleanupFunc()
}
