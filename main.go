package main

import (
	"podloxx-collector/cmd"
	"podloxx-collector/network"
	"podloxx-collector/utils"

	"github.com/mogenius/mo-go/logger"
)

func main() {
	logger.Init()
	utils.LoadDotEnv()
	network.Init()
	cmd.Execute()
}
