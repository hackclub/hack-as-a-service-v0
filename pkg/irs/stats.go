package irs

import (
	"time"

	"github.com/shopspring/decimal"
)

type Stats struct {
	Read     time.Time `json:"read"`
	CpuStats struct {
		CpuUsage struct {
			PercpuUsage []interface{} `json:"percpu_usage"`
			TotalUsage  int64         `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCpuUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int64 `json:"online_cpus"`
	} `json:"cpu_stats"`
	PrecpuStats struct {
		CpuUsage struct {
			TotalUsage int64 `json:"total_usage"`
		} `json:"cpu_usage"`
		SystemCpuUsage int64 `json:"system_cpu_usage"`
		OnlineCpus     int64 `json:"online_cpus"`
	} `json:"precpu_stats"`
	MemoryStats struct {
		Stats struct {
			Cache int64 `json:"cache"`
		} `json:"stats"`
		Usage int64 `json:"usage"`
		Limit int64 `json:"limit"`
	} `json:"memory_stats"`
}

type ProcessedOutput struct {
	Timestamp       time.Time
	MemUsage        decimal.Decimal
	MemUsagePercent float64
	CpuUsage        decimal.Decimal
	CpuUsagePercent float64
}

func (stat *Stats) Process() ProcessedOutput {
	// https://docs.docker.com/engine/api/v1.41/#operation/ContainerStats
	mem_usage := decimal.NewFromInt(stat.MemoryStats.Usage).Sub(decimal.NewFromInt(stat.MemoryStats.Stats.Cache)).Div(decimal.NewFromInt(stat.MemoryStats.Limit))
	cpu_usage := decimal.NewFromInt(stat.CpuStats.CpuUsage.TotalUsage).Sub(decimal.NewFromInt(stat.PrecpuStats.CpuUsage.TotalUsage))
	x := decimal.NewFromInt(stat.CpuStats.SystemCpuUsage).Sub(decimal.NewFromInt(stat.PrecpuStats.SystemCpuUsage))
	if !x.IsZero() {
		cpu_usage = cpu_usage.Div(x)
	}
	cpu_usage = cpu_usage.Mul(decimal.NewFromInt(int64(len(stat.CpuStats.CpuUsage.PercpuUsage))))

	mem_usage_percent, _ := mem_usage.Mul(decimal.NewFromInt(100)).Float64()
	cpu_usage_percent, _ := cpu_usage.Mul(decimal.NewFromInt(100)).Float64()

	output := ProcessedOutput{
		Timestamp:       stat.Read,
		MemUsage:        mem_usage,
		MemUsagePercent: mem_usage_percent,
		CpuUsage:        cpu_usage,
		CpuUsagePercent: cpu_usage_percent,
	}

	return output
}
