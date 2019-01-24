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
	// if the specified tree does not have any children and does not contain a star in the root node
	if int(treeindex+1) <= len(treeArray) {
		return true
	} else {
		return false
	}
}

func cache(treeindex int64) {
	// make a http-post request to the databse requesting the tree
	requesturl := fmt.Sprintf("http://db.nbg1.emile.space/dumptree/%d", treeindex)
	resp, err := http.Get(requesturl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, readerr := ioutil.ReadAll(resp.Body)
	if readerr != nil {
		panic(readerr)
	}

	tree := &structs.Node{}
	jsonUnmarshalErr := json.Unmarshal(body, tree)
	if jsonUnmarshalErr != nil {
		panic(jsonUnmarshalErr)
	}
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

		defer resp.Body.Close()

		// sleep for a given amount of time
		time.Sleep(time.Second * 5)
	}
}

// processstars processes stars as long as the sun is shining!
func processstars(url string) {

	// infinitely get stars and calculate the forces acting on them
	for {

		// make a request to the given url and get the stargalaxy
		resp, err := http.Get(url)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		// read the response containing a list of all stars in json format
		body, err := ioutil.ReadAll(resp.Body)

		// if the response body is not a "Bad Gateway", continue.
		// This problem occurs, when the manager hasn't got enough stars
		if string(body) != "Bad Gateway" {
			stargalaxy := &structs.Stargalaxy{}

			// unmarshal the stargalaxy
			unmarshalErr := json.Unmarshal(body, stargalaxy)
			if unmarshalErr != nil {
				panic(unmarshalErr)
			}

			// if the galaxy is not cached yet, cache it
			if isCached(stargalaxy.Index) == false {
				cache(stargalaxy.Index)
			}

			// calculate the forces acting inbetween all the stars in the galaxy
			star := stargalaxy.Star
			galaxyindex := stargalaxy.Index

			calcallforces(star, galaxyindex)

			// calculate the new position
			// insert the "new" star into the next timestep

			// increase the starProcessed counter
			starsProcessed += 1

			// time.Sleep(time.Second * 100)
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

	force := treeArray[treeindex].CalcAllForces(star, theta)
	log.Printf("[FORCE] Force acting on star: %v \t -> %v", star, force)
}

func main() {
	// start a go method pushing the metrics to the manager
	log.Println("[   ] Starting the metric-pusher")
	go pushMetrics("http://manager.nbg1.emile.space/metrics")

	processstars("http://manager.nbg1.emile.space/providestars/0")
}
