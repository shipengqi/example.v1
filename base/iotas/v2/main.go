package main

import "fmt"

type Priority int

const (
	PriorityKeepalived Priority = iota + 10
	PriorityMetricsServer
	PriorityCoredns
	PriorityCalico
	PriorityFlannel
	PriorityKubeProxy
	PriorityKubeScheduler
	PriorityKubeControllerManager
	PriorityKubeApiserver
	PriorityEtcd
	PriorityKubelet
	PriorityCriO
	PriorityContainerd
)

func main() {
	fmt.Printf("PriorityKeepalived: %d\n", PriorityKeepalived)
	fmt.Printf("PriorityMetricsServer: %d\n", PriorityMetricsServer)
	fmt.Printf("PriorityCoredns: %d\n", PriorityCoredns)
	fmt.Printf("PriorityContainerd: %d\n", PriorityContainerd)
}

// PriorityKeepalived: 1
// PriorityMetricsServer: 2
// PriorityCoredns: 3
// PriorityContainerd: 13
