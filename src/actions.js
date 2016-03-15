import request from './request'
export const ADD_BOOKMARK = 'ADD_BOOKMARK'
export const SET_TOKEN = 'SET_TOKEN'

export const PUSH_TAG = 'PUSH_TAG'
export const POP_TAG = 'POP_TAG'
export const TAG_CLICK = 'TAG_CLICK'
export const UPDATE_TAG_INPUT = 'UPDATE_TAG_INPUT'


const origin = location.origin

export function addBookmark (bookmark) {
  return {
    type: ADD_BOOKMARK,
    bookmark
  }
}

export function setToken (token) {
  return {
    type: SET_TOKEN,
    token
  }
}

export function pushTag (tag) {
  return {
    type: PUSH_TAG,
    tag
  }
}

export function popTag () {
  return {
    type: POP_TAG
  }
}

export function updateTagInput (value) {
  console.log("tag val", value)
  return {
    type: UPDATE_TAG_INPUT,
    value
  }
}

export function onTagClick (tag) {
  return {
    type: TAG_CLICK,
    tag
  }
}

function log (err) {
  console.log(err)
}

function listen (ws) {
  ws.onmessage = function(e) {

  } 
}

export function postBookmark (bookmark, token) {
  return function (dispatch) {
    return request({
      url:`/api/bookmarks`,
      method: "POST",
      headers: {"Authorization": token}
    })
  }
}

export function postCode (code) {

  return function (dispatch) {
    return request({
      url: `/api/user?code=${code}`,
      method: "POST"
    })
    .then(x => {
      let result = JSON.parse(x)
      let {user, token, userid} = result

      localStorage.setItem("loginid", user);
      localStorage.setItem("token", token);
      localStorage.setItem("accountid", userid);

      dispatch(tryLogin(user, token))
    })
    .catch(log)
  }

}

export function tryLogin (storedid) {
  return function (dispatch) {
    return new Promise((resolve, reject) => {
      let wsuri = origin.replace(/^http/, 'ws')
      let ws = new WebSocket(`${wsuri}/ws?user=${storedid}`)
      ws.onerror = function(e) {
        reject(e) 
      }
      ws.onopen = function(e) {
        resolve(ws)
      }
    })
    .then(ws => {
      ws.onmessage = (e) => {
        let x = JSON.parse(e.data)  
        dispatch(addBookmark(x))
      }

      dispatch(setToken(storedid))
    })
    .catch(log)
  } 
}

