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
	SetLocal() (interface{}, error)
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
	if newCbService.seleniumService, err = seleniumwindowscompatibility.NewSeleniumService(newOpts.SeleniumPath, newOpts.Port, sOpts...); err == nil {
		// Connect to the WebDriver instance running locally
		caps := seleniumwindowscompatibility.Capabilities{"browserName": newOpts.BrowserName}
		if newCbService.webDriver, err = seleniumwindowscompatibility.NewRemote(caps, fmt.Sprintf("http://localhost:%d/wd/hub", port)); err == nil {
			if err = newCbService.NavigateToCigarBid(); err == nil {
				return newCbService, err // err will be nil
			}
		}
	}
	// if executiuon made it this far, then there was an error
	return nil, err
}

func (cbs *cbService) NavigateToCigarBid() error {
	var url string
	var err error
	// Check current URL to make sure we're on Cigarbid's website
	if url, err = cbs.webDriver.CurrentURL(); err == nil {
		if strings.Contains(url, cigarBidURL) {
			if cbs.debugMode {
				fmt.Println("CigarBidService.NavigateToCigarBid(): You are on Cigarbid's page")
			}
		} else {
			if cbs.debugMode {
				fmt.Printf("CigarBidService.NavigateToCigarBid(): You are NOT on Cigarbid's page, your URL is %s\n", url)
				fmt.Println("\tTaking you to Cigar Bid...")
			}
			if err = cbs.webDriver.Get(cigarBidURL); err == nil {
				if err = cbs.DisablePushNotifications(); err == nil {
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
	if result, err = cbs.SetLocal("isPushNotificationsEnabled", "false"); err == nil {
		if cbs.debugMode {
			fmt.Printf("CigarBidService.DisablePushNotifications(): result: %v", result)
		}
	}
	return err
}

// AcceptCookies either sets the Local key for the cookie popup,
// or it will click the Accept button if the prompt has appeared
func (cbs *cbService) AcceptCookies() error {
	var err error = nil
	if elem, err := cbs.webDriver.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#onesignal-slidedown-dialog"); err == nil {
		if btn, err := elem.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#onesignal-slidedown-allow-button"); err == nil {
			if err = btn.Click(); err == nil {
				cbs.cookiesAccepted = true
			}
		}
	}
	if err != nil {
		// handle no such element case
		if e, ok := err.(*seleniumwindowscompatibility.Error); ok {
			if e.Err == "no such element" {
				// set Local key instead
				if result, err := cbs.SetLocal("onesignal-notification-prompt", fmt.Sprintf("{\"value\":\"\\\"dismissed\\\"\",\"timestamp\":%d}", (time.Now().Unix()*1000))); err == nil {
					if cbs.debugMode {
						fmt.Printf("CigarBidService.AcceptCookies(): Set onesignal-notification-prompt to dismiss popup, result: %v", result)
					}
					cbs.cookiesAccepted = true
				}
			}
		} else {
			// neither attempt to accept cookies succeeded
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
	if err := cbs.NavigateToCigarBid(); err == nil {
		// Get a reference to the Page Container with the Sign In button
		if elem, err := cbs.webDriver.FindElement(seleniumwindowscompatibility.ByCSSSelector, "#page-container"); err == nil {
			if btn, err := elem.FindElement(seleniumwindowscompatibility.ByCSSSelector, ".btn.btn-success.boostbar-login"); err == nil {
				err = btn.Click()
			}
		}
	}
	return err
}

// Shutdown stops the selenium service and quits the webdriver instance
func (cbs *cbService) Shutdown() error {
	var err error = nil
	if err = cbs.seleniumService.Stop(); err == nil {
		err = cbs.webDriver.Quit()
	}
	return err
}

func (cbs *cbService) SetLocal(key, value string) (interface{}, error) {
	args := make([]interface{}, 2)
	args[0] = key
	args[1] = value
	var result interface{}
	var err error = nil
	if result, err = cbs.webDriver.ExecuteScript("localStorage.setItem(arguments[0],arguments[1])", args); err == nil {
		if cbs.debugMode {
			fmt.Printf("Script Completed Successfully, result: %v", result)
		}
	}
	return result, err
}
