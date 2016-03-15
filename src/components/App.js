import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import appSelector from '../selectors'
import Bookmarks from './Bookmarks'
import TagSearch from './TagSearch'

import {
  addBookmark,
  pushTag,
  popTag,
  updateTagInput,
  onTagClick
} from '../actions'


const mapDispatch = (dispatch) => {
  return {
    onBookmarkAdd: x => dispatch(addBookmark(x)),
    onPushTag: x => dispatch(pushTag(x)),
    onPopTag: x => dispatch(popTag()),
    onTagInput: x => dispatch(updateTagInput(x)),
    onTagClick: x => dispatch(onTagClick(x))
  }
}

class App extends Component {

  render () {

    const {
      bookmarks,
      user,
      onBookmarkAdd,
      onPushTag,
      onPopTag,
      onTagInput,
      onTagClick,
      searchTags,
      searchInput
    } = this.props

    return ( 
      <div>
        <a href="https://github.com/login/oauth/authorize?client_id=1f5639c087a1c263a3ce">Hello </a>
      <TagSearch
        pushTag = {onPushTag}
        popTag = {onPopTag}
        onInput = {onTagInput}
        onTagClick = {onTagClick}
        tags = {searchTags}
        value = {searchInput}
      />
      <Bookmarks bookmarks = {bookmarks}/>       

      </div>
    )
  }
}

export default connect(appSelector, mapDispatch)(App)
