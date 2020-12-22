package common

// LoginCredentials stores username and password
type LoginCredentials struct {
	Username string
	Password string
}

// SeleniumOptions provides a mapper for Selenium's ServiceOption
type SeleniumOptions struct {
	SeleniumPath     string
	ChromeDriverPath string
	GeckoDriverPath  string
	BrowserName      string
	Port             int
	Debug            bool
}

// FreefallBuyOrder records attributes for purchase order
type FreefallBuyOrder struct {
	cigarurl   string  // URL to the product page (must be freefall)
	buylimit   float64 // MAXIMUM price at which to execute a buy order
	buyamount  int     // Number of units to buy
	expiration int64   // Expiration Timestamp in Unix Epoch Time (seconds since Unix Epoch)
}
