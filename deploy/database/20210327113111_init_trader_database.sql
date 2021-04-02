-- +goose Up
CREATE SCHEMA trader;

CREATE TABLE trader.items
(
    id         bigserial                                  NOT NULL CONSTRAINT items_pk PRIMARY KEY,
    name       varchar(200) DEFAULT ''::character varying NOT NULL,
    image_url  varchar(200) DEFAULT ''::character varying NOT NULL,
    creator_id bigint       DEFAULT 0                     NOT NULL,
    updated_at timestamp    DEFAULT now()                 NOT NULL,
    created_at timestamp    DEFAULT now()                 NOT NULL
);

CREATE TABLE trader.orders
(
    id         bigserial               NOT NULL CONSTRAINT orders_pk PRIMARY KEY,
    item_id    bigint    DEFAULT 0     NOT NULL,
    creator_id bigint    DEFAULT 0     NOT NULL,
    order_type smallint  DEFAULT 0     NOT NULL,
    price      decimal   DEFAULT 0     NOT NULL,
    status     smallint  DEFAULT 0     NOT NULL,
    updated_at timestamp DEFAULT now() NOT NULL,
    created_at timestamp DEFAULT now() NOT NULL
);

CREATE TABLE trader.users
(
    id            bigserial                                  NOT NULL CONSTRAINT users_pk PRIMARY KEY,
    display_name  varchar(200) DEFAULT ''::character varying NOT NULL,
    email         varchar(400) DEFAULT ''::character varying NOT NULL,
    password_hash varchar(200) DEFAULT ''::character varying NOT NULL,
    last_login_at timestamp,
    updated_at    timestamp    DEFAULT now()                 NOT NULL,
    created_at    timestamp    DEFAULT now()                 NOT NULL
);

CREATE TABLE trader.transactions
(
    id            bigserial               NOT NULL CONSTRAINT transactions_p PRIMARY KEY,
    make_order_id bigint    DEFAULT 0     NOT NULL,
    take_order_id bigint    DEFAULT 0     NOT NULL,
    final_price   decimal   DEFAULT 0     NOT NULL,
    updated_at    timestamp DEFAULT now() NOT NULL,
    created_at    timestamp DEFAULT now() NOT NULL
);


-- +goose Down
DROP TABLE IF EXISTS trader.items;
DROP TABLE IF EXISTS trader.orders;
DROP TABLE IF EXISTS trader.users;
DROP TABLE IF EXISTS trader.transactions;
