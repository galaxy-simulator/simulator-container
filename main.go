package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
	"time"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	// store a copy of the tree locally
	treeArray        []*structs.Node
	starsProcessed   int
	theta            = 0.1
	currentlyCaching = false
)

// isCached returns true if the tree with the given treeindex is cached and false if not
func isCached(treeindex int64) bool {
	// if the specified tree does not have any children and does not contain a star in the root node
	if int(treeindex+1) <= len(treeArray) {
		return true
	} else {
		return false
	}
}

// cache caches the tree with the given index
func cache(treeindex int64) {
	currentlyCaching = true

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

	// if the treeArray is not long enough, fill it
	for int(treeindex) > len(treeArray) {
		emptyNode := structs.NewNode(structs.NewBoundingBox(structs.NewVec2(0, 0), 10))
		treeArray = append(treeArray, emptyNode)
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
func processstars(url string, core int) {

	// infinitely get stars and calculate the forces acting on them
	for {

		log.Println("------------------------------------------")

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
			if isCached(stargalaxy.Index) == false || currentlyCaching == false {
				log.Println("[Caching]")
				cache(stargalaxy.Index)
				log.Println("[Done Caching!]")
			}

			// calculate the forces acting inbetween all the stars in the galaxy
			star := stargalaxy.Star
			galaxyindex := stargalaxy.Index

			log.Printf("[+++] %v\n", star)

			force := calcallforces(star, galaxyindex)

			// calculate the new position
			star.CalcNewPos(force, 1e15)

			log.Printf("[+++] %v\n", star)

			// insert the "new" star into the next timestep
			insertStar(star, galaxyindex+1)

			// increase the starProcessed counter
			starsProcessed += 1
			log.Printf("[%d][%d] Processed as star!\n", core, starsProcessed)

			// time.Sleep(time.Second * 100)
		} else {
			fmt.Println("Could not get a star from the manager!")
			// Sleep a second and try again
			time.Sleep(time.Second * 1)
		}

	}
}

// insertStar inserts the given star into the tree with the given index
func insertStar(star structs.Star2D, galaxyindex int64) {
	log.Println("[ i ] Inserting the star into the galaxy")

	// if the galaxy does not exist yet, create it
	requrl := "http://db.nbg1.emile.space/nrofgalaxies"
	resp, err := http.Get(requrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	currentAmounfOfGalaxies, _ := strconv.ParseInt(string(body), 10, 64)

	log.Printf("[ i ] Current amount of galaxies: %d\n", currentAmounfOfGalaxies)

	// if there is no galaxy to insert the star into, create it
	if currentAmounfOfGalaxies <= galaxyindex {
		log.Println("[ i ] There is no galaxy to insert into -> creating one")

		// create the new galaxy using a post request
		requestURL := "http://db.nbg1.emile.space/new"
		resp, err := http.PostForm(requestURL,
			url.Values{
				"w": {fmt.Sprintf("%d", 1000000000)},
			},
		)
		if err != nil {
			fmt.Printf("Cound not make a POST request to %s", requestURL)
			body, _ := ioutil.ReadAll(resp.Body)
			fmt.Printf(string(body))
		}
		defer resp.Body.Close()

		log.Println("[ i ] Done creating the new galaxy")
	}

	// insert the star into the galaxy

	log.Printf("[ i ] Inserting the star into galaxy nr. %d\n", galaxyindex)

	// create the new galaxy using a post request
	requestURL := fmt.Sprintf("http://db.nbg1.emile.space/insert/%d", galaxyindex)
	resppost, errpost := http.PostForm(requestURL,
		url.Values{
			"x":  {fmt.Sprintf("%f", star.C.X)},
			"y":  {fmt.Sprintf("%f", star.C.Y)},
			"vx": {fmt.Sprintf("%f", star.V.X)},
			"vy": {fmt.Sprintf("%f", star.V.Y)},
			"m":  {fmt.Sprintf("%f", star.M)},
		},
	)
	if errpost != nil {
		panic(errpost)
		// handle error
	}
	defer resppost.Body.Close()

	log.Printf("[<<<] \t %v", string(requestURL))
	log.Printf("[<<<] \t %v", string(body))

	log.Println("[ i ] Done inserting the star")
}

// calcallforces calculates the forces acting on a given star using the given
// treeindex to define which other stars are in the galaxy
func calcallforces(star structs.Star2D, treeindex int64) structs.Vec2 {
	fmt.Printf("star: %v", star)
	fmt.Printf("treeindex: %d", treeindex)
	fmt.Printf("len(treeArray): %d", len(treeArray))
	force := treeArray[treeindex].CalcAllForces(star, theta)
	return force
}

func main() {
	// start a go method pushing the metrics to the manager
	log.Println("[   ] Starting the metric-pusher")
	go pushMetrics("http://manager.nbg1.emile.space/metrics")

	numCPU := runtime.NumCPU()

	log.Printf("Starting %d go threads", numCPU)

	// for i := 0; i < numCPU-1; i++ {
	// 	go processstars("http://manager.nbg1.emile.space/providestars/0", i)
	// }
	processstars("http://manager.nbg1.emile.space/providestars/0", 7)
}
