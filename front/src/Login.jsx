import React, { Component, Fragment  } from 'react'
import autobind from 'autobind-decorator'
import axios from 'axios'
import { debounce, throttle, trim } from 'lodash'

import { store } from './main'

@autobind
export class Login extends React.Component {

    constructor(props) {
        super(props)

        this.state = {
            username: "",
            password: "",
            alertMessage: "",
            countAttempt: 0
        }
    }

    checkLogin() {
        axios.post('/api/v2/user/login', {
                username: this.state.username,
                password: this.state.password
        }).then((res) => {
            if (res.data.error != null) {
                console.log("login: ", res.data)
                if (!res.data.error) {
                    this.setState({alertMessage: ""})
                    store.login(res.data.result)

                    this.props.history.push('/')
                } else {
                    this.setState({
                        alertMessage: "Wrong password or username",
                        countAttempt: this.state.countAttempt + 1
                    })
                }
            }
        })
    }

    onSubmit(event) {
        event.preventDefault()
        //let checkLogin = debounce(this.checkLogin, 1000)
        this.checkLogin()
    }

    onChangeUsername(event) {
        event.preventDefault()
        const newUsername = trim(event.target.value)
        store.username = newUsername
        this.setState({
            username: newUsername
        })
    }

    onChangePassword(event) {
        event.preventDefault()
        const newPassword = trim(event.target.value)
        this.setState({
            password: newPassword
        })
    }

    showAttempts() {
        if (this.state.countAttempt > 0) {
            return (<span>Attempt {this.state.countAttempt}</span>)
        }
    }

    render() {
        return (
            <Fragment>
                <nav className="navbar navbar-expand-sm sticky-top navbar-dark bg-dark">
                    <button className="btn m-0 p-0">
                        <div className="mr-3">
                            <i className="fa fa-bars fa-lg"></i>
                        </div>
                    </button>

                    <div className="navbar-brand">
                        <i className="fab fa-old-republic fa-lg"></i> S2 Login
                    </div>

                </nav>


                <div className="container-fluid">
                    <div className="row justify-content-center">
                        <div className="col-8 col-sm-6 col-md-4 border p-4 mt-sm-5 ml-3 mr-3">

                            <form onSubmit={this.onSubmit}>
                                <div className="form-group">
                                    <label htmlFor="username">Username:</label>
                                    <input id="username" className="form-control" type="text"  value={this.state.username} onChange={this.onChangeUsername} />
                                </div>

                                <div className="form-group">
                                    <label htmlFor="password">Password:</label>
                                    <input id="password" className="form-control" type="password"  value={this.state.password} onChange={this.onChangePassword} />
                                </div>

                                <div className="text-center mb-3">
                                        {this.state.alertMessage}
                                </div>
                                <div className="text-center mb-3">
                                        {this.showAttempts()}
                                </div>

                                <div className="text-center">
                                    <button onClick={this.onSubmit}  className="btn btn-primary btn-sm">Submit</button>
                                </div>

                            </form>

                        </div>
                    </div>
                </div>
            </Fragment>
        )
    }
}

export default Login
