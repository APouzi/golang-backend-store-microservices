package products

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/APouzi/DBLayer/helpers"
)






type CategoriesRoutesTray struct{
	Db *sql.DB
}






func (route *ProductRoutesTray) ReturnAllPrimeCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT Category_ID,CategoryName, CategoryDescription FROM tblCategoriesPrime"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategorySingleReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategorySingleReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryId, &category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)

}

func (route *ProductRoutesTray) ReturnAllSubCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT Category_ID, CategoryName, CategoryDescription FROM tblCategoriesSub"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategorySingleReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategorySingleReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryId, &category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}

func (route *ProductRoutesTray) ReturnAllFinalCategories(w http.ResponseWriter, r *http.Request){
	query := "SELECT Category_ID, CategoryName, CategoryDescription FROM tblCategoriesFinal"
	rows,err := route.DB.Query(query)
	if err != nil{
		fmt.Println(err)
	}
	category := CategorySingleReturn{}
	categoryList := CategoriesList{}
	categoryList.collection = []CategorySingleReturn{}
	for rows.Next(){
		rows.Scan(&category.CategoryId, &category.CategoryName, &category.CategoryDescription)
		categoryList.collection = append(categoryList.collection, category)
	}
	helpers.WriteJSON(w,http.StatusAccepted, categoryList.collection)
}


func (route *ProductRoutesTray) ReturnCategoryTree(w http.ResponseWriter, r *http.Request){
	query := `SELECT prime.Category_ID, prime.CategoryName, prime.CategoryDescription, sub.Category_ID, sub.CategoryName, sub.CategoryDescription, final.Category_ID, final.CategoryName, final.CategoryDescription
	FROM tblCategoriesPrime AS prime
	JOIN tblCatPrimeSub AS ps ON prime.Category_ID = ps.CatPrimeID
	JOIN tblCategoriesSub AS sub ON ps.CatSubID = sub.Category_ID
	JOIN tblCatSubFinal AS sf ON sub.Category_ID = sf.CatSubID
	JOIN tblCategoriesFinal AS final ON sf.CatFinalID = final.Category_ID`
	rows, err := route.DB.Query(query)
	if err != nil {
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusBadGateway)
	}
	categoryTree := CategoryTree{}
	categoryTree.Categories = map[string]PrimeCategoryTree{}
	for rows.Next() {
		var primeName, subName, finalName, primeDesc, subDesc, finalDesc string
		var primeID, subID, finalID int
		rows.Scan(&primeID, &primeName, &primeDesc, &subID, &subName, &subDesc, &finalID, &finalName, &finalDesc)
		// Build the category tree structure here
		if _, exists := categoryTree.Categories[primeName]; !exists {
			categoryTree.Categories[primeName] = PrimeCategoryTree{
				PrimeCategoryID:   primeID,
				PrimeCategoryName: primeName,
				Categories: map[string]SubCategoryTree{
					subName: {
						SubCategoryID:   subID,
						SubCategoryName: subName,
						Categories: map[string]FinalCategoryTree{
							finalName: {
								FinalCategoryID:   finalID,
								FinalCategoryName: finalName,
							},
						},
					},
				},
			}
		} else {
			primeCategory := categoryTree.Categories[primeName]
			if _, subExists := primeCategory.Categories[subName]; !subExists {
				primeCategory.Categories[subName] = SubCategoryTree{
					SubCategoryID:   subID,
					SubCategoryName: subName,
					Categories: map[string]FinalCategoryTree{
						finalName: {
							FinalCategoryID:   finalID,
							FinalCategoryName: finalName,
						},
					},
				}
			}else{
				subCategory := primeCategory.Categories[subName]
				if _, finalExists := subCategory.Categories[finalName]; !finalExists{
					subCategory.Categories[finalName] = FinalCategoryTree{
						FinalCategoryID:   finalID,
						FinalCategoryName: finalName,
					}
				}
				primeCategory.Categories[subName] = subCategory
				
			}
			categoryTree.Categories[primeName] = primeCategory
		}
		
	}
	helpers.WriteJSON(w, http.StatusAccepted, categoryTree)
}