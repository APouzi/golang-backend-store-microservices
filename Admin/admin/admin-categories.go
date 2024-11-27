package adminendpoints

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Apouzi/Golang-Admin-Service/helpers"
)

type AdminCategoriesRoutes struct{
	DB *sql.DB
}

func InstanceAdminCategoriesRoutes(  ) *AdminCategoriesRoutes {
	r := &AdminCategoriesRoutes{
	}
	return r
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
		fmt.Println(err)
	}


	url := "http://dblayer:8080/category/sub"
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
	}


	url := "http://dblayer:8080/category/primetosub"
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
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesPrime"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)

}

func (route *AdminCategoriesRoutes) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesSub"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}

func (route *AdminCategoriesRoutes) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT CategoryName, CategoryDescription FROM tblCategoriesFinal"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategoryReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategoryReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}