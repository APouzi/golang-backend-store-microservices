package adminendpoints

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Apouzi/Golang-Admin-Service/helpers"
	"github.com/go-chi/chi/v5"
)

func (route *AdminRoutes) CreateProduct(w http.ResponseWriter, r *http.Request) {

	fmt.Println("hit!")

	productFromRequest := &ProductCreate{}

	err := helpers.ReadJSON(w, r, &productFromRequest)
	if err != nil {
		fmt.Println("read json error:", err)
	}
	fmt.Println("create product at admin:", productFromRequest)
	sendOff, err := json.Marshal(productFromRequest)
	if err != nil {
		fmt.Println("There was an error marshalling data:", err)
	}
	createdProductResult, err := http.Post("http://dblayer:8080/db/products/", "application/json", bytes.NewReader(sendOff))
	if err != nil {
		fmt.Println("Error connecting to DBLayer:", err)
		helpers.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error":   "Database service unavailable",
			"details": err.Error(),
		})
		return
	}

	if createdProductResult == nil {
		fmt.Println("Received nil response from DBLayer")
		helpers.WriteJSON(w, http.StatusServiceUnavailable, map[string]string{
			"error": "Database service returned no response",
		})
		return
	}

	prodcreate := &ProductCreateRetrieve{}

	decode := json.NewDecoder(createdProductResult.Body)
	err = decode.Decode(prodcreate)
	if err != nil {
		fmt.Println("Error decoding response from DBLayer:", err)
		helpers.WriteJSON(w, http.StatusInternalServerError, map[string]string{
			"error":   "Failed to parse database response",
			"details": err.Error(),
		})
		return
	}
	fmt.Println("attempted result is", prodcreate)
	helpers.WriteJSON(w, http.StatusAccepted, prodcreate)

}

func (route *AdminRoutes) CreateVariation(w http.ResponseWriter, r *http.Request) {
	ProductID := chi.URLParam(r, "ProductID")
	variation := VariationCreate{}
	helpers.ReadJSON(w, r, &variation)

	// Check if product exists, if not, then return false
	pil := ProductRetrieve{}
	url := "http://dblayer:8080/products/" + ProductID
	resp, err := http.Get(url)
	if err != nil {
		helpers.WriteJSON(w, 500, "Error getting data from database")
	}

	if resp.StatusCode == 404 {
		helpers.ErrorJSON(w, errors.New("could not find the coresponding product to create variation"), 404)
		return
	}

	jDecode := json.NewDecoder(resp.Body)
	if err = jDecode.Decode(&pil); err != nil || pil.ProductID == 0 {
		fmt.Println("There is an error decoding!", err)
	}

	url = "http://dblayer:8080/products/" + ProductID + "/variation"
	varbytes, err := json.Marshal(variation)
	if err != nil {
		helpers.ErrorJSON(w, errors.New("could not martial inputted variation"), 404)
	}
	varReader := bytes.NewReader(varbytes)
	resp, err = http.Post(url, "application/json", varReader)
	if err != nil {
		helpers.WriteJSON(w, 500, "Error getting data from database")
	}
	if resp.StatusCode == 404 {
		helpers.ErrorJSON(w, errors.New("could not find the coresponding product to create variation"), 404)
	}
	verify := variCrtd{}
	jDecode = json.NewDecoder(resp.Body)
	if err = jDecode.Decode(&verify); err != nil || pil.ProductID == 0 {
		fmt.Println("There is an error decoding!", err)
	}

	helpers.WriteJSON(w, 200, verify)

}


func (route *AdminRoutes) EditProduct(w http.ResponseWriter, r *http.Request){
	ProdID := chi.URLParam(r, "ProductID")
	prodEdit := &ProductEdit{}
	helpers.ReadJSON(w,r, prodEdit)
	
	fmt.Println(prodEdit)
	url := "http://dblayer:8080/products/"+ProdID
	fmt.Println("url:", url)
	prodBytes, err := json.Marshal(prodEdit)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("PATCH",url,bytes.NewReader(prodBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	err = responseDecode.Decode(prodEdit)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}


	helpers.WriteJSON(w,http.StatusAccepted,&prodEdit)
	
}



func (route *AdminRoutes) EditVariation(w http.ResponseWriter, r *http.Request){
	r.Header.Get("Authorization")
	VarID := chi.URLParam(r, "VariationID")
	VaritEdit := &VariationEdit{}
	helpers.ReadJSON(w,r, VaritEdit)
	url := "http://dblayer:8080/variation/" + VarID

	varitBytes, err := json.Marshal(VaritEdit)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("PATCH",url,bytes.NewReader(varitBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	varitReturn := &VariationEdit{}
	err = responseDecode.Decode(varitReturn)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}


	helpers.WriteJSON(w, http.StatusAccepted, VaritEdit)
}



func (route *AdminRoutes) AddAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}
	att := Attribute{}
	
	err := helpers.ReadJSON(w,r,&att)
	if err != nil{
		helpers.ErrorJSON(w, err, 500)
		return
	}

	url := "http://dblayer:8080/variation/" + VarID + "/attribute"

	attributeBytes, err := json.Marshal(att)
	if err != nil{
		fmt.Println(err)
	}
	// prodDecode:= bytes.NewReader(prodBytes)
	req, err := http.NewRequest("POST",url,bytes.NewReader(attributeBytes))
	if err != nil{
		fmt.Println("error trying to post to create prime category",err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Send the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	responseDecode := json.NewDecoder(resp.Body)
	attRet := &AddedSendBack{}
	err = responseDecode.Decode(attRet)
	if err != nil{
		fmt.Println("error trying to decode",err)
	}
	
	helpers.WriteJSON(w, 200, attRet)
}

func (route *AdminRoutes) DeleteAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}

	sql, err := route.DB.Exec("DELETE FROM tblProductAttribute WHERE Variation_ID = ? AND AttributeName = ?", VarID, AttName)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1{
		helpers.WriteJSON(w, 200, "Not Deleted")
		return
	}
	
	helpers.WriteJSON(w, 200, "Deleted")
}


