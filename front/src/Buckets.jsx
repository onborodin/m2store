import React, { Component, Fragment } from 'react'
import { Link } from 'react-router-dom'

import autobind from 'autobind-decorator'
import axios from 'axios'
import Humanize from 'humanize-plus'

import { store } from './main'
import { Layout } from './Layout'
import { checkLogin } from './main'
import { Pager } from './Pager'

class Buckets extends Component {

    constructor(props) {
        super(props)

        this.state = {
            buckets: [],
            offset: 0,
            limit: 5,
            total: 0,
            pattern: "*",
            alertMessage: ""
        }
    }

    @autobind
    listBuckets() {
        axios.post('/api/v1/bucket/pagelist', {
                limit: this.state.limit,
                offset: this.state.offset,
                pattern: this.state.pattern
        }).then((res) => {
            if (res.data.error != null) {
                console.log("list buckets response: ", res.data)
                if (!res.data.error) {
                    this.setState({
                        //buckets: res.data.result
                        buckets: res.data.result.buckets,
                        total: res.data.result.total,
                        offset: res.data.result.offset,
                        limit: res.data.result.limit
                    })
                } else {
                    this.setState({
                        alertMessage: "Backend error"
                    })
                }
            }
        }).catch((err) => {
            this.setState({
                alertMessage: "Communication error"
            })
        })
    }

    @autobind
    onChangeLimit(event) {
        const newLimit = Number(event.target.value)
        var newOffset = Math.floor(this.state.offset / newLimit) * newLimit
        this.setState({ limit: newLimit, offset: newOffset }, () => { this.listBuckets() })
    }

    @autobind
    changeOffset(newOffset) {
        this.setState({ offset: newOffset }, () => { this.listBuckets() })
    }

    @autobind
    onChangePattern(event) {
        event.preventDefault()
        const newPattern = event.target.value
        this.setState({ pattern: newPattern }, () => { this.listBuckets() })
    }

    @autobind
    renderTable() {
        return this.state.buckets.map((item, index) => {
            const theItem = item
            return (
                <tr key={index}>
                    <td>{index + 1}</td>
                    <td><Link to={"/files/" + theItem.name}> /{theItem.name} </Link></td>
                    <td>{Humanize.fileSize(theItem.size)}</td>
                </tr>
            )
        })
    }

    render() {
        return (
            <Fragment>
                <Layout>
                    <h5><i className="fas fa-folder"></i> Buckets
                            <i className="fas fa-sync fa-xs ml-2" onClick={this.listBuckets}></i>
                    </h5>
                    <div className="row mb-2">
                        <div className="col">
                            <div className="input-group input-group-sm flex-nowrap">

                                <div className="input-group-prepend">
                                    <div className="input-group-text">{this.state.total}</div>
                                </div>

                                <input type="text" className="form-control" id="button-pattern"  value={this.state.pattern} onChange={this.onChangePattern} />

                                <div className="input-group-append">
                                    <select className="custom-select" id="page-limit" value={this.state.limit} onChange={this.onChangeLimit}>
                                        <option value="5">5</option>
                                        <option value="10">10</option>
                                        <option value="25">25</option>
                                        <option value="50">50</option>
                                    </select>
                                </div>

                            </div>
                        </div>

                    </div>


                    <table className="table table-striped table-hover table-sm">

                        <thead className="thead-light">
                            <tr>
                                <th>#</th>
                                <th>name</th>
                                <th>size</th>
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
        checkLogin()
        this.listBuckets()
    }

    componentDidMount() {
        checkLogin()
        this.setState({
                limit: store.bucketLimit,
                pattern: store.bucketPattern
            },
            () => { this.listBuckets() }
        )
    }

    componentWillUnmount() {
        store.bucketLimit = this.state.limit
        store.bucketPattern = this.state.pattern
    }



}

export default Buckets
