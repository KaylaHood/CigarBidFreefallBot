package main

type LoginCredentials struct {
	username string
	password string
}

// Connects to Cigarbid website and provides methods
// for targeting a specific freefall product, determining
// the lowest price, and making a purchase at that price

type CigarBidFreefallService interface {
	Login(creds LoginCredentials) error
	Logout() error
	FreefallMinimum(url string) (float64, error)
	FreefallPurchaseAtPrice(url string, maxprice float64) (bool, error)
}
