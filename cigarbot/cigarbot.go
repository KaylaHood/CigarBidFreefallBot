package cigarbot

import (
	common "github.com/KaylaHood/CigarBidFreefallBot"
)

const (
	// These paths will be different on your system.
	seleniumPath     = "C:\\WebDriver\\bin\\selenium-server-standalone-3.14.0.jar"
	chromeDriverPath = "C:\\WebDriver\\bin\\chromedriver.exe"
	port             = 9515
	debug            = true
)

var cbs CigarBidService

// RunCigarBot creates a CigarBidService and initiates bot
func RunCigarBot() {
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
	cbs = NewCigarBidService(creds, opts)
	defer cbs.Shutdown()

	cbs.FindFreefallProduct("https://www.cigarbid.com/a/cao-brazilia-lambada/3783948/")
}
