CREATE TABLE IF NOT EXISTS tblCategoriesPrime (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesSub (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);


CREATE TABLE IF NOT EXISTS tblCategoriesFinal (
    Category_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    CategoryName VARCHAR(255) NOT NULL,
    CategoryDescription TEXT
);

CREATE TABLE IF NOT EXISTS tblCatPrimeSub (
  CatPrimeID INT NOT NULL,
  CatSubID INT NOT NULL,
  FOREIGN KEY (CatPrimeID) REFERENCES tblCategoriesPrime (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (Category_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblCatSubFinal (
  CatSubID INT NOT NULL,
  CatFinalID INT NOT NULL,
  FOREIGN KEY (CatSubID) REFERENCES tblCategoriesSub (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (CatFinalID) REFERENCES tblCategoriesFinal (Category_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblCatFinalProd (
  CatFinalID INT NOT NULL,
  Product_ID INT NOT NULL,
  FOREIGN KEY (CatFinalID) REFERENCES tblCategoriesFinal (Category_ID) ON DELETE CASCADE,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
  
);