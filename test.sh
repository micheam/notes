#!/usr/bin/env bash

create_book()
{
  echo
  echo ===================================================
  echo TEST: BOOK CREATION
  outfile=$(mktemp)
  status_code=`command curl -X POST \
    -d "$(jo title=hello)" \ --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED] exit 1
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
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
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
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo $prepared

  # exercise: get book
  outfile=$(mktemp)
  status_code=`command curl -X GET \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
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
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo prepared: $prepared

  # exercise: edit book
  outfile=$(mktemp)
  status_code=`command curl -X PUT \
    -d "$(jo title=world)" \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
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
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books`
  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
  fi
  prepared=`cat $outfile|jq .id | tr -d '"'`
  echo $prepared

  # exercise: delete
  outfile=$(mktemp)
  status_code=`command curl -X DELETE \
    --write-out "%{http_code}" \
    --silent --output "${outfile}" \
    --unix-socket "${HOME}/.notes/notes-localserver.sock" \
    http://localhost/books/$prepared`

  echo STATUS_CODE: ${status_code}
  if [[ $status_code -ne 200 ]]; then
    echo [FAILED]
    exit 1
  fi
  cat $outfile | jq -c .
}

create_book
list_books
get_book
replace_book
delete_book

echo
echo "==================="
echo "  ALL TEST PASSED  "
echo "==================="

exit 0
