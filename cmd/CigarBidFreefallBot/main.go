package main

import (
	"fmt"

	"github.com/KaylaHood/CigarBidFreefallBot/cigarbot"
)

func main() {
	var err error
	defer func() {
		if r := recover(); r != nil {
			err = r.(error)
			panic(fmt.Errorf("CigarBidFreefallBot.main(): Execution finished with an error, err = \n\t%v", err))
		}
	}()
	fmt.Println("Hello, world!")
	cigarbot.StartSelenium()
}
