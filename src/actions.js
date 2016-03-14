import request from './request'
export const ADD_BOOKMARK = 'ADD_BOOKMARK'

const origin = location.origin

export function addBookmark (bookmark) {
  return {
    type: ADD_BOOKMARK,
    bookmark
  }
}

function log (err) {
  console.log(err)
}

function listen (ws) {
  ws.onmessage = function(e) {

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
      ws.onmessage = (e) => log(e)  
      
    })
    .catch(log)
  } 
}

