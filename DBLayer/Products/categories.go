package products

import (
	"encoding/json"
	"net/http"
)


type CategoryCreation struct {
    CategoryName        string
    CategoryDescription string
}

func (prd ProductRoutesTray) AddPrimeCategory(w http.ResponseWriter, r *http.Request) {
	
	decoder := json.NewDecoder(r.Body)
	var category CategoryCreation
	err := decoder.Decode(&category)

	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}


	// Logic to add the prime category to the database


}

func AddSubCategory(w http.ResponseWriter, r *http.Request) {

}

func AddFinalCategory(w http.ResponseWriter, r *http.Request) {

}
