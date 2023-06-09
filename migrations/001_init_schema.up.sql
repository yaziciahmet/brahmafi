CREATE TABLE pool(
    address text primary key
);

CREATE TABLE pool_snapshot(
    id uuid primary key default gen_random_uuid(),
    pool_id text not null,
    token0_balance bigint not null,
    token1_balance bigint not null,
    tick bigint not null,
    block_number bigint not null,
    taken_at timestamptz default now() not null
);

ALTER TABLE pool_snapshot ADD CONSTRAINT pool_fk
FOREIGN KEY (pool_id) REFERENCES pool(address);

CREATE INDEX pool_snapshot_pool_idx ON pool_snapshot(pool_id);

CREATE INDEX pool_snapshot_block_number ON pool_snapshot(block_number);
