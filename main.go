package main

import (
	"github.com/ripross/monitoring_demo/seed"
)

func main() {
	seed.GenerateDatabaseData()
	seed.GenerateTimeSeriesData()
}