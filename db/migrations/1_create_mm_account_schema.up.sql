CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS item
(
    item_id          varchar(255) not null,
    institution_id   varchar(255),
    institution_name varchar(255),
    url              varchar(255),
    primary_color    varchar(255),
    logo             varchar(255),
    tenant_id        varchar(255),
    PRIMARY KEY (item_id)
);

CREATE TABLE IF NOT EXISTS account
(
    account_id               uuid DEFAULT uuid_generate_v4() not null,
    friendly_name            varchar(255),
    name                     varchar(255),
    mask                     varchar(255),
    plaid_id                 varchar(255),
    item_id                  varchar(255),
    official_name            varchar(255),
    available_balance        decimal,
    current_balance          decimal                         NOT NULL,
    starting_balance         decimal,
    iso_currency_code        varchar(255),
    unofficial_currency_code varchar(255),
    account_limit            decimal,
    is_primary_account       boolean,
    is_in_cash_flow          boolean,
    account_type             varchar(255),
    account_sub_type         varchar(255),
    tenant_id                varchar(255),
    PRIMARY KEY (account_id),
    FOREIGN KEY (item_id) REFERENCES item (item_id)
);