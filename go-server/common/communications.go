package common

import (
	"os"
	"sync"
)

var FastAPIURL = ""

func GetFastAPIURL() string {
	return FastAPIURL
}

func SetFastAPIURL() {
	FastAPIURL = os.Getenv("FASTAPI_URL")
}

var (
	GlobalMutex sync.Mutex
)
