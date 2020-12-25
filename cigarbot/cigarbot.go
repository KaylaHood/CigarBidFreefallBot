package cigarbot

import (
	common "github.com/KaylaHood/CigarBidFreefallBot"
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
	var creds = common.LoginCredentials{
		Username: "kaylahood1996@gmail.com",
		Password: "TheseR!TheDroidsURLooking4"}
	var opts = common.SeleniumOptions{
		SeleniumPath:     seleniumPath,
		ChromeDriverPath: chromeDriverPath,
		GeckoDriverPath:  "",
		BrowserName:      "chrome",
		Port:             port,
		Debug:            debug}
	var cbs CigarBidService
	cbs = NewCigarBidService(creds, opts)
	defer cbs.Shutdown()
}
