CREATE TABLE IF NOT EXISTS tblOrders (
    OrderID INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    OrderNumber VARCHAR NOT NULL,
    Guest BOOLEAN NOT Null,
    UserEmail VARCHAR,
    UserID INT,
    AdminNotes VARCHAR NOT NULL,
    OrderProductList JSON NOT NULL,
    FOREIGN KEY (UserID) REFERENCES tblUser (UserID)
)


