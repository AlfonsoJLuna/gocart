-------------------------------------------------------------------------------
-- Config
-------------------------------------------------------------------------------
CREATE TABLE config (
    key         TEXT        PRIMARY KEY COLLATE NOCASE,
    value       TEXT        NOT NULL,
    updated_at  TEXT        NOT NULL DEFAULT (unixepoch())
);

-------------------------------------------------------------------------------
-- Currencies
-------------------------------------------------------------------------------
CREATE TABLE currencies (
    id                  BLOB        PRIMARY KEY,
    created_at          INTEGER     NOT NULL DEFAULT (unixepoch()),
    updated_at          INTEGER     NOT NULL DEFAULT (unixepoch()),
    is_enabled          INTEGER     NOT NULL DEFAULT TRUE,

    iso_code            TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(iso_code) = 3),
    name                TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(name) <= 50),
    name_alt            TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(name_alt) <= 50),

    decimals            INTEGER     NOT NULL DEFAULT 2,
    fx_rate             REAL        NOT NULL DEFAULT 1.0
);

-------------------------------------------------------------------------------
-- Countries
-------------------------------------------------------------------------------
CREATE TABLE countries (
    id                  TEXT        PRIMARY KEY,
    created_at          TEXT        NOT NULL DEFAULT (unixepoch()),
    updated_at          TEXT        NOT NULL DEFAULT (unixepoch()),
    is_enabled          INTEGER     NOT NULL DEFAULT TRUE,

    iso_code            TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(iso_code) = 2),
    name                TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(name) <= 50),
    name_alt            TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(name_alt) <= 50),

    is_eu               INTEGER     NOT NULL DEFAULT FALSE,
    vat_rate            REAL        NOT NULL DEFAULT 0.0,
    currency_id         TEXT        REFERENCES currencies (id) ON DELETE SET NULL
);

-------------------------------------------------------------------------------
-- Regions
-------------------------------------------------------------------------------
CREATE TABLE regions (
    id                  TEXT        PRIMARY KEY,
    created_at          TEXT        NOT NULL DEFAULT (unixepoch()),
    updated_at          TEXT        NOT NULL DEFAULT (unixepoch()),
    is_enabled          INTEGER     NOT NULL DEFAULT TRUE,

    country_id          TEXT        NOT NULL REFERENCES countries (id) ON DELETE CASCADE,
    name                TEXT        NOT NULL COLLATE NOCASE CHECK (length(name) <= 50),
    name_alt            TEXT        NOT NULL COLLATE NOCASE CHECK (length(name_alt) <= 50),
    is_eu               INTEGER     NOT NULL DEFAULT FALSE,
    vat_rate            REAL        NOT NULL DEFAULT 0.0,

    UNIQUE (country_id, name),
    UNIQUE (country_id, name_alt)
);

-------------------------------------------------------------------------------
-- Users
-------------------------------------------------------------------------------
CREATE TABLE users (
    id                      TEXT        PRIMARY KEY,
    created_at              TEXT        NOT NULL DEFAULT (unixepoch()),
    updated_at              TEXT        NOT NULL DEFAULT (unixepoch()),
    is_enabled              INTEGER     NOT NULL DEFAULT TRUE,

    email                   TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(email) <= 50),
    username                TEXT        NOT NULL UNIQUE COLLATE NOCASE CHECK (length(username) <= 30),
    password_hash           TEXT        NOT NULL,

    currency_id             TEXT        REFERENCES currencies (id) ON DELETE SET NULL,
    country_id              TEXT        REFERENCES countries (id) ON DELETE SET NULL,
    region_id               TEXT        REFERENCES regions (id) ON DELETE SET NULL,

    vat_number              TEXT        NOT NULL DEFAULT '' CHECK (length(vat_number) <= 15),
    phone                   TEXT        NOT NULL DEFAULT '' CHECK (length(phone) <= 30),

    billing_full_name       TEXT        NOT NULL DEFAULT '' CHECK (length(billing_full_name) <= 48),
    billing_company_name    TEXT        NOT NULL DEFAULT '' CHECK (length(billing_company_name) <= 48),
    billing_address_1       TEXT        NOT NULL DEFAULT '' CHECK (length(billing_address_1) <= 48),
    billing_address_2       TEXT        NOT NULL DEFAULT '' CHECK (length(billing_address_2) <= 48),
    billing_postal_code     TEXT        NOT NULL DEFAULT '' CHECK (length(billing_postal_code) <= 10),
    billing_city            TEXT        NOT NULL DEFAULT '' CHECK (length(billing_city) <= 37),

    shipping_full_name      TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_full_name) <= 48),
    shipping_company_name   TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_company_name) <= 48),
    shipping_address_1      TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_address_1) <= 48),
    shipping_address_2      TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_address_2) <= 48),
    shipping_postal_code    TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_postal_code) <= 10),
    shipping_city           TEXT        NOT NULL DEFAULT '' CHECK (length(shipping_city) <= 37),

    is_verified             INTEGER     NOT NULL DEFAULT FALSE,
    is_business             INTEGER     NOT NULL DEFAULT FALSE,
    is_subscribed_info      INTEGER     NOT NULL DEFAULT FALSE,
    is_subscribed_promos    INTEGER     NOT NULL DEFAULT FALSE,
    is_admin                INTEGER     NOT NULL DEFAULT FALSE,

    last_login_at           TEXT,
    last_user_change_at     TEXT,
    last_pass_change_at     TEXT,
    last_email_change_at    TEXT
);

CREATE INDEX idx_users_currency_id  ON users (currency_id);
CREATE INDEX idx_users_country_id   ON users (country_id);
CREATE INDEX idx_users_region_id    ON users (region_id);
