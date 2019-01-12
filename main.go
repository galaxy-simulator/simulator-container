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

	log.Println("Simulator container calcNewPos git these values: ")
	log.Printf("(x: %f, y: %f, vx: %f, vy: %f, m: %f)\n", x, y, vx, vy, m)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = fmt.Fprintf(w, "Hello, this is the simu container!")
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/initMassCenter", calcNewPos).Methods("POST")

	// this is an endpoint searching for a star and calculating the forces acting inbetween it and all the other star
	// in it's direct range
	router.HandleFunc("/newpos", calcNewPos).Methods("POST")

	router.HandleFunc("/", indexHandler).Methods("GET")

	log.Fatal(http.ListenAndServe(":8002", router))
}
