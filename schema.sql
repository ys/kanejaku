CREATE EXTENSION "uuid-ossp";

CREATE TABLE metrics(
  id uuid primary key default uuid_generate_v4(),
  key text,
  value float,
  timestamp timestamptz
);

CREATE UNIQUE INDEX index_on_metrics_for_key_and_timestamp ON metrics(key, timestamp);
