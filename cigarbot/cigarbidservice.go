package cigarbot

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	common "github.com/KaylaHood/CigarBidFreefallBot"
	"github.com/KaylaHood/CigarBidFreefallBot/seleniumwindowscompatibility"
)

const (
	cigarBidURL        string = "https://www.cigarbid.com/"
	defaultWaitTimeSec int64  = 1
)

// CigarBidService provides utilities for managing session with CigarBid site
type CigarBidService interface {
	FindFreefallProduct(string)
	navigateToCigarBid()
	focusOnPage()
	Login()
	Shutdown()
}

type cbService struct {
	creds           common.LoginCredentials
	seleniumService *seleniumwindowscompatibility.Service
	webDriver       seleniumwindowscompatibility.WebDriver
	debugMode       bool
	cookiesAccepted bool
}

// NewCigarBidService creates a Webdriver in CHROME and logs in with the given username and password
// then returns the resulting CigarBidService, or the error if one was thrown
func NewCigarBidService(newCreds common.LoginCredentials, newOpts common.SeleniumOptions) CigarBidService {
	newCbService := new(cbService)
	newCbService.creds = newCreds
	newCbService.debugMode = newOpts.Debug
	newCbService.cookiesAccepted = false

	var err error = nil

	sOpts := []seleniumwindowscompatibility.ServiceOption{
		seleniumwindowscompatibility.ChromeDriver(newOpts.ChromeDriverPath),
		seleniumwindowscompatibility.GeckoDriver(newOpts.GeckoDriverPath),
		seleniumwindowscompatibility.Output(os.Stderr), // Output debug information to STDERR
	}
	seleniumwindowscompatibility.SetDebug(newCbService.debugMode)
	newCbService.seleniumService, err = seleniumwindowscompatibility.NewSeleniumService(newOpts.SeleniumPath, newOpts.Port, sOpts...)
	if err != nil {
		panic(fmt.Errorf("CigarBidService.NewCigarBidService(...): Failed to create a new Selenium Service\n\t%v", err))
	}
	// Connect to the WebDriver instance running locally
	caps := seleniumwindowscompatibility.Capabilities{"browserName": newOpts.BrowserName}
	newCbService.webDriver, err = seleniumwindowscompatibility.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
	if err != nil {
		panic(fmt.Errorf("CigarBidService.NewCigarBidService(...): Failed to create a new Remote Web Driver\n\t%v", err))
	}
	newCbService.navigateToCigarBid()
	newCbService.Login()
	return newCbService
}

func (cbs *cbService) FindFreefallProduct(url string) {
	script, err := ioutil.ReadFile("javascript/goto_ff_product.js")
	if err != nil {
		panic(fmt.Errorf("CigarBidService.FindFreefallProduct(): Failed to read goto_ff_product.js\n\t%v", err))
	}
	args := make([]interface{}, 1)
	args[0] = url
	result, err := cbs.webDriver.ExecuteScript(string(script), args)
	if cbs.debugMode {
		fmt.Printf("CigarBidService.FindFreefallProduct(): Go to FF Product Script result: %v\n", result)
	}
	if err != nil {
		panic(fmt.Errorf("CigarBidService.FindFreefallProduct(): Failed to run goto_ff_product.js\n\t%v", err))
	}
	cbs.dismissPopups()
	cbs.sleep(defaultWaitTimeSec)
	cbs.focusOnPage()
	script, err = ioutil.ReadFile("javascript/watch_ff_product.js")
	if err != nil {
		panic(fmt.Errorf("CigarBidService.FindFreefallProduct(): Failed to read watch_ff_product.js\n\t%v", err))
	}
	result, err = cbs.webDriver.ExecuteScript(string(script), make([]interface{}, 0))
	if cbs.debugMode {
		fmt.Printf("CigarBidService.FindFreefallProduct(): Watch FF Product Script result: %v\n", result)
	}
	if err != nil {
		panic(fmt.Errorf("CigarBidService.FindFreefallProduct(): Failed to run watch_ff_product.js\n\t%v", err))
	}
}

