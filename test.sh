#!/usr/bin/env bash

RED='\033[0;31m'
GREEN='\033[0;32m'
NOCOLOR='\033[0m'

success() {
  printf "${GREEN}
██████╗  █████╗ ███████╗███████╗███████╗██████╗ 
██╔══██╗██╔══██╗██╔════╝██╔════╝██╔════╝██╔══██╗
██████╔╝███████║███████╗███████╗█████╗  ██║  ██║
██╔═══╝ ██╔══██║╚════██║╚════██║██╔══╝  ██║  ██║
██║     ██║  ██║███████║███████║███████╗██████╔╝
╚═╝     ╚═╝  ╚═╝╚══════╝╚══════╝╚══════╝╚═════╝ 
${NOCOLOR}\n"
}

failed() {
  printf "${RED}
  █████▒▄▄▄       ██▓ ██▓    ▓█████ ▓█████▄ 
▓██   ▒▒████▄    ▓██▒▓██▒    ▓█   ▀ ▒██▀ ██▌
▒████ ░▒██  ▀█▄  ▒██▒▒██░    ▒███   ░██   █▌
░▓█▒  ░░██▄▄▄▄██ ░██░▒██░    ▒▓█  ▄ ░▓█▄   ▌
░▒█░    ▓█   ▓██▒░██░░██████▒░▒████▒░▒████▓ 
 ▒ ░    ▒▒   ▓▒█░░▓  ░ ▒░▓  ░░░ ▒░ ░ ▒▒▓  ▒ 
 ░       ▒   ▒▒ ░ ▒ ░░ ░ ▒  ░ ░ ░  ░ ░ ▒  ▒ 
 ░ ░     ░   ▒    ▒ ░  ░ ░      ░    ░ ░  ░ 
             ░  ░ ░      ░  ░   ░  ░   ░    
                                     ░      
${NOCOLOR}\n"
}
 
ping()
{
  local status_code=$(command curl -X GET \
    --write-out "%{http_code}" \
    --silent \
    --output /dev/null \
    --connect-timeout 0.1 \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books)

  echo "$status_code"

  if [[ "${status_code}" -ne "000" ]]; then
    return 0
  fi
  return 1
}

create_book()
{
  echo
  echo ===================================================
  echo TEST: BOOK CREATION
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

list_books()
{
  echo
  echo ===================================================
  echo TEST: LIST BOOKS
  outfile=$(mktemp)
  status_code=`command curl -X GET \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

get_book()
{
  echo
  echo ===================================================
  echo TEST: VIEW BOOK

  # prepare test book
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 201 ]]; then
    failed && exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo $prepared

  # exercise: get book
  outfile=$(mktemp)
  status_code=`command curl -X GET \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

replace_book()
{
  echo
  echo ===================================================
  echo TEST: REPLACE BOOK

  # prepare test book
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 201 ]]; then
    failed && exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo prepared: $prepared

  # exercise: edit book
  outfile=$(mktemp)
  status_code=`command curl -X PUT \
    -d "$(jo title=world)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

delete_book()
{
  echo
  echo ===================================================
  echo TEST: DELET BOOK

  # prepare test book
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 201 ]]; then
    failed && exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo $prepared

  # exercise: delete
  outfile=$(mktemp)
  status_code=`command curl -X DELETE \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

create_cont()
{
  echo
  echo ===================================================
  echo TEST: CONTENT CREATION
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello book_id=494a35c5-e010-4ed3-a904-f2c05f186e15)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/localserver.sock" \
    http://localhost/content`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 201 ]]; then
    failed && exit 1
  fi
  cat $outfile | jq -c .
}

if ! ping; then
  failed
  echo 😭 'local server not running'
  exit 1
fi

#create_book
#list_books
#get_book
#replace_book
#delete_book

create_cont

success && exit 0
