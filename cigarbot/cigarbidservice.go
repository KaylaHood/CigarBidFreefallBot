package cigarbot

import (
	"fmt"
	"os"
	"strings"
	"time"

	common "github.com/KaylaHood/CigarBidFreefallBot"
	"github.com/KaylaHood/CigarBidFreefallBot/seleniumwindowscompatibility"
)

const (
	cigarBidURL = "https://www.cigarbid.com/"
)

// CigarBidService provides utilities for managing session with CigarBid site
type CigarBidService interface {
	NavigateToCigarBid() error
	DisablePushNotifications() error
	AcceptCookies() error
	Login() error
	Shutdown() error
	SetLocal(string, string) (interface{}, error)
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
func NewCigarBidService(newCreds common.LoginCredentials, newOpts common.SeleniumOptions) (CigarBidService, error) {
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
	if err == nil {
		// Connect to the WebDriver instance running locally
		caps := seleniumwindowscompatibility.Capabilities{"browserName": newOpts.BrowserName}
		newCbService.webDriver, err = seleniumwindowscompatibility.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port))
		if err == nil {
			err = newCbService.NavigateToCigarBid()
		}
	}
	// if executiuon made it this far, then there was an error
	return newCbService, err
}

func (cbs *cbService) NavigateToCigarBid() error {
	var url string
	var err error
	// Check current URL to make sure we're on Cigarbid's website
	url, err = cbs.webDriver.CurrentURL()
	if err == nil {
		if strings.Contains(url, cigarBidURL) {
			if cbs.debugMode {
				fmt.Println("CigarBidService.NavigateToCigarBid(): You are on Cigarbid's page")
			}
		} else {
			if cbs.debugMode {
				fmt.Printf("CigarBidService.NavigateToCigarBid(): You are NOT on Cigarbid's page, your URL is %s\n", url)
				fmt.Println("\tTaking you to Cigar Bid...")
			}
			err = cbs.webDriver.Get(cigarBidURL)
			if err == nil {
				err = cbs.DisablePushNotifications()
				if err == nil {
					err = cbs.AcceptCookies()
				}
			}
		}
	}
	return err
}

func (cbs *cbService) DisablePushNotifications() error {
	var result interface{}
	var err error = nil
	result, err = cbs.SetLocal("isPushNotificationsEnabled", "false")
	if err == nil {
		if cbs.debugMode {
			fmt.Printf("CigarBidService.DisablePushNotifications(): result: %v\n", result)
		}
	}
	return err
}

// AcceptCookies either sets the Local key for the cookie popup,
// or it will click the Accept button if the prompt has appeared
func (cbs *cbService) AcceptCookies() error {
	var err error = nil
	if cbs.debugMode {
		fmt.Println("CigarBidService.AcceptCookies(): Entered function")
	}
	elem, err := cbs.webDriver.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#onesignal-slidedown-dialog")
	if err == nil {
		if cbs.debugMode {
			fmt.Println("CigarBidService.AcceptCookies(): Found base element")
		}
		btn, err := elem.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#onesignal-slidedown-allow-button")
		if err == nil {
			if cbs.debugMode {
				fmt.Println("CigarBidService.AcceptCookies(): Found button element")
			}
			err = btn.Click()
			if err == nil {
				if cbs.debugMode {
					fmt.Println("CigarBidService.AcceptCookies(): Button clicked")
				}
				cbs.cookiesAccepted = true
			}
		}
	}
	if err != nil {
		// handle no such element case
		if e, ok := err.(*seleniumwindowscompatibility.Error); ok {
			if e.Err == "no such element" {
				// reset err
				err = nil
				// set Local key instead
				result, err := cbs.SetLocal("onesignal-notification-prompt", fmt.Sprintf("{\"value\":\"\\\"dismissed\\\"\",\"timestamp\":%d}", (time.Now().Unix()*1000)))
				if err == nil {
					if cbs.debugMode {
						fmt.Printf("CigarBidService.AcceptCookies(): Set onesignal-notification-prompt to dismiss popup, result: %v\n", result)
					}
					cbs.cookiesAccepted = true
				}
			} else {
				// unrecognized error
				fmt.Printf("CigarBidService.AcceptCookies(): Unrecognized error: %v\n", err)
			}
		} else {
			// neither attempt to accept cookies succeeded
			if cbs.debugMode {
				fmt.Printf("CigarBidService.AcceptCookies(): Error: %v\n", err)
			}
			cbs.cookiesAccepted = false
			return err
		}
	}
	return err // should be nil
}

// Login takes LoginCredentials and a WebDriver instance and attempts to login to CigarBid
// The WebDriver instance MUST BE returned to the same CigarBid page it was on prior to logging in
func (cbs *cbService) Login() error {
	var err error = nil
	err = cbs.NavigateToCigarBid()
	if err == nil {
		// Get a reference to the Page Container with the Sign In button
		elem, err := cbs.webDriver.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#page-container")
		if err == nil {
			btn, err := elem.FindElement(seleniumwindowscompatibility.ByCSSSelector, ".btn.btn-success.boostbar-login")
			if err == nil {
				err = btn.Click()
			}
		}
	}
	return err
}

// Shutdown stops the selenium service and quits the webdriver instance
func (cbs *cbService) Shutdown() error {
	var err error = nil
	err = cbs.seleniumService.Stop()
	if err == nil {
		err = cbs.webDriver.Quit()
	}
	if err == nil && cbs.debugMode {
		fmt.Println("CigarBidService.Shutdown(): successfully shut down service")
	}
	return err
}

func (cbs *cbService) SetLocal(key, value string) (interface{}, error) {
	args := make([]interface{}, 2)
	args[0] = key
	args[1] = value
	var result interface{}
	var err error = nil
	result, err = cbs.webDriver.ExecuteScript("localStorage.setItem(arguments[0],arguments[1])", args)
	if err == nil {
		if cbs.debugMode {
			fmt.Printf("Script Completed Successfully, result: %v\n", result)
		}
	}
	return result, err
}
