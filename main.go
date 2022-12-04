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

//go:embed .env/production.env
var defaultEnvFile string

func main() {
	utils.DefaultEnvFile = defaultEnvFile
	api.HtmlDirFs = htmlDirFs
	logger.Init()
	utils.LoadDotEnv()
	network.Init()
	cmd.Execute()
}
