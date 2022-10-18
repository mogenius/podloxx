package main

import (
	"podloxx-collector/cmd"
	"podloxx-collector/network"

	"github.com/mogenius/mo-go/logger"
)

func main() {
	logger.Init()
	network.Init()
	cmd.Execute()
}
