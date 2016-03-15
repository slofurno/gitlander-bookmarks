import React, { PropTypes } from 'react'

const Bookmark = ({bookmark}) => {
  
  let tagline = bookmark.Tags.reduce((a,c) => a + ", " + c)

  return (
    <div className="card">
      <ul>
        <li>{bookmark.Description}</li>
        <li>{bookmark.Url}</li>
        <li>{tagline}</li>
      </ul> 
    </div>
  )
}

export default Bookmark
