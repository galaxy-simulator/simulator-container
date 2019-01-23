package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	// store a copy of the tree locally
	treeArray      []*structs.Node
	starsProcessed int
	theta          = 0.1
)

// calcNewPos calculates the new position of the star it receives via a POST request
// TODO: Implement it
func calcNewPos(w http.ResponseWriter, r *http.Request) {
	// get the post parameters
	x, _ := strconv.ParseFloat(r.PostFormValue("x"), 64)
	y, _ := strconv.ParseFloat(r.PostFormValue("y"), 64)
	vx, _ := strconv.ParseFloat(r.PostFormValue("vx"), 64)
	vy, _ := strconv.ParseFloat(r.PostFormValue("vy"), 64)
	m, _ := strconv.ParseFloat(r.PostFormValue("m"), 64)

	log.Println("Simulator container calcNewPos got these values: ")
	log.Printf("(x: %f, y: %f, vx: %f, vy: %f, m: %f)\n", x, y, vx, vy, m)
}

// isCached returns true if the tree with the given treeindex is cached and false if not
func isCached(treeindex int64) bool {
	log.Printf("[isCached] Testing if %d is in the local cache\n", treeindex+1)
	log.Printf("[isCached] TreeArray length: %d\n", len(treeArray))

	// if the specified tree does not have any children and does not contain a star in the root node
	if int(treeindex+1) <= len(treeArray) {
		log.Println("[isCached] Yes it is!")
		return true
	} else {
		log.Println("[isCached] Doesn't seem so")
		return false
	}
}

func cache(treeindex int64) {
	log.Println("[ ! ] The tree is not in local cache, requesting it from the database")

	// make a http-post request to the databse requesting the tree
	requesturl := fmt.Sprintf("http://db/dumptree/%d", treeindex)
	log.Println("[   ] Requesting the tree from the database")
	resp, err := http.Get(requesturl)
	if err != nil {
		panic(err)
	}
	log.Println("[   ] No error occurred!")
	defer resp.Body.Close()

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		panic(readerr)
	}

	log.Println("[   ] Unmarshaling the tree and storing it the treeArray")
	tree := &structs.Node{}
	jsonUnmarshalErr := json.Unmarshal(body, tree)
	if jsonUnmarshalErr != nil {
		panic(jsonUnmarshalErr)
	}
	log.Println("[   ] No error occurred!")
	treeArray = append(treeArray, tree)
}

// pushMetrics pushes the metrics to the given host
func pushMetrics(host string) {

	// start an infinite loop
	for {

		hostname, _ := os.Hostname()

		// define a post-request and send it to the given host
		requestURL := fmt.Sprintf("%s", host)
		resp, err := http.PostForm(requestURL,
			url.Values{
				"key":   {fmt.Sprintf("%s{hostname=\"%s\"}", "starsProcessed", hostname)},
				"value": {fmt.Sprintf("%d", starsProcessed)},
			},
		)
		if err != nil {
			fmt.Printf("Cound not make a POST request to %s", requestURL)
		}
		log.Printf("[metrics] Updating the metrics on %s", requestURL)
		log.Printf("[metrics] key=starsProcessed{hostname=\"%s\"}&value=%d", hostname, starsProcessed)

		defer resp.Body.Close()

		// sleep for a given amount of time
		time.Sleep(time.Second * 5)
	}
}

// processstars processes stars as long as the sun is shining!
func processstars(url string) {

	// infinitely get stars and calculate the forces acting on them
	for {

		log.Println("[   ] Getting a star from the manager")
		// make a request to the given url and get the stargalaxy
		resp, err := http.Get(url)
		if err != nil {
			fmt.Println("PANIC")
			panic(err)
		}
		defer resp.Body.Close()
		log.Println("[   ] Done")

		// read the response containing a list of all stars in json format
		log.Println("[   ] Reading the content")
		body, err := ioutil.ReadAll(resp.Body)
		log.Println("[   ] Done")

		// if the response body is not a "Bad Gateway", continue.
		// This problem occurs, when the manager hasn't got enough stars
		if string(body) != "Bad Gateway" {
			stargalaxy := &structs.Stargalaxy{}

			// unmarshal the stargalaxy
			log.Println("[   ] Unmarshaling the stargalaxy")
			unmarshalErr := json.Unmarshal(body, stargalaxy)
			if unmarshalErr != nil {
				panic(unmarshalErr)
			}
			log.Println("[   ] Done")
			log.Printf("[Star] (%f, %f)", stargalaxy.Star.C.X, stargalaxy.Star.C.Y)

			// if the galaxy is not cached yet, cache it
			log.Println("[   ] Testing is the galaxy is cached or not")
			if isCached(stargalaxy.Index) == false {
				log.Println("[   ] It is not -> caching")
				cache(stargalaxy.Index)
			}
			log.Println("[   ] Done")

			log.Println("[   ] Calculating the forces acting")
			// calculate the forces acting inbetween all the stars in the galaxy
			star := stargalaxy.Star
			galaxyindex := stargalaxy.Index

			calcallforces(star, galaxyindex)

			log.Println("[   ] Done")

			// insert the "new" star into the next timestep

			log.Println("[   ] Calculating the new position")
			log.Println("[   ] TODO")

			// increase the starProcessed counter
			starsProcessed += 1

			log.Println("[   ] Waiting 10 seconds...")
			// time.Sleep(time.Second * 100)
			log.Println("[   ] Done")
		} else {
			// Sleep a second and try again
			time.Sleep(time.Second * 1)
		}

	}
}

// calcallforces calculates the forces acting on a given star using the given
// treeindex to define which other stars are in the galaxy
func calcallforces(star structs.Star2D, treeindex int64) {

	// iterate over the tree using Barnes-Hut to determine if the the force should be calculated or not
	log.Printf("[   ] Calculating the forces (%v, *): ", star)

	force := treeArray[treeindex].CalcAllForces(star, theta)
	log.Println("[   ] Done Calculating the forces!")
	log.Printf("[FORCE] Force acting on star: %v \t -> %v", star, force)
}

func main() {
	// start a go method pushing the metrics to the manager
	log.Println("[   ] Starting the metric-pusher")
	go pushMetrics("http://manager/metrics")

	processstars("http://manager/providestars/0")
}
