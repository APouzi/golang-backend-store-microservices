package initializingpopulation

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
)

func PopulateProductTables(db *sql.DB) {

	query, err := ioutil.ReadFile("./sql/Products.sql")

	if err != nil {
		log.Fatal("Error when loading sql file", err)
	}

	_, err = db.Exec(string(query))

	if err != nil {
		log.Fatal("Couldn't complete the execution of the file", err)
	}

	for i := 0.00; i <= 10; i++ {
		_, err = db.Exec("INSERT INTO tblProducts (Product_Name, Product_Description) VALUES(?,?)", "testProductPopulate", "This is a description!")
		if err != nil {
			log.Fatal("Error with tblProducts")
		}
	}
	fmt.Println("init tables!")
	// _, err = db.Exec("INSERT INTO tblCategoriesPrime(CategoryName, CategoryDescription) VALUES(?,?)", "Test Category","This is a description category")
	// if err != nil{
	// 	log.Println("Error inserting tblCategoriesPrime")
	// }
	// _,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(Product_ID, Category_ID) VALUES(1,1)")
	// if err != nil{
	// 	log.Println("Error inserting into tblProductsCategoriesPrime", err)
	// }
	// _,err =  db.Exec("INSERT INTO tblProductsCategoriesPrime(Product_ID, Category_ID) VALUES(2,1)")
	// if err != nil{
	// 	log.Println("Error inserting into tblProductCategoies")
	// }

	// resultProd := database.Product{}
	// row := db.QueryRow("select Product_ID, Product_Name, Product_Description from tblProducts where Product_ID = ?", 4)
	// if row == nil {
	// 	fmt.Println("Nothing returned!")
	// }
	// err = row.Scan(&resultProd.Product_ID, &resultProd.Product_Name, &resultProd.Product_Description)
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// listPrint := []database.Product{}
	// rows, err := db.Query("SELECT tblProducts.Product_ID, tblProducts.Product_Name, tblProducts.Product_Description, tblProducts.Product_Price FROM tblProducts JOIN tblProductsCategoriesPrime ON tblProducts.Product_ID = tblProductsCategoriesPrime.Product_ID JOIN tblCategoriesPrime ON tblProductsCategoriesPrime.Category_ID = tblCategoriesPrime.Category_ID WHERE tblCategoriesPrime.CategoryName = ?", "Test Category" )
	if err != nil {
		log.Fatal("Error with category", err)
	}
	// defer rows.Close()
	// for rows.Next(){
	// 	resultProd2 := database.Product{}
	// 	rows.Scan(&resultProd2.Product_ID, &resultProd2.Product_Name, &resultProd2.Product_Description, &resultProd2.Product_Price)
	// 	listPrint = append(listPrint, resultProd2)
	// }
	// fmt.Println(resultProd)
	// fmt.Println("Population of tables has been completed!")
	// fmt.Println("Categories:", listPrint)
	// for _, v := range listPrint {
	// 	fmt.Println(v.Product_Name)
	// }
	//Test this out!

}