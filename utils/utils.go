package utils

import (
	"fmt"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
	"github.com/mogenius/mo-go/logger"
	"github.com/mogenius/mo-go/utils"
)

var DefaultEnvFile string

func LoadDotEnv() {
	logger.Init()

	homeDirName, homeDirErr := os.UserHomeDir()
	if homeDirErr != nil {
		logger.Log.Error(homeDirErr)
	}

	passwordStrList := []string{"PWD", "PASSWORD"}
	var stage string = "local"
	if len(os.Args) > 0 {
		stage = os.Args[1]
	}
	if os.Getenv("STAGE") != "" {
		stage = os.Getenv("STAGE")
	} else if stage != "" {
		stage = "local"
	}

	if _, err := os.Stat(homeDirName + "/.podloxx/podloxx.env"); err == nil || os.IsExist(err) {
		godotenv.Load(homeDirName + "/.podloxx/podloxx.env")
	} else {
		// load file from embedded
		tes, unmarshallErr := godotenv.Unmarshal(DefaultEnvFile)
		logger.Log.Info(tes)
		if unmarshallErr != nil {
			logger.Log.Fatal("Error unmarshalling .env file")
			logger.Log.Fatal(unmarshallErr)
		}
		// write it to default location
		folderErr := os.Mkdir(homeDirName+"/.podloxx/", 0755)
		if folderErr != nil {
			logger.Log.Fatal("Error creating folder " + homeDirName + "/.podloxx/")
			logger.Log.Fatal(folderErr)
		}

		err := godotenv.Write(tes, homeDirName+"/.podloxx/podloxx.env")
		if err != nil {
			logger.Log.Fatal("Error writing " + homeDirName + "/.podloxx/podloxx.env file")
			logger.Log.Fatal(err)
		}
		godotenv.Load(homeDirName + "/.podloxx/podloxx.env")
	}

	tes, readErr := godotenv.Read(homeDirName + "/.podloxx/podloxx.env")
	if readErr != nil {
		logger.Log.Fatal("Error reading " + homeDirName + "/.podloxx/podloxx.env file")
		logger.Log.Fatal(readErr)
	}

	for key, element := range tes {
		if utils.Contains(passwordStrList, strings.ToUpper(key)) {
			if len(element) == 0 {
				element = os.Getenv(key)
			}
			maxStr := 5
			strLen := len(element)
			if strLen < 5 {
				maxStr = 2
			}
			logger.Log.Notice("Key:", key, "=>", "Element:", fmt.Sprintf("%s%s", element[0:maxStr], strings.Repeat("*", utf8.RuneCountInString(element[maxStr:strLen-1]))))
		} else {
			logger.Log.Notice("Key:", key, "=>", "Element:", element)
		}
	}
}
