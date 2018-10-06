-- +goose Up
-- SQL in this section is executed when the migration is applied.

-- The intent of this table was to use TimescaleDB
CREATE TABLE metrics_development (
  time TIMESTAMPTZ NOT NULL,
  node_name TEXT NOT NULL,
  cpu_used INT NOT NULL,
  mem_used INT NOT NULL
);

CREATE TABLE metrics_test (
  LIKE metrics_development
  INCLUDING ALL
);

-- If you look at the docs, you do something similar: https://docs.timescale.com/v0.12/getting-started/creating-hypertables
-- SELECT create_hypertable('metrics_development', 'time');
-- SELECT create_hypertable('metrics_test', 'time');

-- But it took too long to figure out even, with the offical docker image.
-- Probably need to call CREATE EXTENSION IF NOT EXISTS timescaledb CASCADE; somewhere?

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE metrics_development, metrics_test;
