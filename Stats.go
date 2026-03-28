package main

type Stats struct {
	CPUStats struct {
		CPUUsage struct {
			TotalUsage uint64 `json:"total_usage"`
		} `json:"cpu_usage"`
	} `json:"cpu_stats"`

	MemoryStats struct {
		Usage uint64 `json:"usage"`
	} `json:"memory_stats"`
}