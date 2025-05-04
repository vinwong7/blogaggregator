package main

import (
	"fmt"

	"github.com/vinwong7/blogaggregator/internal/config"
)

func main() {
	configFile, _ := config.Read()
	fmt.Println(configFile)
	configFile.SetUser("vinwong7")
	configFile, _ = config.Read()
	fmt.Println(configFile)

}
