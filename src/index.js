import React from 'react'
import { createStore, applyMiddleware } from 'redux'
import { Provider } from 'react-redux'
import thunk from 'redux-thunk'
import { render } from 'react-dom'

import rootReducer from './reducers'
import App from './components/App'

import {
  tryLogin,
  postCode
} from './actions'

let store = createStore(
  rootReducer,
  applyMiddleware(thunk)
)

let unsubscribe = store.subscribe(() =>
  console.log(store.getState())
)

var storedid = localStorage.getItem("loginid");
var storedtoken =  localStorage.getItem("token");
var accountid = localStorage.getItem("accountid");

function mapQueryString(s){
      return s.split('&')
  .map(function(kvp){
        return kvp.split('=');
  }).reduce(function(sum,current){
        var key = current[0];
        var value = current[1];
        sum[key] = value;
        return sum;
  },{});
}

var qs = mapQueryString(location.search.substr(1));
var code = qs["code"];

history.pushState('','','/');

if (storedtoken) {
  store.dispatch(tryLogin(accountid))
} else if (code) {
  store.dispatch(postCode(code)) 

} else {
  console.log("login")
}

render(
  <Provider store={store}>
    <App/>
  </Provider>,
  document.getElementById('root')
)
