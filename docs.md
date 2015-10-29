/ws?user={userid} -H token

api/follow?follow={userid}&user={userid} -H token

POST api/bookmarks

GET api/user
-returns all userids

POST api/user
-returns new user {user,token}

user
-id

bookmark
-id
-userid
-url

subscription
-id
-userid
-subid

things
type
data

relationship
user-user
user-bookmark
