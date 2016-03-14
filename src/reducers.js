import { combineReducers } from 'redux'
import {
  ADD_BOOKMARK
} from './actions'


function bookmarks (state = [], action) {

  switch (action.type) {

  default:
    return state
  }
}

function searchTags (state = [], action) {
  switch (action.type) {

    default:
      return state
  }
}

const emptyBookmark = {
  Id: "",
  Owner: "",
  Url: "",
  Description: "", 
  Tags: [],
  Summary: "",
  Time: 0
}

function modalInput (state = emptyBookmark, action) {
  switch (action.type) {

    default:
      return state
  }
}


const rootReducer = combineReducers({
  bookmarks,
  searchTags,
  modalInput
})

export default rootReducer
