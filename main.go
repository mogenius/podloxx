package main

import (
	"embed"
	"podloxx-collector/api"
	"podloxx-collector/cmd"
	"podloxx-collector/network"
	"podloxx-collector/utils"

	"github.com/mogenius/mo-go/logger"
)

//go:embed ui/dist/podloxx/*
var htmlDirFs embed.FS

func main() {
	logger.Init()
	utils.LoadDotEnv()
	api.HtmlDirFs = htmlDirFs
	network.Init()
	cmd.Execute()
}