func (route *AdminRoutes) UpdateAttribute(w http.ResponseWriter, r *http.Request){
	VarID := chi.URLParam(r,"VariationID")
	AttName := chi.URLParam(r, "AttributeName")
	if VarID == ""{
		helpers.ErrorJSON(w, errors.New("please input VariationID"),400)
		return
	}
	AttRead := Attribute{} 
	helpers.ReadJSON(w,r,&AttRead)
	fmt.Println(AttRead.Attribute,"variatio id and atttribute name", VarID,AttName)
	sql, err := route.DB.Exec("UPDATE tblProductAttribute SET AttributeName = ? WHERE Variation_ID = ? AND AttributeName = ?",AttRead.Attribute ,VarID, AttName)
	if err != nil{
		helpers.ErrorJSON(w,err, 400)
		return
	}

	nRows, _ := sql.RowsAffected()
	if nRows < 1{
		helpers.WriteJSON(w, 200, "No Updated Happened")
		return
	}
	
	helpers.WriteJSON(w, 200, "Attribute Updated")
}


func(route *AdminRoutes) CreateProductSize(w http.ResponseWriter, r *http.Request){
	variationID := chi.URLParam(r, "VariationID")
	productID := chi.URLParam(r, "ProductID")

	// Ensure DateCreated and ModifiedDate are initialized to now
	resp, err := http.Post("http://dblayer:8080/products/"+productID+"/variation/"+variationID+"/size","application/json", r.Body)
	if err != nil {
		fmt.Println("There was an error posting product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error posting product size"), 500)
		return
	}
	defer resp.Body.Close()

	var respData struct {
		SizeID int64 `json:"size_id"`
	}
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&respData); err != nil {
		fmt.Println("Error decoding response from DBLayer:", err)
		helpers.ErrorJSON(w, errors.New("failed to parse database response"), 500)
		return
	}
	sizeID := respData.SizeID
	if sizeID == 0 {
		fmt.Println("There was an error getting last insert id for product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error getting last insert id for product size"), 500)
		return
	}

	helpers.WriteJSON(w, http.StatusCreated, sizeID)
}  

func (route *AdminRoutes) DeleteProductSize(w http.ResponseWriter, r *http.Request){
	sizeID := chi.URLParam(r, "ProductSizeID")
	prdCheck := ProductSize{}


	response, err := http.Get("http://dblayer:8080/product-size/"+sizeID)
	if err != nil {
		fmt.Println("There was an error retrieving product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error retrieving product size"), 500)
		return
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		fmt.Println("Failed to retrieve product size, status code:", response.StatusCode)
		helpers.ErrorJSON(w, errors.New("failed to retrieve product size"), response.StatusCode)
		return
	}
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&prdCheck); err != nil {
		fmt.Println("Error decoding response from DBLayer:", err)
		helpers.ErrorJSON(w, errors.New("failed to parse database response"), 500)
		return
	}

	if prdCheck.SizeID == nil {
		fmt.Println("Product size not found for deletion")
		helpers.ErrorJSON(w, errors.New("product size not found for deletion"), 404)
		return
	}

	req, err := http.NewRequest("DELETE", "http://dblayer:8080/product-size/"+sizeID, nil)
	if err != nil {
		fmt.Println("There was an error creating delete request for product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error creating delete request for product size"), 500)
		return
	}

	client := &http.Client{}
	deleteResp, err := client.Do(req)
	if err != nil {
		fmt.Println("There was an error sending delete request for product size:", err)
		helpers.ErrorJSON(w, errors.New("there was an error deleting product size"), 500)
		return
	}
	defer deleteResp.Body.Close()

	if deleteResp.StatusCode != http.StatusOK {
		fmt.Println("Delete product size failed, status code:", deleteResp.StatusCode)
		helpers.ErrorJSON(w, errors.New("failed to delete product size"), deleteResp.StatusCode)
		return
	}

	verifyResp, err := http.Get("http://dblayer:8080/productsize/" + sizeID)
	if err != nil {
		fmt.Println("There was an error verifying product size deletion:", err)
		helpers.ErrorJSON(w, errors.New("could not verify product size deletion"), 500)
		return
	}
	defer verifyResp.Body.Close()

	if verifyResp.StatusCode == http.StatusOK {
		verifyCheck := ProductSize{}
		verifyDecoder := json.NewDecoder(verifyResp.Body)
		if err := verifyDecoder.Decode(&verifyCheck); err == nil && verifyCheck.SizeID != nil {
			helpers.ErrorJSON(w, errors.New("product size was not deleted"), 500)
			return
		}
	}

	helpers.WriteJSON(w, http.StatusOK, "Product size deleted successfully")
}
func (route *AdminRoutes) DeleteProductVariation(w http.ResponseWriter, r *http.Request){
}
func (route *AdminRoutes) DeleteProduct(w http.ResponseWriter, r *http.Request){
}