-- +goose Up
-- SQL in this section is executed when the migration is applied.
CREATE TABLE metrics_development (
  time TIMESTAMPTZ NOT NULL,
  node_name TEXT NOT NULL,
  process_name TEXT NULL,
  cpu_used REAL NOT NULL,
  mem_used REAL NOT NULL
);

CREATE UNIQUE INDEX names_idx ON metrics_development (node_name, process_name);

CREATE TABLE metrics_test (
  LIKE metrics_development
  INCLUDING ALL
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE metrics_development, metrics_test;
