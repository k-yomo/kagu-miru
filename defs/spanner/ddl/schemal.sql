CREATE TABLE item_categories (
    id STRING(256) NOT NULL,
    name STRING(256) NOT NULL,
    level INT64 NOT NULL,
    parent_id STRING(256),
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE TABLE items (
    id STRING(256) NOT NULL,
    name STRING(256) NOT NULL,
    description STRING(16384) NOT NULL,
    status INT64 NOT NULL,
    url STRING(1024) NOT NULL,
    affiliate_url STRING(1024) NOT NULL,
    price INT64 NOT NULL,
    image_urls ARRAY<STRING(1024)> NOT NULL,
    average_rating FLOAT64 NOT NULL,
    review_count INT64 NOT NULL,
    category_id STRING(256) NOT NULL,
    jan_code STRING(256),
    platform STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (category_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE INDEX items_by_updated_at ON items (updated_at DESC);
