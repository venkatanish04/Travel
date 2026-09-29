package raft

type Peer struct {
	ID      string
	Address string
}

type Transport struct{ Peers []Peer }

func NewTransport(peers []Peer) *Transport { return &Transport{Peers: peers} }
