import React, { Fragment, Component } from 'react'
import $ from 'jquery'
import autobind from 'autobind-decorator'
import axios from 'axios'

import { UserUpdate } from './UserUpdate'
import { UserDelete } from './UserDelete'

export class UserOption extends Component {

    //constructor(props) {
    //    super(props)
    //}

    @autobind
    showForm() {
        $("#option-user-" + this.props.user.username).modal('show')
    }

    @autobind
    hideForm() {
        $("#option-user-" + this.props.user.username).modal('hide')
    }

    render() {
        return (
            <Fragment>

                <a onClick={this.showForm}>{this.props.children}</a>

                <div className="modal fade bd-example-modal-sm" id={"option-user-" + this.props.user.username}
                                                                    tabIndex="-1" role="dialog"  ref="form">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">

                                <div className="modal-header">
                                    <h5 className="modal-title">Database [{this.props.user.username}] option</h5>
                                    <button type="button" className="close"
                                                                onClick={this.hideForm}>&times;</button>
                                </div>

                                <div className="modal-body">
                                            <UserUpdate user={this.props.user}
                                                            updateCallback={this.props.updateCallback}
                                                            hideCallback={this.hideForm} />
                                            <div className="dropdown-divider"></div>
                                            <UserDelete user={this.props.user}
                                                            updateCallback={this.props.updateCallback}
                                                            hideCallback={this.hideForm} />

                                </div>

                                <div className="modal-footer">
                                    <button type="button" className="btn btn-sm btn-secondary"
                                                            onClick={this.hideForm}>Close</button>
                                </div>

                        </div>
                    </div>
                </div>
            </Fragment>
        )
    }
}

export default UserOption
