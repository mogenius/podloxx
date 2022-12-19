package utils

import (
	"fmt"
	"os"
	"os/exec"
	"podloxx/logger"
	"runtime"
	"strings"
	"unicode/utf8"

	"github.com/joho/godotenv"
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
	if len(os.Args) > 1 {
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
		if Contains(passwordStrList, strings.ToUpper(key)) {
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

func Contains(s []string, str string) bool {
	for _, v := range s {
		if strings.Contains(str, v) {
			return true
		}
	}
	return false
}

func OpenBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}

	if err != nil {
		fmt.Errorf("error while opening browser, %v", err)
	}
}
