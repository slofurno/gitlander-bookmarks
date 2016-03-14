import React, { Component, PropTypes } from 'react'
import { connect } from 'react-redux'
import appSelector from '../selectors'

import {
  addBookmark
} from '../actions'


const mapDispatch = (dispatch) => {
  return {}
}

class App extends Component {

  render () {
    return ( <a href="https://github.com/login/oauth/authorize?client_id=1f5639c087a1c263a3ce">Hello </a> )
  }
}

export default connect(appSelector, mapDispatch)(App)
