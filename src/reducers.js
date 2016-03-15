import { combineReducers } from 'redux'
import {
  ADD_BOOKMARK,
  SET_TOKEN,

  PUSH_TAG,
  POP_TAG,
  UPDATE_TAG_INPUT,
  TAG_CLICK
  
} from './actions'


function bookmarks (state = [], action) {
  switch (action.type) {
  case ADD_BOOKMARK:
    return state.concat([action.bookmark])
  default:
    return state
  }
}

function searchInput (state = "", action) {
  switch (action.type) {
  case UPDATE_TAG_INPUT:
    return action.value
  case PUSH_TAG:
    return ""
  default:
    return state
  }
}

function searchTags (state = {tags:[], value:""}, action) {
  const {tags, value} = state
  let {tag} = action
  switch (action.type) {
  case UPDATE_TAG_INPUT:
    return {
      value: action.value,
      tags
    }
  case PUSH_TAG:
    return {
      value: "",
      tags: tags.indexOf(value) === -1 ? tags.concat([value]) : tags 
    }
  case POP_TAG:
    return {
      value,
      tags: tags.slice(0, -1) 
    }
  case TAG_CLICK:
    return {
      value,
      tags: tags.filter(x => x !== tag)
    }
  default:
    return state 
  }
}

function user (state = {}, action) {
  switch (action.type) {
  case SET_TOKEN:
    return Object.assign({}, state, {token:action.token})
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
  modalInput,
  user
})

export default rootReducer
