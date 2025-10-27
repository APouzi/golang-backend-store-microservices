-- ecommerce_core.sql
-- E-commerce order + payment + fulfillment schema (MySQL 8.0+)

START TRANSACTION;

-- ============================================================================
-- Customers & Addresses
-- ============================================================================

CREATE TABLE IF NOT EXISTS customers (
  customer_id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  email              VARCHAR(320) NOT NULL,
  full_name          VARCHAR(255),
  created_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  UNIQUE KEY ux_customers_email (email)
);

CREATE TABLE IF NOT EXISTS addresses (
  address_id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  full_name          VARCHAR(255),
  line1              VARCHAR(255) NOT NULL,
  line2              VARCHAR(255),
  city               VARCHAR(128) NOT NULL,
  region             VARCHAR(128),
  postal_code        VARCHAR(32) NOT NULL,
  country            CHAR(2) NOT NULL,
  phone              VARCHAR(32),
  created_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6)
);

-- ============================================================================
-- Orders
-- ============================================================================

CREATE TABLE IF NOT EXISTS orders (
  order_id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_number       VARCHAR(32) NOT NULL,
  customer_id        BIGINT UNSIGNED,
  email              VARCHAR(320) NOT NULL,
  billing_address_id BIGINT UNSIGNED,
  shipping_address_id BIGINT UNSIGNED,

  currency           CHAR(3) NOT NULL,--- USD, HKD, EUR, MAD
  subtotal_cents     BIGINT NOT NULL DEFAULT 0, --- $19.99 = 1999 cents
  discount_cents     BIGINT NOT NULL DEFAULT 0,
  shipping_cents     BIGINT NOT NULL DEFAULT 0,
  tax_cents          BIGINT NOT NULL DEFAULT 0,
  total_cents        BIGINT NOT NULL, --- subtotal - discount + tax + shipping

  status             ENUM('created','awaiting_payment','paid','fulfilled','cancelled','closed') NOT NULL DEFAULT 'created',

  placed_at          DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  provider           VARCHAR(32),
  provider_order_id  VARCHAR(128),
  metadata           JSON NOT NULL,

  created_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
  updated_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),

  UNIQUE KEY ux_orders_order_number (order_number),
  KEY ix_orders_status (status),
  KEY ix_orders_provider (provider, provider_order_id),
  KEY ix_orders_email (email),
  KEY ix_orders_placed_at (placed_at),

  FOREIGN KEY (customer_id) REFERENCES customers(customer_id),
  FOREIGN KEY (billing_address_id) REFERENCES addresses(address_id),
  FOREIGN KEY (shipping_address_id) REFERENCES addresses(address_id)
);

-- ============================================================================
-- Order Items
-- ============================================================================

CREATE TABLE IF NOT EXISTS order_items (
  order_item_id         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id              BIGINT UNSIGNED NOT NULL,
  product_id            BIGINT,
  sku                   VARCHAR(128),
  title                 VARCHAR(255) NOT NULL,
  variation             JSON,

  qty                   INT NOT NULL CHECK (qty > 0),
  currency              CHAR(3) NOT NULL,

  unit_price_cents      BIGINT NOT NULL,
  line_subtotal_cents   BIGINT NOT NULL,
  line_discount_cents   BIGINT NOT NULL DEFAULT 0,
  line_tax_cents        BIGINT NOT NULL DEFAULT 0,
  line_total_cents      BIGINT NOT NULL,

  KEY ix_order_items_order (order_id),
  KEY ix_order_items_sku (sku),

  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
  FOREIGN KEY (product_id) REFERENCES tblProductSize (Size_ID) ON DELETE CASCADE
);

-- ============================================================================
-- Payments & Refunds
-- ============================================================================

