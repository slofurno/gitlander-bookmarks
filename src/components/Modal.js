import React, { Component, PropTypes } from 'react'

export default class JobModal extends Component {
  closeModal (e) {
    e.stopPropagation()
		this.props.cancelModal()
  }

  stopP (e) {
    e.stopPropagation()
  }


  handleClick(e) {
		e.preventDefault()
    const name = this.refs.name.value
    const city = this.refs.city.value
    const post = this.refs.post.value
    const status = this.refs.status.value
		const {id} = this.props.selectedJob
    this.props.updateJob({
			id, name, city, post, status
		})
  }

  render () {

    const {updateJob, deleteJob, selectedJob, addStatus} = this.props

    if (selectedJob.id < 0) {
      return (<div></div>)
    }

    let background = {
      position: "fixed",
      display: "relative",
      top: "0",
      bottom: "0",
      left: "0",
      right: "0",
      backgroundColor: "RGBA(0,0,0,.4)"
    }

    let modalContainer = {
      position: "absolute",
      margin: "auto",
      top: "0",
      bottom: "0",
      left: "0",
      right: "0",
    }

    return (
      <div style={background} onClick={e => this.closeModal(e)}>
        <div className="modal" style={modalContainer} onClick={this.stopP}>
          <div className="card job-card flex column">
            <h3>{selectedJob.lastStatus}</h3>
            <input 
              className="job-display" 
              type="text"
              ref="name" 
              defaultValue={selectedJob.name}
              placeholder="company name" 
            />
            <input 
              className="job-display" 
              type="text" 
              ref="city" 
              defaultValue={selectedJob.city}
              placeholder="city" 
            />
            <textarea 
              rows="8"
              className="job-display flex grow" 
              ref="post"
              defaultValue={selectedJob.post}
            ></textarea>
            <input 
              className="job-display" 
              type="text"
              ref="status" 
              defaultValue={selectedJob.status}
            />
            <a className="link" href="#" onClick={(e) => this.handleClick(e)}>Save</a>
            <a className="link" href="#" onClick={(e) => addStatus({job:selectedJob.id, status:"tevs"})}>Move >>></a>
            <a className="link" href="#" onClick={(e) => deleteJob(selectedJob)}>Delete</a>
          </div>
        </div>
      </div>
    )
  }

}
