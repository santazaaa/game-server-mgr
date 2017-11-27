package utils

type PortMgr struct {
	startPort int
	portCount int
	availablePorts []int
}

// Init init port manager
func (p *PortMgr) Init(start int, count int) {
	p.startPort = start
	p.portCount = count
	for i := 0; i < count; i++ {
		nextPort := start + i
		p.availablePorts = append(p.availablePorts, nextPort)
	}
}

// GetNext get next available port from list
func (p *PortMgr) GetNext() int {
	if len(p.availablePorts) == 0 {
		return -1
	}
	port := p.availablePorts[0]
	println(port)
	p.availablePorts = p.availablePorts[1:]
	return port
}

// Free return an unused port to the available port list
func (p *PortMgr) Free(port int) {
	p.availablePorts = append(p.availablePorts, port)
}