import React, { PropTypes } from 'react'
import Bookmark from './Bookmark'

const Bookmarks = ({bookmarks}) => {

  let list = bookmarks.map((x,i) => <Bookmark bookmark={x} key={i}/>)

  return (
    <div> 
      {list}
    </div>
  )
}

export default Bookmarks

