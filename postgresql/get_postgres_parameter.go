package main

import (
	"fmt"
	"github.com/go-park-mail-ru/2019_1_SleeplessNights/shared/config"
	"os"
)

func main() {
	fmt.Println(config.Get("postgres." + os.Args[1]))
}