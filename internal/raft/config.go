package raft

type Config struct {
	ID      string
	Address string
	Peers   map[string]string
}

func DefaultConfigs() map[string]Config {
	return map[string]Config{
		"node1": {ID: "node1", Address: "localhost:50051", Peers: map[string]string{"node2": "localhost:50052", "node3": "localhost:50053"}},
		"node2": {ID: "node2", Address: "localhost:50052", Peers: map[string]string{"node1": "localhost:50051", "node3": "localhost:50053"}},
		"node3": {ID: "node3", Address: "localhost:50053", Peers: map[string]string{"node1": "localhost:50051", "node2": "localhost:50052"}},
	}
}
