CREATE TABLE pool(
    id uuid primary key,
    address text not null,
    chain_id text not null,
    chain_name text not null
);

CREATE TABLE pool_snapshot(
    id uuid primary key default gen_random_uuid(),
    pool_id uuid not null,
    token0_balance int not null,
    token1_balance int not null,
    tick int not null,
    block_number int not null
);
