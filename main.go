package main

import (
	"embed"
	"podloxx/api"
	"podloxx/cmd"
	"podloxx/logger"
	"podloxx/network"
	"podloxx/utils"
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
