package cigarbot

import (
	"github.com/KaylaHood/CigarBidFreefallBot/seleniumwindowscompatibility"
)

const (
	// These paths will be different on your system.
	seleniumPath     = "C:\\WebDriver\\bin\\selenium-server-standalone-3.14.0.jar"
	chromeDriverPath = "C:\\WebDriver\\bin\\chromedriver.exe"
	port             = 9515
	debug            = true
)

var s seleniumwindowscompatibility.Service
var wd seleniumwindowscompatibility.WebDriver

// StartSelenium Start a Selenium WebDriver server instance
func StartSelenium() {
	// Start a Selenium WebDriver server instance (if one is not already
	// running)

}
