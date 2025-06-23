CREATE TABLE IF NOT EXISTS tblLocation (
    location_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(10, 8),
    street_address VARCHAR(300)
);

CREATE TABLE IF NOT EXISTS tblLocationProduct (
    location_product_id INT NOT NULL,
    product_id INT NOT NULL,
    inventory_id INT NOT NULL,
    PRIMARY KEY (location_product_id, product_id),
    FOREIGN KEY (location_product_id) REFERENCES tblLocation(location_id),
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(Variation_ID),
    FOREIGN KEY (inventory_id) REFERENCES tblInventoryProductDetail(inventory_id)
);

CREATE TABLE IF NOT EXISTS tblInventoryProductDetail (
    inventory_id INT AUTO_INCREMENT PRIMARY KEY,
    quantity INT NOT NULL,
    product_id INT NOT NULL,
    location_id INT NOT NULL,
    description TEXT,
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(Variation_ID),
    FOREIGN KEY (location_id) REFERENCES tblLocation(location_id)
);

CREATE TABLE IF NOT EXISTS tblInventoryLocationTransfers (
    transfers_id INT AUTO_INCREMENT PRIMARY KEY,
    source_location_id INT NOT NULL,
    destination_location_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    transfer_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    status VARCHAR(15),
    FOREIGN KEY (source_location_id) REFERENCES tblLocation(location_id),
    FOREIGN KEY (destination_location_id) REFERENCES tblLocation(location_id),
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(Variation_ID)
);
