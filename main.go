package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"git.darknebu.la/GalaxySimulator/structs"
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

	// Recursively calculate the center of mass of a node by calculating the center of mass of the four children of
	// that node.
	// The formula used is the following:
	// ( (x_i * m_i)/(m_total) , (y_i * m_i)/(m_total) )

	// starting at the root node of the tree

}

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

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, this is the simu container!")
}

func calcallforcesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("The calcallforcesHandler was accessed!")

	vars := mux.Vars(r)
	treeindex, _ := strconv.ParseInt(vars["treeindex"], 10, 0)

	if r.Method == "GET" {
		_, _ = fmt.Fprintf(w, "Make a post request to this endpoint to calc some forces!")
		fmt.Println(treeindex)
	} else {
		// get the post parameters
		x, _ := strconv.ParseFloat(r.PostFormValue("x"), 64)
		y, _ := strconv.ParseFloat(r.PostFormValue("y"), 64)
		vx, _ := strconv.ParseFloat(r.PostFormValue("vx"), 64)
		vy, _ := strconv.ParseFloat(r.PostFormValue("vy"), 64)
		m, _ := strconv.ParseFloat(r.PostFormValue("m"), 64)

		log.Println("Simulator container calcallforces got these values: ")
		log.Printf("(x: %f, y: %f, vx: %f, vy: %f, m: %f)\n", x, y, vx, vy, m)
		_, _ = fmt.Fprintf(w, "calculating forces...")
		_, _ = fmt.Fprintf(w, "Simu here, calculating the forces acting on the star (%f, %f)", x, y)
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/", indexHandler).Methods("GET")
	router.HandleFunc("/newpos", calcNewPos).Methods("POST")
	router.HandleFunc("/initMassCenter", calcNewPos).Methods("POST")
	router.HandleFunc("/calcallforces/{treeindex}", calcallforcesHandler).Methods("GET", "POST")

	log.Fatal(http.ListenAndServe(":80", router))
}
