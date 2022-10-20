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

func LoadDotEnv() {
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

	envPath := fmt.Sprintf(".env/%s.env", stage)
	tes, readErr := godotenv.Read(envPath)
	if readErr != nil {
		logger.Log.Fatal("Error loading .env file")
	}
	logger.Init()
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

	err := godotenv.Load(envPath)

	if err != nil {
		logger.Log.Fatal("Error loading .env file")
	}
}
