package raft

import "testing"

func TestDefaultConfigs(t *testing.T) {
	configs := DefaultConfigs()
	if len(configs) != 3 {
		t.Fatalf("expected 3 nodes, got %d", len(configs))
	}
	if configs["node1"].Address != "localhost:50051" {
		t.Fatalf("unexpected node1 address: %s", configs["node1"].Address)
	}
	if len(configs["node1"].Peers) != 2 {
		t.Fatalf("expected 2 peers, got %d", len(configs["node1"].Peers))
	}
}

func TestKubernetesConfigsUseServiceDNS(t *testing.T) {
	configs := KubernetesConfigs()
	if configs["node1"].Address != "node1:50051" {
		t.Fatalf("unexpected node1 Kubernetes address: %s", configs["node1"].Address)
	}
	if configs["node1"].Peers["node2"] != "node2:50052" || configs["node1"].Peers["node3"] != "node3:50053" {
		t.Fatalf("unexpected Kubernetes peers: %#v", configs["node1"].Peers)
	}
	if got := ConfigsForEnvironment("kubernetes")["node2"].Address; got != "node2:50052" {
		t.Fatalf("environment selection returned %q", got)
	}
}
