package helpers

import (
	"encoding/hex"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/sessions"
	"golang.org/x/crypto/bcrypt"
)

const ()

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(formPassword, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(formPassword))
	return err == nil
}

func GetTokenFromSession(r *http.Request) ([]byte, error) {
	session, err := sessions.Get(r)

	if err != nil {
		return nil, err
	}

	decodedSecret, err := hex.DecodeString(session.CSRFSecret)
	if err != nil {
		return nil, err
	}

	return decodedSecret, nil
}

func UserAgentIsBot(userAgent string) bool {
	var botUserAgents = []string{
		"Googlebot",          // Google Bot
		"Bingbot",            // Bing Bot
		"Slurp",              // Yahoo Bot
		"DuckDuckBot",        // DuckDuckGo Bot
		"Bot",                // Generic Bot
		"crawler",            // Generic Crawler
		"spider",             // Generic Spider
		"facebookexternalhit", // Facebook Bot
		"twitterbot",         // Twitter Bot
		"linkedinbot",        // LinkedIn Bot
		"ahrefsbot",          // Ahrefs Bot
		"SEMrushBot",         // SEMRush Bot
		"moz",                // Moz Bot
		"pinterest",          // Pinterest Bot
		"redditbot",          // Reddit Bot
		"baiduspider",        // Baidu Bot
		"yandexbot",          // Yandex Bot
		"exabot",             // ExaBot
		"twitterbot",         // Twitter Bot
		"applebot",           // Apple Bot
		"bingpreview",        // Bing Preview Bot
		"googleusercontent",  // Google User Content Bot
		"curl",               // Curl Bot (often used for automated scripts)
		"wget",               // Wget Bot (often used for automated scripts)
		"archive.org_bot",    // Internet Archive Bot
		"linkdex",            // Linkdex Bot
		"seznambot",          // Seznam Bot
		"conduit",            // Conduit Bot
		"zyborg",             // Zyborg Bot
		"semalt",             // Semalt Bot
		"yahoo",              // Yahoo Bot
		"dotbot",             // Moz Dotbot
		"nimbostratus",       // Nimbostratus Bot
		"surveybot",          // Survey Bot
		"adsbot-google",      // Google Ads Bot
		"searchbot",          // Generic Search Bot
		"scrapy",             // Scrapy Bot
		"sitebot",            // Site Bot
		"webcrawler",         // Web Crawler
		"spiderbot",          // Spider Bot
		"sloth",              // Sloth Bot
		// Add more known bot user-agents as needed
	}
	
	for _, botAgent := range botUserAgents {
		if strings.Contains(userAgent, botAgent) {
			return true
		}
	}
	return false
}
