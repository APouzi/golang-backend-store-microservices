package products

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/APouzi/DBLayer/helpers"
	"github.com/go-chi/chi/v5"
)

type ProductRoutesTray struct{
	DB *sql.DB
	getAllProductsStmt *sql.Stmt
	getOneProductStmt *sql.Stmt
	getOneVariationProductStment *sql.Stmt
	getProductPrimeCategoryByID *sql.Stmt
	getAllProductByCategoryStmt *sql.Stmt
	getAllProductByCategoryPrimeStmt *sql.Stmt
	getAllProductByCategorySubStmt *sql.Stmt
	getAllProductByCategoryFinalStmt *sql.Stmt
}

func GetProductRouteInstance(dbInst *sql.DB) *ProductRoutesTray{
	
	routeMap := prepareProductRoutes(dbInst)
	prd_tray := &ProductRoutesTray{
		getAllProductsStmt: routeMap["getAllProducts"],
		getOneProductStmt: routeMap["getOneProducts"],
		getOneVariationProductStment: routeMap["getOneVariationProducts"],
		getProductPrimeCategoryByID: routeMap["getProductPrimeCategoryByID"],
		getAllProductByCategoryStmt: routeMap["GetAllProductByCategoryStmt"],
		getAllProductByCategoryPrimeStmt: routeMap["GetAllProductByCategoryPrimeStmt"],
		getAllProductByCategorySubStmt: routeMap["GetAllProductByCategorySubStmt"],
		getAllProductByCategoryFinalStmt: routeMap["GetAllProductByCategoryFinalStmt"],
		DB: dbInst,
	}
	return prd_tray
}

func prepareProductRoutes(dbInst *sql.DB) map[string]*sql.Stmt{
	sqlStmentsMap := make(map[string]*sql.Stmt)
	
	getAllPrdStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, Product_Description FROM tblProducts")
	if err != nil {
		log.Fatal(err)
	}
	
	GetOneProductStmt, err := dbInst.Prepare("SELECT Product_ID, Product_Name, Product_Description, PRIMARY_IMAGE, Date_Created, Modified_Date FROM tblProducts WHERE Product_ID = ?")
	if err != nil {
		log.Fatal(err)
	}
	
	GetOneVariationStmt, err := dbInst.Prepare("SELECT Variation_ID, Product_ID, Variation_Name, Variation_Description, Variation_Price, PRIMARY_IMAGE FROM tblProductVariation WHERE Variation_ID = ?")
	if err != nil {
		log.Fatal(err)
	}

	getAllFinalStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, CategoryName FROM AllProductsInFinalView LIMIT 10 OFFSET ?")
	if err != nil {
		log.Fatal(err)
	}

	getAllSubStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, CategoryName FROM AllProductsInSubView LIMIT 10 OFFSET ?")
	if err != nil {
		log.Fatal(err)
	}

	getAllPrimeStment, err := dbInst.Prepare("SELECT Product_ID, Product_Name, CategoryName FROM AllProductsInPrimeView LIMIT 10 OFFSET ?")
	if err != nil {
		log.Fatal(err)
	}


	// GetAllProductsPrimeCategoryByID, err := dbInst.Prepare("SELECT tblProducts.Product_ID, tblProducts.Product_Name FROM tblProducts JOIN tblCatFinalProd ON tblCatFinalProd.Product_ID = tblProducts.Product_ID JOIN tblCategoriesFinal ON tblCategoriesFinal.Category_ID = tblCatFinalProd.CatFinalID JOIN tblCatSubFinal ON tblCatSubFinal.CatFinalID = tblCategoriesFinal.Category_ID JOIN tblCategoriesSub ON tblCategoriesSub.Category_ID = tblCatSubFinal.CatSubID JOIN tblCatPrimeSub ON tblCatPrimeSub.CatSubID = tblCategoriesSub.Category_ID JOIN tblCategoriesPrime ON tblCategoriesPrime.Category_ID = tblCatPrimeSub.CatPrimeID WHERE tblCategoriesPrime.Category_ID = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	
	sqlStmentsMap["getAllProducts"] = getAllPrdStment
	sqlStmentsMap["getOneProducts"] = GetOneProductStmt
	sqlStmentsMap["getOneVariationProducts"] = GetOneVariationStmt
	sqlStmentsMap["GetAllProductByCategoryFinalStmt"] = getAllFinalStment
	sqlStmentsMap["GetAllProductByCategorySubStmt"] = getAllSubStment
	sqlStmentsMap["GetAllProductByCategoryPrimeStmt"] = getAllPrimeStment
	// sqlStmentsMap["getProductPrimeCategoryByID"] = GetAllProductsPrimeCategoryByID
	
	return sqlStmentsMap
}

