cat bookmark.json | curl -H "Authorization: 12345" -X POST --data-binary @- http://localhost:555/api/bookmarks
