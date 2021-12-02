package main

import "fmt"

type Priority int

const (
	PriorityDefault    Priority = 0
	PriorityKeepalived          = (iota * 10000) + 1
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
	fmt.Printf("PriorityDefault: %d\n", PriorityDefault)
	fmt.Printf("PriorityKeepalived: %d\n", PriorityKeepalived)
	fmt.Printf("PriorityMetricsServer: %d\n", PriorityMetricsServer)
	fmt.Printf("PriorityCoredns: %d\n", PriorityCoredns)
	fmt.Printf("PriorityContainerd: %d\n", PriorityContainerd)
}

// PriorityDefault: 0
// PriorityKeepalived: 10001
// PriorityMetricsServer: 20001
// PriorityCoredns: 30001
// PriorityContainerd: 130001
