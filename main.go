package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
)

var (
	// store a copy of the tree locally
	treeArray      []*structs.Node
	starsProcessed int
)

func nextpos(deltat float64, star structs.Star2D) {
	// ...
}

func initMassCenter(w http.ResponseWriter, r *http.Request) {
	log.Println("[   ] Initializing the Mass Center of all the nodes in the tree")

	// getting the tree index
	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)
	_, _ = fmt.Fprintln(w, treeindex)
}

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

// indexHandler returns a simple overview of the available functions
func indexHandler(w http.ResponseWriter, r *http.Request) {
	infostring := `Hello, this is the simu container!

/
/newpos
/initMassCenter
/calcallforces/{treeindex} Calculates all the forces acting on a star given via a POST
/metrics 
`
	_, _ = fmt.Fprintf(w, infostring)
}

// calcallforcesHandler calculates all the forces acting on a given star
//
// 1. read the tree index
// 2. If the tree is not in the local cache
// 2.1. Get the tree
// 3. Calculate all the forces acting the given theta
// 4. Write the forces to a file
func calcallforcesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The calcallforcesHandler was accessed!")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	log.Println("Read the treeindex")

	// if the endpoint was accessed using a GET request, display how to use the endpoint correctly
	if r.Method == "GET" {
		log.Println("The request method was GET")
		_, _ = fmt.Fprintf(w, "Make a post request to this endpoint to calc some forces!")
		fmt.Println(treeindex)

	} else {
		log.Println("The request method was POST")

		// find out if the tree is all ready cached, if not, cache it
		if isCached(treeindex) == false {
			log.Println("[ ! ] The tree is not in local cache, requesting it from the database")

			// make a http-post request to the databse requesting the tree
			requesturl := fmt.Sprintf("http://%s/dumptree/%d", "db", treeindex)
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

		log.Println("[   ] Getting the star values from the post form")
		// get the star the forces should be calculated on
		x, _ := strconv.ParseFloat(r.PostFormValue("x"), 64)
		y, _ := strconv.ParseFloat(r.PostFormValue("y"), 64)
		vx, _ := strconv.ParseFloat(r.PostFormValue("vx"), 64)
		vy, _ := strconv.ParseFloat(r.PostFormValue("vy"), 64)
		m, _ := strconv.ParseFloat(r.PostFormValue("m"), 64)
		theta, _ := strconv.ParseFloat(r.PostFormValue("theta"), 64)

		log.Println("[   ] Simulator container calcallforces got these values: ")
		log.Printf("(x: %f, y: %f, vx: %f, vy: %f, m: %f, theta: %f)\n", x, y, vx, vy, m, theta)
		_, _ = fmt.Fprintf(w, "calculating forces...\n")
		_, _ = fmt.Fprintf(w, "Simu here, calculating the forces acting on the star (%f, %f)", x, y)

		star := structs.Star2D{
			C: structs.Vec2{
				X: x,
				Y: y,
			},
			V: structs.Vec2{
				X: vx,
				Y: vy,
			},
			M: m,
		}

		// iterate over the tree using Barnes-Hut to determine if the the force should be calculated or not
		log.Printf("[   ] Calculating the forces (%v, *): ", star)

		log.Println(treeArray[treeindex].GenForestTree(treeArray[treeindex]))

		force := treeArray[treeindex].CalcAllForces(star, theta)
		log.Println("[   ] Done Calculating the forces!")
		log.Printf("[   ] The force acting on star %v is %v", star, force)

		_, _ = fmt.Fprintf(w, "The force acting on star %v is %v", star, force)
		writefilerr := ioutil.WriteFile("out.txt", []byte(fmt.Sprintf("%v, %v", star, force)), 0644)
		if writefilerr != nil {
			panic(writefilerr)
		}
	}

	starsProcessed += 1
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

// metricHandler returns a list of the simulators metrics
func metricHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "stars_processed %d", starsProcessed)
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
		log.Printf("Updating the metrics on %s", requestURL)
		log.Printf("key=starsProcessed{hostname=\"%s\"}&value=%d", hostname, starsProcessed)

		defer resp.Body.Close()

		// sleep for a given amount of time
		time.Sleep(time.Second * 5)
	}
}

// randomUpdateStarsProcessed is a test function intended to run in an own go-method. It randomly increases the
// starsProcessed counter faking actual calculations.
func randomUpdateStarsProcessed() {

	// create a new random source to get the random values from
	randomSource := rand.New(rand.NewSource(time.Now().UnixNano()))

	// increase the starsProcessed counter and wait a random amount of time in the interval [0, 10)
	for {
		starsProcessed += 1
		fmt.Printf("Updated starsprocessed: %d\n", starsProcessed)

		// sleep for a random time
		randomInt := randomSource.Intn(10)
		time.Sleep(time.Duration(randomInt) * time.Second)
	}
}

func main() {
	// start a go method pushing the metrics to the manager
	log.Println("[   ] Starting the metric-pusher")
	go pushMetrics("http://manager:80/metrics")

	// randomly update the stars processed counter
	go randomUpdateStarsProcessed()

	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/newpos", calcNewPos).Methods("POST")
	router.HandleFunc("/initMassCenter", calcNewPos).Methods("POST")
	router.HandleFunc("/calcallforces/{treeindex}", calcallforcesHandler).Methods("GET", "POST")

	router.HandleFunc("/metrics", metricHandler).Methods("GET")

	fmt.Println("[   ] Simulator Container up")
	log.Fatal(http.ListenAndServe(":80", router))
}