CREATE TABLE IF NOT EXISTS payments (
  payment_id           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id             BIGINT UNSIGNED NOT NULL,
  provider             VARCHAR(32) NOT NULL,
  provider_payment_id  VARCHAR(128) NOT NULL,
  method_brand         VARCHAR(32),
  last4                VARCHAR(4),
  status               ENUM('authorized','captured','failed','refunded','partially_refunded')
                         NOT NULL,
  amount_cents         BIGINT NOT NULL,
  currency             CHAR(3) NOT NULL,
  raw_response         JSON NOT NULL,

  created_at           DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  UNIQUE KEY ux_payments_provider_ref (provider, provider_payment_id),
  KEY ix_payments_order (order_id),
  KEY ix_payments_status (status),

  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS refunds (
  refund_id            BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  payment_id           BIGINT UNSIGNED NOT NULL,
  amount_cents         BIGINT NOT NULL,
  reason               VARCHAR(255),
  provider_refund_id   VARCHAR(128),
  created_at           DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  KEY ix_refunds_payment (payment_id),

  FOREIGN KEY (payment_id) REFERENCES payments(payment_id) ON DELETE CASCADE
);

-- ============================================================================
-- Promotions & Discounts
-- ============================================================================

CREATE TABLE IF NOT EXISTS promotions (
  promotion_id       BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  code               VARCHAR(64) UNIQUE,
  type               ENUM('fixed','percent','bogo') NOT NULL,
  value              DECIMAL(12,4) NOT NULL,
  max_uses           INT,
  starts_at          DATETIME(6),
  ends_at            DATETIME(6),
  metadata           JSON NOT NULL,

  KEY ix_promotions_active (starts_at, ends_at)
);

CREATE TABLE IF NOT EXISTS order_discounts (
  order_discount_id  BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id           BIGINT UNSIGNED NOT NULL,
  promotion_id       BIGINT UNSIGNED,
  code               VARCHAR(64),
  amount_cents       BIGINT NOT NULL,
  allocation         JSON NOT NULL,

  KEY ix_order_discounts_order (order_id),
  KEY ix_order_discounts_code (code),

  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE,
  FOREIGN KEY (promotion_id) REFERENCES promotions(promotion_id)
);

-- ============================================================================
-- Shipments & Shipment Items
-- ============================================================================

CREATE TABLE IF NOT EXISTS shipments (
  shipment_id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id           BIGINT UNSIGNED NOT NULL,
  carrier            VARCHAR(64),
  service            VARCHAR(64),
  tracking_number    VARCHAR(128),
  shipped_at         DATETIME(6),
  status             ENUM('pending','shipped','delivered','returned')
                       NOT NULL DEFAULT 'pending',

  KEY ix_shipments_order (order_id),
  KEY ix_shipments_tracking (tracking_number),
  KEY ix_shipments_status (status),

  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS shipment_items (
  shipment_item_id   BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  shipment_id        BIGINT UNSIGNED NOT NULL,
  order_item_id      BIGINT UNSIGNED NOT NULL,
  qty                INT NOT NULL CHECK (qty > 0),

  KEY ix_shipment_items_shipment (shipment_id),
  KEY ix_shipment_items_order_item (order_item_id),

  FOREIGN KEY (shipment_id) REFERENCES shipments(shipment_id) ON DELETE CASCADE,
  FOREIGN KEY (order_item_id) REFERENCES order_items(order_item_id) ON DELETE CASCADE
);

-- ============================================================================
-- Inventory Movements
-- ============================================================================

CREATE TABLE IF NOT EXISTS inventory_movements (
  movement_id        BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  sku                VARCHAR(128) NOT NULL,
  variation_key      VARCHAR(255),
  qty_delta          INT NOT NULL,
  reason             ENUM('reserve','sale','cancel','refund','adjustment') NOT NULL,
  ref_type           ENUM('order','refund','manual'),
  ref_id             BIGINT UNSIGNED,
  created_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  KEY ix_inventory_movements_sku (sku),
  KEY ix_inventory_movements_created (created_at),
  KEY ix_inventory_movements_ref (ref_type, ref_id)
);

-- ============================================================================
-- Order Events
-- ============================================================================

CREATE TABLE IF NOT EXISTS order_events (
  order_event_id     BIGINT UNSIGNED NOT NULL AUTO_INCREMENT PRIMARY KEY,
  order_id           BIGINT UNSIGNED NOT NULL,
  type               VARCHAR(64) NOT NULL,
  details            JSON NOT NULL,
  created_at         DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),

  KEY ix_order_events_order (order_id),
  KEY ix_order_events_type (type),
  KEY ix_order_events_created (created_at),

  FOREIGN KEY (order_id) REFERENCES orders(order_id) ON DELETE CASCADE
);

COMMIT;
