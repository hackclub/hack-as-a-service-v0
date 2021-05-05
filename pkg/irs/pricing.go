package irs

import "github.com/shopspring/decimal"

func (output *ProcessedOutput) price() decimal.Decimal {
	// Pricing (https://hackclub.slack.com/archives/C01N3B30TFB/p1619626597140600)
	cpu_cost := decimal.RequireFromString("0.000001929012345679").Mul(output.CpuUsage)
	cpu_cost = decimal.Max(decimal.Zero, cpu_cost)

	// TODO: memory pricing

	return cpu_cost
}
