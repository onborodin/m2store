import React, { Component, Fragment } from 'react'
import autobind from 'autobind-decorator'
import axios from 'axios'
import dayjs from 'dayjs'

import { clone } from 'lodash'

import { Layout } from './Layout'
import { checkLogin } from './main'

import { UserCreate } from './UserCreate'
import { UserOption } from './UserOption'


import { Pager } from './Pager'

export class Users extends Component {

    constructor(props) {
        super(props)

        this.state = {
            users: [],
            offset: 0,
            limit: 5,
            total: 0,
            alertMessage: ""
        }
    }

    @autobind
    listUsers() {
        axios.post('/api/v2/user/list', {
                limit: this.state.limit,
                offset: this.state.offset
        }).then((res) => {
            if (res.data.error != null) {
                let users = []
                if (res.data.result.users != null) {
                    users = res.data.result.users
                }

                if (!res.data.error) {
                    this.setState({
                        users: users,
                        total: res.data.result.total,
                        offset: res.data.result.offset,
                        limit: res.data.result.limit
                    })
                } else {
                    this.setState({
                        alertMessage: "Backend error. " + res.data.message
                    })
                }
            }
        }).catch((err) => {
            this.setState({
                alertMessage: "Communication error"
            })
        })
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
    changeOffset(newOffset) {
        this.setState({ offset: newOffset }, () => { this.listUsers() })
    }

    @autobind
    onChangeLimit(event) {
        const newLimit = Number(event.target.value)
        var newOffset = Math.floor(this.state.offset / newLimit) * newLimit
        this.setState({ limit: newLimit, offset: newOffset }, () => { this.listUsers() })
    }


    @autobind
    renderTable() {
        return this.state.users.map((user, index) => {
            const { id, username, isadmin } = user
            const theUser = user
            return (
                <tr key={id}>
                    <td>{index + this.state.offset + 1}</td>
                    <td><UserOption user={theUser} updateCallback={this.listUsers}>{username}</UserOption></td>
                    <td>{this.renderIsAdmin(isadmin)}</td>
                </tr>
            )
        })
    }

    @autobind
    renderIsAdmin(isadmin) {
        if (isadmin) {
            return <i className="fas fa-user-plus"></i>;
        } else {
            return <i className="fas fa-microchip"></i>;
        }
    }

    render() {
        return (
            <Fragment>
                <Layout>
                    <div className="container-fluid">
                        <div className="row mb-3">
                                <h5>
                                    <i className="fas fa-users"></i>
                                    <span> Users </span>
                                    <a onClick={this.listUsers}><i className="fas fa-sync fa-sm"></i></a>
                                    <small> {dayjs().format('hh:mm:ss')}</small>
                                </h5>
                                <div className="ml-auto">
                                    <UserCreate updateCallback={this.listUsers}/>
                                </div>
                        </div>
                    </div>

                    <div className="d-inline-flex mb-1">
                        <select className="custom-select" id="page-limit" value={this.state.limit} onChange={this.onChangeLimit}>
                            <option value="5">5</option>
                            <option value="10">10</option>
                            <option value="25">25</option>
                            <option value="50">50</option>
                        </select>
                    </div>

                    <table className="table table-striped table-hover table-sm">

                        <thead className="thead-light">
                            <tr>
                                <th>#</th>
                                <th>login name</th>
                                <th>mode</th>
                            </tr>
                        </thead>
                        <tbody>
                            {this.renderTable()}
                        </tbody>
                    </table>

                    <Pager total={this.state.total} limit={this.state.limit} offset={this.state.offset} callback={this.changeOffset} />

                </Layout>
            </Fragment>
        )
    }

    componentDidMount() {
        checkLogin("admin")
        this.listUsers()
    }
}

export default Users
