package admin

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"

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

func (prdRoutes *ProductRoutesTray) CreateProductMultiChain(w http.ResponseWriter, r *http.Request) {
	transaction, err := prdRoutes.DB.Begin()
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




func (prdRoutes *ProductRoutesTray) CreateProductVariation(w http.ResponseWriter, r *http.Request) {

	ProductID := chi.URLParam(r, "ProductID")
	variation := VariationCreate{}
	helpers.ReadJSON(w,r, &variation)
	var prodID int64
	err := prdRoutes.DB.QueryRow("SELECT Product_ID FROM tblProducts WHERE Product_ID = ?",ProductID).Scan(&prodID)
	if err == sql.ErrNoRows{
		msg := ProdExist{}
		msg.ProductExists = false
		msg.Message = "Product provided does not exist"
		helpers.WriteJSON(w,http.StatusAccepted,msg)
		log.Println("Variation Creation failed, Product doesn't exist")
		return
	}
	// Implement the returns for this to allow for proper exiting 

	var varit sql.Result
	if variation.PrimaryImage != "" {
		varitCrt := variCrtd{}
		varit, err = prdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price) VALUES(?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price)
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
	varit, err = prdRoutes.DB.Exec("INSERT INTO tblProductVariation(Product_ID, Variation_Name, Variation_Description, Variation_Price, PRIMARY_IMAGE) VALUES(?,?,?,?,?)", ProductID,variation.Name, variation.Description, variation.Price, variation.PrimaryImage)
	if err != nil{
		fmt.Println("insert into tblProductVariation failed")
		fmt.Println(err)
		helpers.ErrorJSON(w, errors.New("insert into tblProductVariation failed, could not retrieve varitation id"),400)
	}
}


type ProdInvLocCreation struct{
	VarID int64 `json:"Variation_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}

type PILCreated struct{
	InvID int64 `json:"Inv_ID"`
	Quantity int `json:"Quantity"`
	Location string `json:"Location"`
}

func (prdRoutes *ProductRoutesTray) CreateInventoryLocation(w http.ResponseWriter, r *http.Request) {

	pil := ProdInvLocCreation{}
	helpers.ReadJSON(w,r,&pil)

	res ,err:= prdRoutes.DB.Exec("INSERT INTO tblProductInventoryLocation(Variation_ID, Quantity, Location_At) VALUES(?,?,?)", pil.VarID,pil.Quantity,pil.Location)
	
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
	pilReturn := PILCreated{}
	pilReturn.InvID = pilID
	pilReturn.Quantity = pil.Quantity
	pilReturn.Location = pil.Location
	helpers.WriteJSON(w, http.StatusAccepted, pil)
}
