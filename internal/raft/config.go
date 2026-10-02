package raft

type Config struct {
	ID      string
	Address string
	Peers   map[string]string
}

func LocalConfigs() map[string]Config {
	return map[string]Config{
		"node1": {ID: "node1", Address: "localhost:50051", Peers: map[string]string{"node2": "localhost:50052", "node3": "localhost:50053"}},
		"node2": {ID: "node2", Address: "localhost:50052", Peers: map[string]string{"node1": "localhost:50051", "node3": "localhost:50053"}},
		"node3": {ID: "node3", Address: "localhost:50053", Peers: map[string]string{"node1": "localhost:50051", "node2": "localhost:50052"}},
	}
}

func DockerConfigs() map[string]Config {
	return map[string]Config{
		"node1": {ID: "node1", Address: "node1:50051", Peers: map[string]string{"node2": "node2:50052", "node3": "node3:50053"}},
		"node2": {ID: "node2", Address: "node2:50052", Peers: map[string]string{"node1": "node1:50051", "node3": "node3:50053"}},
		"node3": {ID: "node3", Address: "node3:50053", Peers: map[string]string{"node1": "node1:50051", "node2": "node2:50052"}},
	}
}

func KubernetesConfigs() map[string]Config {
	return map[string]Config{
		"node1": {ID: "node1", Address: "node1:50051", Peers: map[string]string{"node2": "node2:50052", "node3": "node3:50053"}},
		"node2": {ID: "node2", Address: "node2:50052", Peers: map[string]string{"node1": "node1:50051", "node3": "node3:50053"}},
		"node3": {ID: "node3", Address: "node3:50053", Peers: map[string]string{"node1": "node1:50051", "node2": "node2:50052"}},
	}
}

func DefaultConfigs() map[string]Config {
	return LocalConfigs()
}

func ConfigsForEnvironment(env string) map[string]Config {
	switch env {
	case "docker":
		return DockerConfigs()
	case "kubernetes":
		return KubernetesConfigs()
	default:
		return LocalConfigs()
	}
}
