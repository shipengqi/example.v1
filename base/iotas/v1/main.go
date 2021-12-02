package main

import "fmt"

type Priority int

const (
	PriorityKeepalived    Priority = 1 << iota // 1 << 0 which is 00000001
	PriorityMetricsServer                      // 1 << 1 which is 00000010
	PriorityCoredns                            // 1 << 2 which is 00000100
	PriorityCalico                             // 1 << 3 which is 00001000
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
// PriorityCoredns: 4
// PriorityContainerd: 4096
