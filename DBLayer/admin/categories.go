package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)
func GetCategoriesRouteInstance (dbInstance *sql.DB) *CategoriesRoutesTray {
	route := &CategoriesRoutesTray{
		DB: dbInstance,
	}
	return route
}

type CategoriesRoutesTray struct {
	DB *sql.DB
}

func (adminProdRoutes *CategoriesRoutesTray) CreatePrimeCategory(w http.ResponseWriter, r *http.Request) {
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("in dblayer: CreatePrimeCategory", category_read)
	result, err := adminProdRoutes.DB.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription)
	if err != nil {
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w, 200, catReturn)
	if err != nil {
		fmt.Println("There was an error trying to send data back", err)
	}

}

func (route *CategoriesRoutesTray) CreateSubCategory(w http.ResponseWriter, r *http.Request) {
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil {
		fmt.Println(err)
	}

	result, err := route.DB.Exec("INSERT INTO tblCategoriesSub(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription)
	if err != nil {
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w, 200, catReturn)
	if err != nil {
		fmt.Println("There was an error trying to send data back", err)
	}
}

func (route *CategoriesRoutesTray) CreateFinalCategory(w http.ResponseWriter, r *http.Request) {
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil {
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCategoriesFinal(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription)
	if err != nil {
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}

	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w, 200, catReturn)
	if err != nil {
		fmt.Println("There was an error trying to send data back", err)
	}
}

func (route *CategoriesRoutesTray) ConnectPrimeToSubCategory(w http.ResponseWriter, r *http.Request) {
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w, r, &FinalSub)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed reading json"), 500)
		return
	}
	fmt.Println("FinalSub in ConnectPrimeToSubCategory:", FinalSub)
	result, err := route.DB.Exec("INSERT INTO tblCatPrimeSub(CatPrimeID,  CatSubID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed inserting into tblCatPrimeSub"), 500)
		return
	}
	_, err = result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed getting last insert id"), 500)
		return
	}

	helpers.WriteJSON(w, http.StatusAccepted, FinalSub)
}

func (route *CategoriesRoutesTray) ConnectSubToFinalCategory(w http.ResponseWriter, r *http.Request) {
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalSub := CatToCat{}
	err := helpers.ReadJSON(w, r, &FinalSub)
	if err != nil {
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatSubFinal(CatSubID, CatFinalID) VALUES(?,?)", FinalSub.CatStart, FinalSub.CatEnd)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed inserting into tblCatSubFinal"), 500)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed getting last insert id"), 500)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, FinalSub)
}

func (route *CategoriesRoutesTray) ConnectFinalToProdCategory(w http.ResponseWriter, r *http.Request) {
	// Frontend will have the names and ids, so I PROBABLY wont need to do a search regarding the names of category to get ids
	FinalProd := CatToProd{}
	err := helpers.ReadJSON(w, r, &FinalProd)
	if err != nil {
		fmt.Println(err)
	}
	result, err := route.DB.Exec("INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)", FinalProd.Cat, FinalProd.Prod)

	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed inserting into tblCatFinalProd"), 500)
		return
	}

	_, err = result.LastInsertId()
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("failed getting last insert id"), 500)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, FinalProd)
}

func (route *CategoriesRoutesTray) InsertIntoFinalProd(w http.ResponseWriter, r *http.Request) {
	ReadCatR := ReadCat{}
	err := helpers.ReadJSON(w, r, &ReadCatR)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("InsertIntoCategory ReadCatR", ReadCatR)
	FinalProd := "INSERT INTO tblCatFinalProd(CatFinalID, Product_ID) VALUES(?,?)"
	route.DB.Exec(FinalProd, 1, ReadCatR.Category)
}

