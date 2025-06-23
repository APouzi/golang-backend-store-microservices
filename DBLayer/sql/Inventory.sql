CREATE TABLE IF NOT EXISTS location (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    latitude DECIMAL(10, 8),
    longitude DECIMAL(10, 8),
    street_address VARCHAR(300),
);

CREATE TABLE IF NOT EXISTS location_product (
    location_id INT NOT NULL,
    product_id INT NOT NULL,
    inventory_id INT NOT NULL,
    PRIMARY KEY (location_id, product_id),
    FOREIGN KEY (location_id) REFERENCES location(id),
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(product_id),
    FOREIGN KEY (inventory_id) REFERENCES inventory_product_detail(inventory_id)
);



CREATE TABLE IF NOT EXISTS inventory_product_detail (
    inventory_id INT AUTO_INCREMENT PRIMARY KEY,
    quantity INT NOT NULL,
    product_id INT NOT NULL,
    location_id INT NOT NULL,
    description TEXT,
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(product_id),
    FOREIGN KEY (location_id) REFERENCES location(id)
);


CREATE TABLE IF NOT EXISTS transfers (
    transfers_id INT AUTO_INCREMENT PRIMARY KEY,
    source_location_id INT NOT NULL,
    destination_location_id INT NOT NULL,
    product_id INT NOT NULL,
    quantity INT NOT NULL,
    transfer_date DATETIME DEFAULT CURRENT_TIMESTAMP,
    description TEXT,
    status VARCHAR(15),
    FOREIGN KEY (source_location_id) REFERENCES location(id),
    FOREIGN KEY (destination_location_id) REFERENCES location(id),
    FOREIGN KEY (product_id) REFERENCES tblProductVariation(product_id)
);








