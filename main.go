package main

import (
	"io"
	"io/ioutil"
	"encoding/json"
	"net/http"
	"time"

	"github.com/patrickmn/go-cache"
	"github.com/imroc/req"
	"github.com/zenazn/goji"
)

type Game struct {
	Name      string `json:"name"`
	Exec      string `json:"exec_name"`
	Website   string `json:"website_url"`
	Store     string `json:"store_url"`
	StartedAt int    `json:"started_at"` 
}

var cacheDuration time.Duration
var myCache *cache.Cache
const (
	CACHE_KEY_MUSIC = "wow-music"
	CACHE_KEY_GAME = "wow-game"
)


func main() {

	// Create a cache with a default expiration time of 3 minutes, and which
	// purges expired items every 5 minutes
	cacheDuration = 3*time.Minute
	myCache = cache.New(cacheDuration, 5*time.Minute)

	// Add routes to the global handler
	goji.Get("/listening", whatMusic)
	goji.Get("/healthcheck", healthCheck)

	goji.Get("/playing", whatGame)
	goji.Post("/playing", startedGame)
	goji.Delete("/playing", closedGame)

	// Use a custom 404 handler
	goji.NotFound(NotFound)

	// Call Serve() at the bottom of your main() function, and it'll take
	// care of everything else for you, including binding to a socket (with
	// automatic support for systemd and Einhorn) and supporting graceful
	// shutdown on SIGINT. Serve() is appropriate for both development and
	// production.
	goji.Serve()
}

// whatMusic route (GET "/"). Print a list of greets.
func whatMusic(httpResp http.ResponseWriter, httpReq *http.Request) {

	httpResp.Header().Set("Content-Type", "application/json")

	musicData, expiryTime, found := myCache.GetWithExpiration(CACHE_KEY_MUSIC)
	if found {
		httpResp.Header().Set("Expires", getHTTPTime(expiryTime))
		io.WriteString(httpResp, musicData.(string))
		return
	}

	rawResp, _ := req.Get("https://api.listenbrainz.org/1/user/smurfpandey/playing-now")
	resp := rawResp.Response()
		
	bodyBytes, _ := ioutil.ReadAll(resp.Body)
    
	bodyString := string(bodyBytes)
	
	// save in cache
	myCache.Set(CACHE_KEY_MUSIC, bodyString, cache.DefaultExpiration)
	
	expiryTime = time.Now().Add(cacheDuration)
	httpResp.Header().Set("Expires", getHTTPTime(expiryTime))
	io.WriteString(httpResp, bodyString)
}

// NotFound is a 404 handler.
func NotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Umm... have you tried turning it off and on again?", 404)
}

// Healthcheck handler
func healthCheck(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Ok")
	return
}

func getHTTPTime(yourTime time.Time) string {
	return yourTime.UTC().Format(http.TimeFormat)
}

func startedGame(httpResp http.ResponseWriter, httpReq *http.Request) {
	decoder := json.NewDecoder(httpReq.Body)
    var gameBeingPlayed Game
    err := decoder.Decode(&gameBeingPlayed)
    if err != nil {
        http.Error(httpResp, "Umm... are you sure this is the correct data?", 400)
	}
	
	myCache.Set(CACHE_KEY_GAME, &gameBeingPlayed, 5*time.Hour)
	
	io.WriteString(httpResp, "Ok")
}

func whatGame(httpResp http.ResponseWriter, httpReq *http.Request) {
	cacheData, found := myCache.Get(CACHE_KEY_GAME)
	if found {
		gameBeingPlayed := cacheData.(*Game)
		strGame, err := json.Marshal(gameBeingPlayed)
		if err != nil {
			http.Error(httpResp, "Ohh... game data is corrupted?", 500)
			return
		}

		httpResp.Header().Set("Content-Type", "application/json")

		io.WriteString(httpResp, string(strGame))
		return
	}

	io.WriteString(httpResp, "Nothing")
}

func closedGame(httpResp http.ResponseWriter, httpReq *http.Request) {
	myCache.Delete(CACHE_KEY_GAME)
	io.WriteString(httpResp, "Ok")
}