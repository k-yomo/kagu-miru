CREATE TABLE item_categories (
    id STRING(256) NOT NULL,
    name STRING(256) NOT NULL,
    level INT64 NOT NULL,
    parent_id STRING(256),
    image_url STRING(1024),
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE TABLE amazon_browse_nodes (
    id STRING(256) NOT NULL,
    name STRING(256) NOT NULL,
    level INT64 NOT NULL,
    parent_id STRING(256),
    item_category_id STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES amazon_browse_nodes (id),
    FOREIGN KEY (item_category_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE TABLE rakuten_item_genres (
    id INT64 NOT NULL,
    name STRING(256) NOT NULL,
    level INT64 NOT NULL,
    parent_id INT64,
    item_category_id STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES rakuten_item_genres (id),
    FOREIGN KEY (item_category_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE TABLE yahoo_shopping_item_categories (
    id INT64 NOT NULL,
    name STRING(256) NOT NULL,
    level INT64 NOT NULL,
    parent_id INT64,
    item_category_id STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (parent_id) REFERENCES yahoo_shopping_item_categories (id),
    FOREIGN KEY (item_category_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE TABLE rakuten_tag_groups (
    id INT64 NOT NULL,
    name STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL
) PRIMARY KEY(id);

CREATE TABLE rakuten_tags (
    id INT64 NOT NULL,
    name STRING(256) NOT NULL,
    tag_group_id INT64 NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (tag_group_id) REFERENCES rakuten_tag_groups (id)
) PRIMARY KEY(id);

CREATE TABLE items (
    id STRING(256) NOT NULL,
    group_id STRING(256),
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
    brand_name STRING(256),
    colors ARRAY<STRING(256)>,
    width_range ARRAY<INT64>,
    depth_range ARRAY<INT64>,
    height_range ARRAY<INT64>,
    jan_code STRING(256),
    platform STRING(256) NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    FOREIGN KEY (category_id) REFERENCES item_categories (id)
) PRIMARY KEY(id);

CREATE INDEX items_by_updated_at ON items (updated_at);
CREATE INDEX items_by_group_id ON items (group_id);