var allowed = map[string]string{
	"variation_id":"Variation_ID",
    "product_id": "Product_ID",
    "variation_name": "Variation_Name",
    "sku": "SKU",
    "upc": "UPC",
}

func (prdRoutes *ProductRoutesTray) GetOneProductVariationByParamEndPoint(w http.ResponseWriter, r *http.Request){
	q := r.URL.Query()
	var filters []string
	// var args []interface{}
	// argPos := 1 // for positional placeholders like $1
	for key, values := range q {
		key = strings.ToLower(key)
		fmt.Println("key and value",key,values)
		if _, ok := allowed[key]; !ok {
			helpers.ErrorJSON(w,errors.New("unknown paramater in Variation table"),http.StatusBadRequest)
			return // ignore unknown
		}
		// "key is the key and values is an array
		filters = append(filters, fmt.Sprintf("%s = %s", allowed[key], values[0]))
	}

	query := "SELECT Variation_ID, Product_ID, Variation_Name, Variation_Description, Variation_Price, SKU, UPC FROM tblProductVariation"
	query += " WHERE "

	if len(filters) > 0 {
		query += strings.Join(filters, " AND ")
	}
	rows, err := prdRoutes.DB.Query(query)
	if err != nil{
		fmt.Println("HANLDE THIS")
	}
	ListProducts := []ProductVariation{}
	prodJSON := ProductVariation{}
	for rows.Next(){
		
		err := rows.Scan(
			&prodJSON.VariationID,
			&prodJSON.ProductID,
			&prodJSON.VariationName,
			&prodJSON.VariationDescription,
			&prodJSON.VariationPrice,
			&prodJSON.SKU,
			&prodJSON.UPC,
		)

		if err != nil{
			fmt.Println("Scanning Error:",err)
		}
		ListProducts = append(ListProducts, prodJSON)
	}

	helpers.WriteJSON(w, http.StatusAccepted,ListProducts)

}
func (prdRoutes *ProductRoutesTray) GetAllProductsEndPoint(w http.ResponseWriter, r *http.Request) {

	rows, _ := prdRoutes.getAllProductsStmt.Query()
	ListProducts := []ProductJSONRetrieve{}
	prodJSON := ProductJSONRetrieve{}
	defer rows.Close()
	for rows.Next(){
		
		err := rows.Scan(
			&prodJSON.Product_ID,
			&prodJSON.Product_Name,
			&prodJSON.Product_Description,
		)

		if err != nil{
			fmt.Println("Scanning Error:",err)
		}
		ListProducts = append(ListProducts, prodJSON)
	}

	helpers.WriteJSON(w,200,ListProducts)

}

type ProductJSON struct {
	Product_ID          int     `json:"Product_ID"`
	Product_Name        string  `json:"Product_Name"`
	Product_Description string  `json:"Product_Description"`
	PRIMARY_IMAGE       string  `json:"PRIMARY_IMAGE,omitempty"`
	ProductDateAdded   string  `json:"DateAdded"`
	ModifiedDate       string `json:"ModifiedDate"`
}

