package main

import (
	"github.com/tkanos/go-rest-api-sample/config"
)

var (
	appConfig *config.Config
)

func init() {
	appConfig = config.GetConfig()
}

func main() {

}
