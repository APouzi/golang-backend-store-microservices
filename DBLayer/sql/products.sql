CREATE TABLE IF NOT EXISTS tblProducts (
  Product_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_Name VARCHAR(255) NOT NULL,
  Product_Description TEXT,
  PRIMARY_IMAGE VARCHAR(255) NULL,
  Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Modified_Date DATETIME NULL
);

CREATE TABLE IF NOT EXISTS tblProductVariation (
  Variation_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_ID INT NOT NULL,
  Variation_Name VARCHAR(255) NOT NULL,
  Variation_Description TEXT,
  PRIMARY_IMAGE VARCHAR(255) NULL,
  Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Modified_Date DATETIME NULL,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblProductSize(
  Size_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Size_Name VARCHAR(50) NOT NULL,
  Size_Description TEXT,
  Variation_ID INT NOT NULL,
  Variation_Price DECIMAL(10,2) NOT NULL,
  SKU VARCHAR(50),
  UPC VARCHAR(50),
  PRIMARY_IMAGE VARCHAR(255) NULL,
  Date_Created TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  Modified_Date DATETIME NULL,
  FOREIGN KEY (Variation_ID) REFERENCES tblProductVariation (Variation_ID) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS tblProductTaxCode(
  TaxCode_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  TaxCode_Name VARCHAR(50) NOT NULL,
  TaxCode_Description TEXT,
  TaxCode VARCHAR(50) NOT NULL,
  Provider ENUM('stripe','other') NOT NULL -- Payment provider for the tax code, for now only stripe is supported. 
);


CREATE TABLE IF NOT EXISTS tblProductSizeTaxCode (
  Size_ID INT NOT NULL,
  TaxCode_ID INT NOT NULL,
  PRIMARY KEY (Size_ID, TaxCode_ID),
  FOREIGN KEY (Size_ID) REFERENCES tblProductSize (Size_ID) ON DELETE CASCADE,
  FOREIGN KEY (TaxCode_ID) REFERENCES tblProductTaxCode (TaxCode_ID) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS tblProductInventoryLocation (
  Inv_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Variation_ID INT NOT NULL,
  Location_ID INT,
  Quantity INT NOT NULL,
  Location_At VARCHAR(255) NOT NULL,
  FOREIGN KEY (Variation_ID) REFERENCES tblProductVariation (Variation_ID) ON DELETE CASCADE,
  FOREIGN KEY (Location_ID) REFERENCES tblLocation (Location_ID) ON DELETE CASCADE
);


CREATE TABLE IF NOT EXISTS tblProductVariationImages (
  ImageID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Product_ID INT NOT NULL,
  ImageURL VARCHAR(255) NOT NULL,
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
);
CREATE TABLE IF NOT EXISTS tblProductAttribute (
  AttributeID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  Variation_ID INT NOT NULL,
  AttributeName VARCHAR(255) NOT NULL,
  FOREIGN KEY (Variation_ID) REFERENCES tblProductVariation (Variation_ID) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS tblDiscount (
  Discount_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  DiscountCode VARCHAR(255) NOT NULL,
  DiscountPercentage DECIMAL(5,2) NOT NULL,
  DiscountStartDate DATE,
  DiscountEndDate DATE
);

CREATE TABLE IF NOT EXISTS tblProductDiscount (
  Product_ID INT NOT NULL,
  Discount_ID INT NOT NULL,
  PRIMARY KEY (Product_ID, Discount_ID),
  FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE,
  FOREIGN KEY (Discount_ID) REFERENCES tblDiscount (Discount_ID) ON DELETE CASCADE
);




-- CREATE TABLE IF NOT EXISTS tblVariation (
--   Variation_ID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
--   Product_ID INT NOT NULL,
--   VariationName VARCHAR(255) NOT NULL,
--   VariationDescription TEXT,
--   FOREIGN KEY (Product_ID) REFERENCES tblProducts (Product_ID) ON DELETE CASCADE
-- );


