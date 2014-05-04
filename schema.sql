CREATE EXTENSION "uuid-ossp";

CREATE TABLE metrics(
    id uuid primary key default uuid_generate_v4(),
    key text,
    value float,
    timestamp timestamptz DEFAULT now()
);

CREATE INDEX index_on_metrics_for_key_and_timestamp ON metrics(key, timestamp);
CREATE INDEX index_on_metrics_1 ON metrics(key, round_timestamp(timestamp, 1));
CREATE INDEX index_on_metrics_10 ON metrics(key, round_timestamp(timestamp, 10));
CREATE INDEX index_on_metrics_30 ON metrics(key, round_timestamp(timestamp, 30));
CREATE INDEX index_on_metrics_60 ON metrics(key, round_timestamp(timestamp, 60));
CREATE INDEX index_on_metrics_300 ON metrics(key, round_timestamp(timestamp, 300));

-- Found on https://stackoverflow.com/questions/14300004/postgresql-equivalent-of-oracles-percentile-cont-function/14309370#14309370

CREATE OR REPLACE FUNCTION array_sort (ANYARRAY)
RETURNS ANYARRAY LANGUAGE SQL
AS $$
SELECT ARRAY(
    SELECT $1[s.i] AS "foo"
    FROM
    generate_series(array_lower($1,1), array_upper($1,1)) AS s(i)
    ORDER BY foo
);
$$;

CREATE OR REPLACE FUNCTION percentile_cont(myarray float[], percentile float)
RETURNS real AS
$$

DECLARE
ary_cnt INTEGER;
row_num float;
crn float;
frn float;
calc_result float;
new_array float[];
BEGIN
    ary_cnt = array_length(myarray,1);
    row_num = 1 + ( percentile * ( ary_cnt - 1 ));
    new_array = array_sort(myarray);

    crn = ceiling(row_num);
    frn = floor(row_num);

    if crn = frn and frn = row_num then
        calc_result = new_array[row_num];
    else
        calc_result = (crn - row_num) * new_array[frn]
        + (row_num - frn) * new_array[crn];
    end if;

    RETURN calc_result;
END;
$$
LANGUAGE 'plpgsql' IMMUTABLE;


CREATE FUNCTION _final_median(anyarray) RETURNS float8 AS $$
WITH q AS
(
    SELECT val
    FROM unnest($1) val
    WHERE VAL IS NOT NULL
    ORDER BY 1
),
cnt AS
(
    SELECT COUNT(*) AS c FROM q
)
SELECT AVG(val)::float8
FROM
(
    SELECT val FROM q
    LIMIT  2 - MOD((SELECT c FROM cnt), 2)
    OFFSET GREATEST(CEIL((SELECT c FROM cnt) / 2.0) - 1,0)
) q2;
$$ LANGUAGE sql IMMUTABLE;

CREATE AGGREGATE median(anyelement) (
    SFUNC=array_append,
    STYPE=anyarray,
    FINALFUNC=_final_median,
    INITCOND='{}'
);

-- http://www.jamiebegin.com/rounding-datetimes-and-timestamps-in-postgresql/
CREATE OR REPLACE FUNCTION round_timestamp(
ts timestamptz
,round_secs int
) RETURNS timestamptz AS $$
DECLARE
_mystamp timestamp;
_round_secs decimal;
BEGIN

_round_secs := round_secs::decimal;

_mystamp := timestamptz 'epoch'
       + ROUND((EXTRACT(EPOCH FROM ts))::int / _round_secs) * _round_secs
       * INTERVAL '1 second';

RETURN _mystamp;

END; $$ LANGUAGE plpgsql IMMUTABLE;
