import React, { Fragment, Component } from 'react'
import $ from 'jquery'
import { trim, ltrim } from 'validator'
import autobind from 'autobind-decorator'
import axios from 'axios'

export class UserCreate extends React.Component {

    constructor(props) {
        super(props)
        this.initState = {
            username: "",
            password: "",
            isadmin: false,
            usernameIsValid: false,
            passwordIsValid: false,
            formIsValid: false,
            usernameMessage: "",
            passwordMessage: "",
            alertMessage: ""
        }
        this.state = this.initState

        this.minUsernameLength = 4
        this.minPasswordLength = 4
    }

    createUser() {
        axios.post('/api/v2/user/create', {
                username: this.state.username,
                password: this.state.password,
                isadmin: this.state.isadmin
        }).then((res) => {
            if (res.data.error != null) {
                if (!res.data.error) {
                    this.hideForm()
                    this.clearForm()
                    this.props.updateCallback()
                } else {
                    this.setState({ alertMessage: "Backend error. " +  res.data.message })
                }
            }
        }).catch((err) => {
                    this.setState({ alertMessage: "Communication error" })
        })
    }

    showForm() {
        $("#create-user").modal('show')
    }

    hideForm() {
        $("#create-user").modal('hide')
    }

    clearForm() {
        this.setState(this.initState)
    }

    @autobind
    onSubmit(event) {
        event.preventDefault()
        if (!this.state.formIsValid) {
            return
        }
        this.createUser()
    }

    validateUsername() {
        if (this.state.username.length > this.minUsernameLength) {
            this.setState({
                    usernameIsValid: true,
                    usernameMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    usernameIsValid: false,
                    usernameMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validatePassword() {
        if (this.state.password.length > this.minPasswordLength) {
            this.setState({
                    passwordIsValid: true,
                    passwordMessage: ""
                },
                () => { this.validateForm() }
            )
        } else {
            this.setState({
                    passwordIsValid: false,
                    passwordMessage: "The field must be filled"
                },
                () => { this.validateForm() }
            )
        }
    }

    validateForm() {
        if (this.state.usernameIsValid && this.state.passwordIsValid) {
            this.setState({ formIsValid: true })
        } else {
            this.setState({ formIsValid: false })
        }
    }


    @autobind
    onChangeUsername(event) {
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                username: newValue
            },
            () => { this.validateUsername() }
        )
    }

    @autobind
    onChangePassword(event) {
        event.preventDefault()
        const newValue = trim(event.target.value)
        this.setState({
                password: newValue
            },
            () => { this.validatePassword() }
        )
    }

    @autobind
    onIsAdmin(event) {
        this.setState({ isadmin: event.target.checked })
    }

    showAlert() {
        if (this.state.alertMessage != "") {
            return (
                <div className="alert alert-warning border mx-4" role="alert">
                  <div className="text-center">{this.state.alertMessage}</div>
                </div>
            )
        }
    }

    @autobind
    makeId(base) {
        return 'user-create-' + base
    }

    render() {
        return (
            <React.Fragment>

                <a onClick={this.showForm}>
                    <i className="fas fa-plus fa-lg"></i>
                </a>

                <div className="modal fade" id="create-user" tabIndex="-1" role="dialog"  ref="form">
                    <div className="modal-dialog" role="document">
                        <div className="modal-content">

                            <form acceptCharset="UTF-8" onSubmit={this.onSubmit} ref="form">

                                <div className="modal-header">
                                    <h5 className="modal-title">Create user</h5>
                                    <button type="button" className="close" onClick={this.hideForm}>
                                        <span>&times;</span>
                                    </button>
                                </div>

                                <div className="modal-body">

                                    <div className="form-group">
                                        <label htmlFor="username">Username:</label>
                                        <input type="text" className="form-control" id="username" value={this.state.username}  onChange={this.onChangeUsername}/>
                                        <small className="form-text text-muted">{this.state.usernameMessage}</small>
                                    </div>

                                    <div className="form-group">
                                        <label htmlFor="password">Password:</label>
                                        <input type="password" className="form-control" id="password" onChange={this.onChangePassword}/>
                                        <small className="form-text text-muted">{this.state.passwordMessage}</small>
                                    </div>

                                    <div className="form-group form-check">
                                        <input id={this.makeId("isadmin")} className="form-check-input"
                                                    type="checkbox" onChange={this.onIsAdmin} defaultChecked={this.state.isadmin} />
                                        <label className="form-check-label" htmlFor={this.makeId("isadmin")}> As admin</label>
                                    </div>

                                </div>

                                {this.showAlert()}

                                <div className="modal-footer">
                                    <button type="button" className="btn btn-sm btn-secondary" onClick={this.hideForm}>Close</button>
                                    <button type="submit" className={this.state.formIsValid ? "btn btn-sm btn-primary" : "btn btn-sm btn-primary disabled"} onClick={this.onSubmit}>Create</button>
                                </div>

                            </form>

                        </div>
                    </div>
                </div>
            </React.Fragment>
        )
    }
}

export default UserCreate
