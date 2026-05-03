package adminendpoints

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Apouzi/Golang-Admin-Service/helpers"
	"github.com/go-chi/chi/v5"
)

type AdminCategoriesRoutes struct{
	DB *sql.DB
}

func InstanceAdminCategoriesRoutes(  ) *AdminCategoriesRoutes {
	r := &AdminCategoriesRoutes{
	}
	return r
}









func (route *AdminCategoriesRoutes) ReturnCategoryTree(w http.ResponseWriter, r *http.Request){
	 catTree := CategoryTree{}
	 url := "http://dblayer:8080/category/tree"
	 resp, err := http.Get(url)
	 if err != nil {
		 fmt.Println("Error fetching category tree:", err)
		 helpers.ErrorJSON(w, errors.New("failed to fetch category tree"), 500)
		 return
	 }
	 defer resp.Body.Close()

	 err = json.NewDecoder(resp.Body).Decode(&catTree)
	 if err != nil {
		 fmt.Println("Error decoding category tree:", err)
		 helpers.ErrorJSON(w, errors.New("failed to decode category tree"), 500)
		 return
	 }

	 helpers.WriteJSON(w, 200, catTree)
}

func (route *AdminRoutes) DeletePrimeCategory(w http.ResponseWriter, r *http.Request){
	CatID := chi.URLParam(r,"CatPrimeID")
	if CatID == ""{
		fmt.Println("No CatPrimeID wasn't pulled")
		return
	}
	url := "http://dblayer:8080/category/prime/" + CatID
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("Failed to create delete request for CatPrimeID")
		helpers.ErrorJSON(w, errors.New("failed to create delete request"), 500)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeID")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed deletion in CatPrimeID, status code:", resp.StatusCode)
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:true}
	helpers.WriteJSON(w,200,sendBack)
}
func (route *AdminRoutes) DeleteSubCategory(w http.ResponseWriter, r *http.Request){
	CatID := chi.URLParam(r,"CatSubID")
	if CatID == ""{
		fmt.Println("No CatSubID wasn't pulled")
		return
	}
	
	url := "http://dblayer:8080/category/sub/" + CatID
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("Failed to create delete request for CatSubID")
		helpers.ErrorJSON(w, errors.New("failed to create delete request"), 500)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("Failed deletion in CatSubID")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed deletion in CatSubID, status code:", resp.StatusCode)
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	
	sendBack := DeletedSendBack{SendBack:true}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *AdminRoutes) DeleteFinalCategory(w http.ResponseWriter, r *http.Request){
	CatID := chi.URLParam(r,"CatFinalID")
	if CatID == ""{
		fmt.Println("No CatFinalID wasn't pulled")
		return
	}
	
	url := "http://dblayer:8080/category/final/" + CatID
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		fmt.Println("Failed to create delete request for CatFinalID")
		helpers.ErrorJSON(w, errors.New("failed to create delete request"), 500)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil{
		fmt.Println("Failed deletion in CatFinalID")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Failed deletion in CatFinalID, status code:", resp.StatusCode)
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:true}
	helpers.WriteJSON(w,200,sendBack)
}




func (route *AdminCategoriesRoutes) CreatePrimeCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	url := "http://dblayer:8080/category/prime"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w,  200, catret)
}

func (route *AdminCategoriesRoutes) CreateSubCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println("error reading json",err)
	}

	fmt.Println("category payload:", category_read.CategoryName, category_read.CategoryDescription)


	url := "http://dblayer:8080/category/sub"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println("error marshalling json",err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}

func (route *AdminCategoriesRoutes) CreateFinalCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/final"
	catBytes, err := json.Marshal(category_read)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CategoryReturn{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}



func (route *AdminCategoriesRoutes) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed reading json"), 500)
		return
	}
	fmt.Println("final sub:",FinalSub)

	url := "http://dblayer:8080/category/primetosub"
	catBytes, err := json.Marshal(FinalSub)
	if err != nil{
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed marshalling json"), 500)
		return
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
		helpers.ErrorJSON(w, errors.New("failed posting json"), 500)
		return
	}

	catret := &CatToCat{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
		helpers.ErrorJSON(w, errors.New("failed decoding json"), 500)
		return	
	}
	
	helpers.WriteJSON(w, http.StatusCreated, catret)
}

func (route *AdminCategoriesRoutes) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w,r, &FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/subtofinal"
	catBytes, err := json.Marshal(FinalSub)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CatToCat{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}

	helpers.WriteJSON(w, http.StatusAccepted, catret)
}


func (route *AdminCategoriesRoutes) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request){
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalProd := CatToProd{}
	err := helpers.ReadJSON(w,r, &FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	url := "http://dblayer:8080/category/finaltoprod"
	catBytes, err := json.Marshal(FinalProd)
	if err != nil{
		fmt.Println(err)
	}
	catDecode:= bytes.NewReader(catBytes)
	resp, err := http.Post(url,"application/json",catDecode)
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	catret := &CatToProd{}
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(catret)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}

	helpers.WriteJSON(w, http.StatusAccepted, FinalProd)
}


func (route *AdminCategoriesRoutes) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request){
	ReadCatR := ReadCat{}
	err := helpers.ReadJSON(w,r,&ReadCatR)
	if err != nil{
		fmt.Println(err)
	}
	fmt.Println("InsertIntoCategory ReadCatR",ReadCatR)
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
	route.DB.Exec(FinalProd, 1,ReadCatR.Category)
}



func (route *AdminCategoriesRoutes) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
	url := "http://dblayer:8080/category/prime"
	resp, err := http.Get(url)
	if err != nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	categoryList := []CategoryReturn{}
	err = json.NewDecoder(resp.Body).Decode(&categoryList)
	if err != nil {
		fmt.Println("Error decoding prime categories:", err)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList)

}

func (route *AdminCategoriesRoutes) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
	url := "http://dblayer:8080/category/sub"
	resp, err := http.Get(url)
	if err != nil{
		fmt.Println(err)
	}
	defer resp.Body.Close()
	categoryList := []CategoryReturn{}
	err = json.NewDecoder(resp.Body).Decode(&categoryList)
	if err != nil {
		fmt.Println("Error decoding prime categories:", err)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList)
}

func (route *AdminCategoriesRoutes) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
	url := "http://dblayer:8080/category/final"
	resp, err := http.Get(url)
	if err != nil{
		fmt.Println(err)
	}
	categoryList := []CategoryReturn{}
	err = json.NewDecoder(resp.Body).Decode(&categoryList)
	if err != nil {
		fmt.Println("Error decoding prime categories:", err)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList)
}