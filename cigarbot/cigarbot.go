package cigarbot

import (
	"fmt"

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
		Password: "Cigars4MeAndNug"}
	var opts = common.SeleniumOptions{
		SeleniumPath:     seleniumPath,
		ChromeDriverPath: chromeDriverPath,
		GeckoDriverPath:  "",
		BrowserName:      "chrome",
		Port:             port,
		Debug:            debug}
	var cbs CigarBidService
	cbs, err := NewCigarBidService(creds, opts)
	if err != nil {
		fmt.Printf("There was an error when creating a New Cigar Bid Service!! This was the error:\n\t%v", err)
	} else {
		fmt.Println("There were NO errors encountered. Congratulations!")
	}
	defer cbs.Shutdown()

	err = cbs.Login()
}
