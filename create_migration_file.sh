#!/usr/bin/env bash
set -e

CMD_GOOSE='go run github.com/pressly/goose/v3/cmd/goose@latest'
DIR='./internal/persistence/sqlite/migrations'
DIALECT='sqlite3'
DBNAME='test.db'

${CMD_GOOSE} \
  -dir "${DIR}" \
  ${DIALECT} \
  ${DBNAME} \
  create \
  ${1?migration file name} \
  sql