func (cbs *cbService) navigateToCigarBid() {
	err := cbs.webDriver.Get(cigarBidURL)
	cbs.dismissPopups()
	cbs.sleep(defaultWaitTimeSec)
	cbs.focusOnPage()
	if err != nil {
		panic(fmt.Errorf("CigarBidService.navigateToCigarBid(): Failed to navigate to Cigar Bid\n\t%v", err))
	}
}

func (cbs *cbService) sleep(seconds int64) {
	if cbs.debugMode {
		fmt.Printf("CigarBidService.sleep(): Waiting for %d seconds..., current time: %s\n", seconds, time.Now().Format("15:04:05.00000"))
	}
	time.Sleep(time.Duration(seconds) * time.Second)
	if cbs.debugMode {
		fmt.Printf("CigarBidService.sleep(): Finished waiting for %d seconds..., current time: %s\n", seconds, time.Now().Format("15:04:05.00000"))
	}
}

func (cbs *cbService) focusOnPage() {
	script, err := ioutil.ReadFile("javascript/focus_on_page.js")
	if err != nil {
		panic(fmt.Errorf("CigarBidService.focusOnPage(): Failed to read focus_on_page.js\n\t%v", err))
	}
	result, err := cbs.webDriver.ExecuteScript(string(script), make([]interface{}, 0))
	if cbs.debugMode {
		fmt.Printf("CigarBidService.focusOnPage(): Script result: %v\n", result)
	}
	if err != nil {
		panic(fmt.Errorf("CigarBidService.focusOnPage(): Failed to run focus_on_page.js\n\t%v", err))
	}
}

// dismissPopups will run dismiss_popups.js on the current page
func (cbs *cbService) dismissPopups() {
	script, err := ioutil.ReadFile("javascript/dismiss_popups.js")
	if err != nil {
		panic(fmt.Errorf("CigarBidService.DismissPopups(): Failed to read dismiss_popups.js\n\t%v", err))
	}
	result, err := cbs.webDriver.ExecuteScript(string(script), make([]interface{}, 0))
	if cbs.debugMode {
		fmt.Printf("CigarBidService.DismissPopups(): Script result: %v\n", result)
	}
	if err != nil {
		panic(fmt.Errorf("CigarBidService.DismissPopups(): Failed to run dismiss_popups.js\n\t%v", err))
	}
}

// Login takes LoginCredentials and a WebDriver instance and attempts to login to CigarBid
// The WebDriver instance MUST BE returned to the same CigarBid page it was on prior to logging in
func (cbs *cbService) Login() {
	script, err := ioutil.ReadFile("javascript/login.js")
	if err != nil {
		panic(fmt.Errorf("CigarBidService.Login(): Failed to read login.js\n\t%v", err))
	}
	args := make([]interface{}, 2)
	args[0] = cbs.creds.Username
	args[1] = cbs.creds.Password
	result, err := cbs.webDriver.ExecuteScript(string(script), args)
	if cbs.debugMode {
		fmt.Printf("CigarBidService.Login(): Script result: %v\n", result)
	}
	if err != nil {
		panic(fmt.Errorf("CigarBidService.Login(): Failed to run login.js\n\t%v", err))
	}
	if cbs.debugMode {
		url, _ := cbs.webDriver.CurrentURL()
		fmt.Printf("CibarBidService.Login(): Logged in successfully, current URL is %s\n", url)
	}
	cbs.sleep(defaultWaitTimeSec)
	cbs.webDriver.Refresh()
	cbs.sleep(defaultWaitTimeSec)
	cbs.dismissPopups()
	cbs.focusOnPage()
}

// Shutdown stops the selenium service and quits the webdriver instance
func (cbs *cbService) Shutdown() {
	// Ignore errors, just stop and quit
	cbs.seleniumService.Stop()
	cbs.webDriver.Quit()
}
