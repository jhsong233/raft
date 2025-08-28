package raft

import "errors"

func (rn *RawNode) Bootstrap(peers []Peer) error {
	if len(peers) == 0 {
		return errors.New("must provide at least one peer to Bootstrap")
	}

	rn.raft.becomeFollower(1, None)
	return  nil
}