func (route *CategoriesRoutesTray) DeletePrimeCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatPrimeName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	_, err := route.DB.Exec("DELETE FROM tblCategoriesPrime WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *CategoriesRoutesTray) DeleteSubCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatSubName")
	if CatName == ""{
		fmt.Println("No CatSubName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesSub WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatSubName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}
	
	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}


func (route *CategoriesRoutesTray) DeleteFinalCategory(w http.ResponseWriter, r *http.Request){
	CatName := chi.URLParam(r,"CatFinalName")
	if CatName == ""{
		fmt.Println("No CatPrimeName wasn't pulled")
		return
	}
	
	_, err := route.DB.Exec("DELETE FROM tblCategoriesFinal WHERE CategoryName = ?", CatName)
	if err != nil{
		fmt.Println("Failed deletion in CatPrimeName")
		helpers.ErrorJSON(w, errors.New("failed deletion in table"), 500)
		return
	}

	sendBack := DeletedSendBack{SendBack:false}
	helpers.WriteJSON(w,200,sendBack)
}

func (route *CategoriesRoutesTray) EditPrimeCategory(w http.ResponseWriter, r *http.Request){
	editPC := CategoryEdit{}
	err := helpers.ReadJSON(w,r,&editPC)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	if editPC.CategoryName == nil && editPC.CategoryDescription == nil {
		helpers.ErrorJSON(w, errors.New("missing fields in request, nothing to update"), http.StatusBadRequest)
		return
	}
	var sb strings.Builder 
	sb.WriteString("UPDATE tblCategoriesPrime SET ")
	if editPC.CategoryName != nil {
		sb.WriteString("CategoryName = ?, ")
	}
	if editPC.CategoryDescription != nil {
		sb.WriteString("CategoryDescription = ? ")
	}
	sb.WriteString("WHERE Category_ID = ?")
	var query string = sb.String()
	_, err = route.DB.Exec(query,
		editPC.CategoryName, editPC.CategoryDescription, editPC.CategoryId)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadGateway)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, editPC)
}

func (route *CategoriesRoutesTray) EditSubCategory(w http.ResponseWriter, r *http.Request){
	editSC := CategoryEdit{}
	err := helpers.ReadJSON(w,r,&editSC)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	if editSC.CategoryName == nil && editSC.CategoryDescription == nil {
		helpers.ErrorJSON(w, errors.New("missing fields in request, nothing to update"), http.StatusBadRequest)
		return
	}
	var sb strings.Builder
	sb.WriteString("UPDATE tblCategoriesSub SET ")
	if editSC.CategoryName != nil {
		sb.WriteString("CategoryName = ?, ")
	}
	if editSC.CategoryDescription != nil {
		sb.WriteString("CategoryDescription = ? ")
	}
	sb.WriteString("WHERE Category_ID = ?")
	var query string = sb.String()
	_, err = route.DB.Exec(query,
		editSC.CategoryName, editSC.CategoryDescription, editSC.CategoryId)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadGateway)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, editSC)
}
func (route *CategoriesRoutesTray) EditFinalCategory(w http.ResponseWriter, r *http.Request){
	editFC := CategoryEdit{}
	err := helpers.ReadJSON(w,r,&editFC)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadRequest)
		return
	}
	if editFC.CategoryName == nil && editFC.CategoryDescription == nil {
		helpers.ErrorJSON(w, errors.New("missing fields in request, nothing to update"), http.StatusBadRequest)
		return
	}
	var sb strings.Builder 
	sb.WriteString("UPDATE tblCategoriesFinal SET ")
	if editFC.CategoryName != nil {
		sb.WriteString("CategoryName = ?, ")
	}
	if editFC.CategoryDescription != nil {
		sb.WriteString("CategoryDescription = ? ")
	}
	sb.WriteString("WHERE Category_ID = ?")
	var query string = sb.String()
	_, err = route.DB.Exec(query,
		editFC.CategoryName, editFC.CategoryDescription, editFC.CategoryId)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w, err, http.StatusBadGateway)
		return
	}
	helpers.WriteJSON(w, http.StatusAccepted, editFC)
}