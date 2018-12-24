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
	z, _ := strconv.ParseFloat(r.PostFormValue("z"), 64)

	log.Printf("(%f, %f, %f)", x, y, z)
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/initMassCenter", calcNewPos).Methods("POST")

	// this is an endpoint searching for a star and calculating the forces acting inbetween it and all the other star
	// in it's direct range
	router.HandleFunc("/newpos", calcNewPos).Methods("POST")
	log.Fatal(http.ListenAndServe(":8345", router))
}
