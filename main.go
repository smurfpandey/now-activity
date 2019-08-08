package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/imroc/req"
	"github.com/zenazn/goji"
)

var myCache *cache.Cache

func main() {

	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	myCache = cache.New(3*time.Minute, 6*time.Minute)

	// Add routes to the global handler
	goji.Get("/whats-playing", whatsPlaying)

	// Use a custom 404 handler
	goji.NotFound(NotFound)

	// Call Serve() at the bottom of your main() function, and it'll take
	// care of everything else for you, including binding to a socket (with
	// automatic support for systemd and Einhorn) and supporting graceful
	// shutdown on SIGINT. Serve() is appropriate for both development and
	// production.
	goji.Serve()
}

// whatsPlaying route (GET "/"). Print a list of greets.
func whatsPlaying(w http.ResponseWriter, r *http.Request) {

	foo, found := myCache.Get("wow-music")
	if found {
		io.WriteString(w, foo.(string))
		return
	}

	rawResp, _ := req.Get("https://api.listenbrainz.org//1/user/smurfpandey/playing-now")
	resp := rawResp.Response()
		
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
    
	bodyString := string(bodyBytes)
	
	// save in cache
	myCache.Set("wow-music", bodyString, cache.DefaultExpiration)
	
	io.WriteString(w, bodyString)
}

// NotFound is a 404 handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
}