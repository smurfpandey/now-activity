package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/zenazn/goji"
	"github.com/imroc/req"
)

func main() {
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
	rawResp, _ := req.Get("https://api.listenbrainz.org//1/user/smurfpandey/playing-now")
	resp := rawResp.Response()
	fmt.Println(resp.StatusCode)
	
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
    
    bodyString := string(bodyBytes)
	
	io.WriteString(w, bodyString)
}

// NotFound is a 404 handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
}