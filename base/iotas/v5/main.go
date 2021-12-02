package main

import "fmt"

type Priority int

const (
	PriorityDefault    Priority = iota * 10000
	PriorityKeepalived
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
// PriorityKeepalived: 10000
// PriorityMetricsServer: 20000
// PriorityCoredns: 30000
// PriorityContainerd: 130000
