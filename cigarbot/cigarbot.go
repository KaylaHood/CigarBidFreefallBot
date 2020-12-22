package cigarbot

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/KaylaHood/CigarBidFreefallBot/seleniumwindowscompatibility"
)

const (
	// These paths will be different on your system.
	seleniumPath     = "C:\\WebDriver\\bin\\selenium-server-standalone-3.14.0.jar"
	chromeDriverPath = "C:\\WebDriver\\bin\\chromedriver.exe"
	port             = 9515
)

var s seleniumwindowscompatibility.Service
var wd seleniumwindowscompatibility.WebDriver

// Start a Selenium WebDriver server instance
func StartSelenium() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running)

	opts := []seleniumwindowscompatibility.ServiceOption{
		seleniumwindowscompatibility.ChromeDriver(chromeDriverPath), // Specify the path to ChromeDriver in order to use Chrome
		seleniumwindowscompatibility.Output(os.Stderr),              // Output debug information to STDERR
	}
	seleniumwindowscompatibility.SetDebug(true)
	myService, err := seleniumwindowscompatibility.NewSeleniumService(seleniumPath, port, opts...)
	if err != nil {
		panic(err) // panic is used only as an example and is not otherwise recommended
	}
	defer myService.Stop()

	// Connect to the WebDriver instance running locally
	caps := seleniumwindowscompatibility.Capabilities{"browserName": "chrome"}
	wd, err := seleniumwindowscompatibility.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(err)
	}
	defer wd.Quit()

	// Navigate to the simple playground interface
	if err := wd.Get("http://play.golang.org/?simple=1"); err != nil {
		panic(err)
	}

	// Get a reference to the text box containing code
	elem, err := wd.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#code")
	if err != nil {
		panic(err)
	}
	// Remove the boilerplate code already in the text box
	if err := elem.Clear(); err != nil {
		panic(err)
	}

	// Enter some new code in text box
	err = elem.SendKeys(`
		package main
		import "fmt"

		func main() {
			fmt.Println("Hello WebDriver!\n")
		}
	`)
	if err != nil {
		panic(err)
	}

	// Click the run button
	btn, err := wd.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#run")
	if err != nil {
		panic(err)
	}
	if err := btn.Click(); err != nil {
		panic(err)
	}

	// Wait for the program to finish running and get the output
	outputDiv, err := wd.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#output")
	if err != nil {
		panic(err)
	}

	var output string
	for {
		output, err = outputDiv.Text()
		if err != nil {
			panic(err)
		}
		if output != "Waiting for remote server..." {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	fmt.Printf("%s", strings.Replace(output, "\n\n", "\n", -1))

	// Example Output:
	// Hello WebDriver!
	//
	// Program exited.
}
