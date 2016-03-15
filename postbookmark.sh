#!/bin/sh

cat json | curl -H "Authorization: $1" -X POST --data-binary @- http://localhost:555/api/bookmarks
