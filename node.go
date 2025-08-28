package raft

type Node interface {
	Tick()
}

type Peer struct {
	ID      uint64
	Context []byte
}

func setupNode(c *Config, peers []Peer) *node {
	if len(peers) == 0 {
		panic("no peers given; use RestartNode instead")
	}
	rn, err := NewRawNode(c)
	if err != nil {
		panic(err)
	}
	err = rn.Bootstrap(peers)
	if err != nil {
		c.Logger.Warningf("error occurred during starting a new node: %v", err)
	}

	n := newNode(rn)
	return &n
}

func StartNode(c *Config, peers []Peer) Node {
	n := setupNode(c, peers)
	go n.run()
	return n
}

type node struct {
	tickc chan struct{}
	done  chan struct{}
	rn    *RawNode
}

func newNode(rn *RawNode) node {
	return node{
		tickc: make(chan struct{}, 128),
		rn:    rn,
	}
}

func (n *node) run() {
	for {
		select {
		case <-n.tickc:
			n.rn.Tick()
		}
	}
}

// Tick increments the internal logical clock for this Node. Election timeouts
// and heartbeat timeouts are in units of ticks.
func (n *node) Tick() {
	select {
	case n.tickc <- struct{}{}:
	case <-n.done:
	default:
		n.rn.raft.logger.Warningf("%x A tick missed to fire. Node blocks too long!", n.rn.raft.id)
	}
}
