package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

type ProductRoutesTray struct{
	DB *sql.DB
}

func GetProductRouteInstance(dbInst *sql.DB) *ProductRoutesTray{
	
	// routeMap := prepareProductRoutes(dbInst)
	return &ProductRoutesTray{
		// getAllProductsStmt: routeMap["getAllProducts"],
		DB: dbInst,
	}
}

func (adminProdRoutes *ProductRoutesTray) CreateProductMultiChain(w http.ResponseWriter, r *http.Request) {
	transaction, err := adminProdRoutes.DB.Begin()
	if err != nil {
		log.Println("Error creating a transation in CreateProduct")
		log.Println(err)
	}

	productRetrieve := &ProductCreate{}

	helpers.ReadJSON(w, r, &productRetrieve)
	fmt.Println("product retrieve at db:",productRetrieve)
	tRes, err := transaction.Exec("INSERT INTO tblProducts(Product_Name, Product_Description) VALUES(?,?)", productRetrieve.Name, productRetrieve.Description)
	if err != nil {
		fmt.Println("transaction at tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
	}
	prodID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProduct has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	tRes, err = transaction.Exec("INSERT INTO tblProductVariation(Product_ID,Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", prodID, productRetrieve.VariationName, productRetrieve.VariationDescription, productRetrieve.VariationPrice)
	if err != nil {
		fmt.Println("transaction at tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}

	ProdVarID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println("retrieval of LastInsertID of tblProductVariation has failed")
		fmt.Println(err)
		transaction.Rollback()
		return
	}
	PCR := ProductCreateRetrieve{
		ProductID: prodID,
		VarID:     ProdVarID,
	}
	if productRetrieve.LocationAt == "" {

		err = transaction.Commit()
		if err != nil {
			fmt.Println(err)
			transaction.Rollback()
			return
		}
		helpers.WriteJSON(w, http.StatusAccepted, &PCR)
		return
	}

	tRes, err = transaction.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_AT) VALUES(?,?,?)", ProdVarID, productRetrieve.VariationQuantity, productRetrieve.LocationAt)
	if err != nil {
		fmt.Println("transaction at tblProductInventory has failed")
		fmt.Println(err)
	}
	invID, err := tRes.LastInsertId()
	if err != nil {
		fmt.Println(err)
	}
	PCR.ProdInvLoc = invID
	err = transaction.Commit()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("print of prod created",PCR)
	helpers.WriteJSON(w,http.StatusAccepted,PCR)
}




func (adminProdRoutes *ProductRoutesTray) CreateProductVariation(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	fmt.Println("productID", ProductID)
	variation := VariationCreate{PrimaryImage: ""}
	helpers.ReadJSON(w,r, &variation)
	// var prodID int64
	varitCrt := variCrtd{}
	// var varit sql.Result
	if variation.PrimaryImage != "" {
		
		varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price)
		if err != nil{
			log.Println("insert into tblProductVariation failed")
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed"),400)
			return
		}
		varitCrt.VariationID, err = varit.LastInsertId()
		if err != nil{
			log.Println(err)
			helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
			return
		}
		helpers.WriteJSON(w, http.StatusCreated,varitCrt)
	}
	prodid, err := strconv.Atoi(ProductID)
	varit, err := adminProdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", prodid,variation.Name, variation.Description, variation.Price)
	if err != nil{
		fmt.Println("insert into tblProductVariation failed")
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
	}
	varitCrt.VariationID, err = varit.LastInsertId()
	
	helpers.WriteJSON(w, http.StatusCreated,varitCrt)
}




func (adminProdRoutes *ProductRoutesTray) CreateInventoryLocation(w http.ResponseWriter, r *http.Request) {

	prodInvLoc := ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&prodInvLoc)
	
	res ,err:= adminProdRoutes.DB.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_At) VALUES(?,?,?)", prodInvLoc.VarID,prodInvLoc.Quantity,prodInvLoc.Location)
	
	if err != nil{
		fmt.Println("failed to create tblProductInventoryLocation")
		fmt.Println(err)
		helpers.ErrorJSON(w,err,http.StatusForbidden)
		return
	}

	pilID, err := res.LastInsertId()
	
	if err != nil{
		fmt.Println("result of tblProductInventoryLocation failed")
	}
	pilReturn := PILCreated{InvID:pilID, Quantity: prodInvLoc.Quantity, Location: prodInvLoc.Location }
	helpers.WriteJSON(w, http.StatusAccepted, pilReturn)
}


func (adminProdRoutes *ProductRoutesTray) CreatePrimeCategory(w http.ResponseWriter, r *http.Request) {
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	fmt.Println("in dblayer: CreatePrimeCategory",category_read)
	result, err := adminProdRoutes.DB.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}

	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w,200,catReturn)
	if err != nil{
		fmt.Println("There was an error trying to send data back",err)
	}

}

func (route *ProductRoutesTray) CreateSubCategory(w http.ResponseWriter, r *http.Request){
	category_read := CategoryInsert{}
	err := helpers.ReadJSON(w, r, &category_read)
	if err != nil{
		fmt.Println(err)
	}

	
	result, err := route.DB.Exec("INSERT INTO tblCategoriesSub(CategoryName, CategoryDescription) VALUES(?,?)", category_read.CategoryName, category_read.CategoryDescription )
	if err != nil{
		fmt.Println(err)
	}
	resultID, err := result.LastInsertId()
	if err != nil{
		fmt.Println(err)
	}
	
	catReturn := CategoryReturn{CategoryId: resultID, CategoryName: category_read.CategoryName, CategoryDescription: category_read.CategoryDescription}
	err = helpers.WriteJSON(w,200,catReturn)
	if err != nil{
		fmt.Println("There was an error trying to send data back",err)
	}
}
