CREATE TABLE IF NOT EXISTS orders(
    orderid TEXT PRIMARY KEY,
    productid TEXT NOT NULL,
    consumerid TEXT NOT NULL,
    vendorid TEXT NOT NULL,
    amount INT NOT NULL,
    orderprice INT
);