func (prdRoutes *ProductRoutesTray) GetOneProductEndPoint(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hit getoneproductendpoint")
	productID, err :=  strconv.Atoi(chi.URLParam(r,"ProductID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	rows :=prdRoutes.getOneProductStmt.QueryRow(productID)
	prodJSON := ProductJSON{}
	
	err = rows.Scan(
		&prodJSON.Product_ID, 
		&prodJSON.Product_Name, 
		&prodJSON.Product_Description,  
		&prodJSON.PRIMARY_IMAGE,
		&prodJSON.ProductDateAdded,
		&prodJSON.ModifiedDate,
	)
	if err == sql.ErrNoRows{
		helpers.WriteJSON(w,404,"Cant find product son")
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,prodJSON)

}

type VariationRetrieve struct{
	Variation_ID int64 `json:"Variation_ID"`
	ProductID int64 `json:"Product_ID"`
	Name string `json:"Variation_Name"`
	Description string `json:"Variation_Description"`
	Price float32 `json:"Variation_Price"`
	PrimaryImage sql.NullString `json:"PRIMARY_IMAGE,omitempty"`

}

type FailedDBQuery struct {
	Msg string `json:"message"`
}

func (prdRoutes *ProductRoutesTray) GetOneVariationEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit getoneproductendpoint")
	productVariationID, err :=  strconv.Atoi(chi.URLParam(r,"VariationID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}
	fmt.Println("hit getoneproductendpoint", productVariationID)
	// row := route.DB.QueryRow("SELECT Variation_ID FROM tblProductVariation WHERE Variation_ID = ?",pil.VarID)
	rows := prdRoutes.getOneVariationProductStment.QueryRow(productVariationID)
	variationJSON := &VariationRetrieve{}
	
	err = rows.Scan(
		&variationJSON.Variation_ID,
		&variationJSON.ProductID, 
		&variationJSON.Name, 
		&variationJSON.Description,  
		&variationJSON.Price,
		&variationJSON.PrimaryImage,
	)
	if err == sql.ErrNoRows{
		fmt.Println("Doesnt work:",err)
		helpers.ErrorJSON(w, errors.New("could not get variation"),404)
		return
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,&variationJSON)

}

type VariationRet struct {
    Variation_ID  *int64   `json:"variation_id,omitempty"`
    ProductID    *int64   `json:"product_id,omitempty"`
    Name         *string  `json:"name,omitempty"`
    Description  *string  `json:"description,omitempty"`
    Price        *float64 `json:"price,omitempty"`
    PrimaryImage *string  `json:"primary_image,omitempty"`
}


func (prdRoutes *ProductRoutesTray) SearchProductsEndPoint(w http.ResponseWriter, r *http.Request){
	query := r.URL.Query().Get("q")
	if query == ""{
		fmt.Println("no search term")
		helpers.ErrorJSON(w, errors.New("no search term"),404)
		return
	}
	fmt.Println(r.URL.Query().Get("q"))
	fmt.Println("query from product end point", query)
	sqlQuery := `
		SELECT p.Product_ID, p.Product_Name, p.Product_Description,
			   pv.Variation_ID, pv.Variation_Name, pv.Variation_Description, pv.Variation_Price,
			   pil.Inv_ID, pil.Quantity, pil.Location_At
		FROM tblProducts p
		LEFT JOIN tblProductVariation pv ON p.Product_ID = pv.Product_ID
		LEFT JOIN tblProductInventoryLocation pil ON pv.Variation_ID = pil.Variation_ID
		WHERE p.Product_Name LIKE ? OR pv.Variation_Name LIKE ?
	`

	//the question marks in the sql statement are replaced by the values of the query string, specifically the two parameters after the query string.
	rows, err := prdRoutes.DB.Query(sqlQuery, "%"+query+"%", "%"+query+"%")
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var product_list []map[string]interface{}
	for rows.Next() {
		var product Product
		var variation VariationRet
		var inventory struct {
			Inv_ID     *int64  `json:"inv_id,omitempty"`
			Quantity  *int64  `json:"quantity,omitempty"`
			LocationAt  *string `json:"location,omitempty"`
		}

		

		err := rows.Scan(
			&product.Product_ID, &product.Product_Name, &product.Product_Description,
			&variation.Variation_ID, &variation.Name, &variation.Description, &variation.Price,
			&inventory.Inv_ID, &inventory.Quantity, &inventory.LocationAt,
		)
		if err != nil {
			helpers.ErrorJSON(w, fmt.Errorf("error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		

		result := map[string]interface{}{
			"product": product,
		}

		result["variation"] = variation


		result["inventory"] = inventory

		product_list = append(product_list, result)
	}

	if err = rows.Err(); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("error iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, product_list)
}


func (prdRoutes *ProductRoutesTray) GetProductAndVariationsByProductID(w http.ResponseWriter, r *http.Request){
	productID := chi.URLParam(r,"productID")
	if productID == ""{
		fmt.Println("no search term")
		helpers.ErrorJSON(w, errors.New("no search term"),404)
		return
	}
	fmt.Println(r.URL.Query().Get("q"))
	fmt.Println("query from product end point", productID)
	sqlQuery := `
		SELECT p.Product_ID, p.Product_Name, p.Product_Description,
			    pv.Variation_ID,pv.Variation_Name, pv.Variation_Description, pv.Variation_Price,
			   pil.Inv_ID, pil.Quantity, pil.Location_At
		FROM tblProducts p
		LEFT JOIN tblProductVariation pv ON p.Product_ID = pv.Product_ID
		LEFT JOIN tblProductInventoryLocation pil ON pv.Variation_ID = pil.Variation_ID
		WHERE p.Product_ID = ?
	`

	//the question marks in the sql statement are replaced by the values of the query string, specifically the two parameters after the query string.
	rows, err := prdRoutes.DB.Query(sqlQuery, productID)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var product_list []map[string]interface{}
	for rows.Next() {
		var product Product
		var variation VariationRet
		var inventory struct {
			Inv_ID     *int64
			Quantity  *int64
			LocationAt *string
		}

		

		err := rows.Scan(
			&product.Product_ID, &product.Product_Name, &product.Product_Description,
			&variation.Variation_ID, &variation.Name, &variation.Description, &variation.Price,
			&inventory.Inv_ID, &inventory.Quantity, &inventory.LocationAt,
		)
		if err != nil {
			helpers.ErrorJSON(w, fmt.Errorf("error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		

		result := map[string]interface{}{
			"product": product,
		}
			result["variation"] = variation


			result["inventory"] = inventory

		product_list = append(product_list, result)
	}

	if err = rows.Err(); err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("error iterating rows: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, product_list)
}


func (prdRoutes *ProductRoutesTray) GetProductAndVariationsPaginated(w http.ResponseWriter, r *http.Request) {
    // Parse pagination params from query
    pageStr := r.URL.Query().Get("page")
    pageSizeStr := r.URL.Query().Get("page_size")


	page, err := strconv.Atoi(pageStr)

	pageSize, err := strconv.Atoi(pageSizeStr)

    if page < 1 {
        page = 1
    }
    if pageSize < 1 {
        pageSize = 10
    }

    offset := (page - 1) * pageSize

    // Count total records
    var totalCount int
    countQuery := `SELECT COUNT(*) FROM tblProducts p
                   LEFT JOIN tblProductVariation pv ON p.Product_ID = pv.Product_ID
                	`
    err = prdRoutes.DB.QueryRow(countQuery).Scan(&totalCount)
    if err != nil {
        helpers.ErrorJSON(w, fmt.Errorf("failed to get count: %v", err), http.StatusInternalServerError)
        return
    }

    totalPages := (totalCount + pageSize - 1) / pageSize
	fmt.Println("total pages",totalPages)
    // Main query with LIMIT/OFFSET
    sqlQuery := `
        SELECT p.Product_ID, p.Product_Name, p.Product_Description,
               pv.Variation_ID, pv.Variation_Name, pv.Variation_Description, ps.Size_ID,
               ps.Size_Name, ps.Size_Description, ps.Variation_ID, ps.Variation_Price, ps.SKU, ps.UPC, ps.PRIMARY_IMAGE
        FROM tblProducts p
        LEFT JOIN tblProductVariation pv ON p.Product_ID = pv.Product_ID
		LEFT JOIN tblProductSize ps ON ps.Variation_ID = pv.Variation_ID 
        LIMIT ? OFFSET ?
    `

    rows, err := prdRoutes.DB.Query(sqlQuery, pageSize, offset)
    if err != nil {
        helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
        return
    }
    defer rows.Close()
	// var products []ProductVariationPaginated = []ProductVariationPaginated{}
    var productsMapForJoins map[int64]ProductPaginated = map[int64]ProductPaginated{}
	dataSendBack := []ProductPaginated{}
    for rows.Next() {
        var product rowData

        err := rows.Scan(
            &product.ProductID, &product.ProductName, &product.ProductDescription,
            &product.VariationID, &product.VariationName, &product.VariationDesc,
            &product.SizeID,&product.SizeName, &product.SizeDesc, &product.VariationID, &product.SizeVariationPrice,&product.SKU, &product.UPC,&product.PrimaryImage,
        )
		if _,ok := productsMapForJoins[product.ProductID]; !ok{
			productsMapForJoins[product.ProductID] =  ProductPaginated{
				ProductID: product.ProductID,
				ProductName: product.ProductName,
				ProductDescription: product.ProductDescription,
				Variations: map[int64]Variation{},
			}
		}
		
		if product.VariationID == nil || *product.VariationID == 0{
				continue
		}
		if _, ok:= productsMapForJoins[product.ProductID].Variations[*product.VariationID]; !ok {
		productsMapForJoins[product.ProductID].Variations[*product.VariationID] = Variation{
				Name: product.VariationName,
				Description: product.VariationDesc,
				ProductSize: map[int64]ProductSize{},
		}
			
		}
		if product.SizeID == nil || *product.SizeID == 0 {
			fmt.Println("product size doesn't exist doesn't exist, continue",product.VariationID)
			continue
		}
		if _,ok := productsMapForJoins[product.ProductID].Variations[*product.VariationID].ProductSize[*product.SizeID]; !ok{
			productsMapForJoins[product.ProductID].Variations[*product.VariationID].ProductSize[*product.SizeID] = ProductSize{
				SizeID:product.SizeID,
				SizeName:product.SizeName,
				SizeDescription:product.SizeDesc,
				VariationPrice: &product.SizeVariationPrice.Float64, 
				SKU: product.SKU, 
				UPC: product.UPC, 
			}
		}
	
		if err != nil {
			helpers.ErrorJSON(w, fmt.Errorf("error scanning row: %v", err), http.StatusInternalServerError)
			return
		}

		if err = rows.Err(); err != nil {
			helpers.ErrorJSON(w, fmt.Errorf("error iterating rows: %v", err), http.StatusInternalServerError)
			return
		}

    
	}
	for _, v := range productsMapForJoins {
		dataSendBack = append(dataSendBack, v)
	}
	sort.Slice(dataSendBack, func(i, j int) bool {
		return dataSendBack[i].ProductID  < dataSendBack[j].ProductID
	})
	response := PaginatedResponse{
        Data:       dataSendBack,
        Page:       page,
        PageSize:   pageSize,
        TotalCount: totalCount,
        TotalPages: totalPages,
    }
	helpers.WriteJSON(w, http.StatusOK, response)

}




type CategoryRetrieval struct{
	Product_ID sql.NullInt64 `json:"product_id"`
	Product_Name sql.NullString `json:"product_name"`
	CategoryName sql.NullString `json:"category_name"`
}

type CategorySendOff struct{
	Product_ID int64 `json:"product_id"`
	Product_Name string `json:"product_name"`
	CategoryName string `json:"category_name"`
}


func (prdRoutes *ProductRoutesTray) GetAllProductsInFinalCategoryViewEndPoint(w http.ResponseWriter, r *http.Request) {

	CatDBRet := &CategoryRetrieval{}
	
	category_name := r.URL.Query().Get("final_category_name")
	page, err :=  strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("getting page failed: %v", err), http.StatusInternalServerError)
		return
	}

	page = 10 * (page - 1)

	res, err := prdRoutes.getAllProductByCategoryFinalStmt.Query(page)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	SendOffList := []CategorySendOff{}
	CatDBSend := &CategorySendOff{}
	for res.Next(){
		
		err := res.Scan(
			&CatDBRet.Product_ID,
			&CatDBRet.Product_Name,
			&CatDBRet.CategoryName,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		CatDBSend.Product_ID = CatDBRet.Product_ID.Int64
		CatDBSend.Product_Name = CatDBRet.Product_Name.String
		CatDBSend.CategoryName = CatDBRet.CategoryName.String

		SendOffList = append(SendOffList, *CatDBSend)
	}

	fmt.Println("Category Name:", category_name)
	
	helpers.WriteJSON(w,http.StatusAccepted,SendOffList)
}

func (prdRoutes *ProductRoutesTray) GetAllProductsInSubCategoryViewEndPoint(w http.ResponseWriter, r *http.Request) {

	CatDBRet := &CategoryRetrieval{}
	
	category_name := r.URL.Query().Get("sub_category_name")
	page, err :=  strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("getting page failed: %v", err), http.StatusInternalServerError)
		return
	}

	page = 10 * (page - 1)

	res, err := prdRoutes.getAllProductByCategorySubStmt.Query(page)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	SendOffList := []CategorySendOff{}
	CatDBSend := &CategorySendOff{}
	for res.Next(){
		
		err := res.Scan(
			&CatDBRet.Product_ID,
			&CatDBRet.Product_Name,
			&CatDBRet.CategoryName,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		CatDBSend.Product_ID = CatDBRet.Product_ID.Int64
		CatDBSend.Product_Name = CatDBRet.Product_Name.String
		CatDBSend.CategoryName = CatDBRet.CategoryName.String

		SendOffList = append(SendOffList, *CatDBSend)
	}

	fmt.Println("Category Name:", category_name)
	
	helpers.WriteJSON(w,http.StatusAccepted,SendOffList)
}

func (prdRoutes *ProductRoutesTray) GetAllProductsInPrimeCategoryViewEndPoint(w http.ResponseWriter, r *http.Request) {

	CatDBRet := &CategoryRetrieval{}
	
	category_name := r.URL.Query().Get("prime_category_name")
	page, err :=  strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("getting page failed: %v", err), http.StatusInternalServerError)
		return
	}

	page = 10 * (page - 1)

	res, err := prdRoutes.getAllProductByCategoryPrimeStmt.Query(page)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	SendOffList := []CategorySendOff{}
	CatDBSend := &CategorySendOff{}
	for res.Next(){
		
		err := res.Scan(
			&CatDBRet.Product_ID,
			&CatDBRet.Product_Name,
			&CatDBRet.CategoryName,
		)
		if err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		CatDBSend.Product_ID = CatDBRet.Product_ID.Int64
		CatDBSend.Product_Name = CatDBRet.Product_Name.String
		CatDBSend.CategoryName = CatDBRet.CategoryName.String

		SendOffList = append(SendOffList, *CatDBSend)
	}

	fmt.Println("Category Name:", category_name)
	
	helpers.WriteJSON(w,http.StatusAccepted,SendOffList)
}



func (prdRoutes *ProductRoutesTray) GetOneProductSizeEndPoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hit GetOneProductSizeEndPoint")
	productSizeID, err :=  strconv.Atoi(chi.URLParam(r,"SizeID"))
	if err != nil{
		fmt.Println("String to Int failed:", err)
	}

	row := prdRoutes.DB.QueryRow("SELECT Size_ID, Size_Name, Size_Description, Variation_ID, Variation_Price, SKU, UPC, Date_Created, Modified_Date FROM tblProductSize WHERE Size_ID = ?",productSizeID)

	// row := prdRoutes.getOneVariationProductStment.QueryRow(productVariationID)
	variationJSON := &ProductSize{}
	
	err = row.Scan(
		&variationJSON.SizeID,
		&variationJSON.SizeName, 
		&variationJSON.SizeDescription, 
		&variationJSON.VariationID,  
		&variationJSON.VariationPrice,
		&variationJSON.SKU,
		&variationJSON.UPC,
		&variationJSON.DateCreated,
		&variationJSON.ModifiedDate,

	)
	if err == sql.ErrNoRows{
		fmt.Println("Doesnt work:",err)
		helpers.ErrorJSON(w, errors.New("could not get variation"),404)
		return
	}
	if err != nil{
		fmt.Println("scanning error:",err)
	}
	
	helpers.WriteJSON(w,200,&variationJSON)

}


func (prdRoutes *ProductRoutesTray) GetAllProductTaxCodeEndPoint(w http.ResponseWriter, r *http.Request) {
	var taxCodes []ProductTaxCode
	rows, err := prdRoutes.DB.Query("SELECT TaxCode_ID, TaxCode_Name, TaxCode_Description, TaxCode FROM tblProductTaxCode")
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var taxCode ProductTaxCode
		if err := rows.Scan(&taxCode.TaxCodeID, &taxCode.TaxCodeName, &taxCode.TaxCodeDescription, &taxCode.TaxCode); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		taxCodes = append(taxCodes, taxCode)
	}

	helpers.WriteJSON(w, http.StatusOK, taxCodes)
}

func(prdRoutes *ProductRoutesTray) GetOneProductTaxCodeEndpoint(w http.ResponseWriter, r *http.Request) {
	taxCodeID, err := strconv.Atoi(chi.URLParam(r, "TaxCodeID"))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid tax code ID: %v", err), http.StatusBadRequest)
		return
	}

	row := prdRoutes.DB.QueryRow("SELECT TaxCode_ID, TaxCode_Name, TaxCode_Description, TaxCode FROM tblProductTaxCode WHERE TaxCode_ID = ?", taxCodeID)

	var taxCode ProductTaxCode
	err = row.Scan(&taxCode.TaxCodeID, &taxCode.TaxCodeName, &taxCode.TaxCodeDescription, &taxCode.TaxCode)
	if err == sql.ErrNoRows {
		helpers.ErrorJSON(w, errors.New("tax code not found"), http.StatusNotFound)
		return
	}
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}

	helpers.WriteJSON(w, http.StatusOK, taxCode)  
}

func (prdRoutes *ProductRoutesTray) GetAllProductTaxCodeEndPointFromProductSizeIntermediary(w http.ResponseWriter, r *http.Request) {
	productSizeID, err := strconv.Atoi(chi.URLParam(r, "SizeID"))
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("invalid product size ID: %v", err), http.StatusBadRequest)
		return
	}

	rows, err := prdRoutes.DB.Query("SELECT t.TaxCode_ID, t.TaxCode_Name, t.TaxCode_Description, t.TaxCode FROM tblProductTaxCode t JOIN tblProductSizeTaxCode pstc ON t.TaxCode_ID = pstc.TaxCode_ID WHERE pstc.Size_ID = ?", productSizeID)
	if err != nil {
		helpers.ErrorJSON(w, fmt.Errorf("database query failed: %v", err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var taxCodes []ProductTaxCode
	for rows.Next() {
		var taxCode ProductTaxCode
		if err := rows.Scan(&taxCode.TaxCodeID, &taxCode.TaxCodeName, &taxCode.TaxCodeDescription, &taxCode.TaxCode); err != nil {
			fmt.Println("Error scanning row:", err)
			continue
		}
		taxCodes = append(taxCodes, taxCode)
	}

	helpers.WriteJSON(w, http.StatusOK, taxCodes)
}
