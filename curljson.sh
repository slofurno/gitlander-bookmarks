#!/bin/sh

cat json | curl -H "Authorization: 2799535" -X POST --data-binary @- http://localhost:555/api/bookmarks
