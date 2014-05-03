#kanejaku

Experiment to track metrics "à la graphite" but lighter

## Getting started

### Dependencies

* Postgresql >= 9.3.

### Run it on Heroku

```
$ hk create kanejaku
$ hk set BUILDPACK_URL=https://github.com/kr/heroku-buildpack-go.git#go1.2
$ hk addon-add heroku-postgresql
$ hk psql -c '\i schema.sql'
$ git push heroku master
```

Application connects itself to the url pointed by env variable `DATABASE_URL`.