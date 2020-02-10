import React, { Component, Fragment } from 'react'
import { Link } from 'react-router-dom'

import autobind from 'autobind-decorator'
import axios from 'axios'
import Humanize from 'humanize-plus'
import moment from 'moment'

import { store, history } from './main'
import { Layout } from './Layout'
import { checkLogin } from './main'
import { Pager } from './Pager'

class Files extends Component {

    constructor(props) {
        super(props)

        this.state = {
            files: [],
            bucket: "",
            offset: 0,
            limit: 10,
            total: 0,
            pattern: "*",
            alertMessage: ""
        }
    }

   @autobind
    listFiles() {
        axios.post('/api/v1/file/pagelist', {
                limit: this.state.limit,
                offset: this.state.offset,
                pattern: this.state.pattern,
                bucket: this.props.match.params[0]
        }).then((res) => {
            if (res.data.error != null) {
                console.log("list files response: ", res.data)
                if (!res.data.error) {
                    this.setState({
                        //files: res.data.result
                        files: res.data.result.files,
                        total: res.data.result.total,
                        offset: res.data.result.offset,
                        limit: res.data.result.limit,
                        bucket: res.data.result.bucket
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
        event.preventDefault()
        const newLimit = Number(event.target.value)
        var newOffset = Math.floor(this.state.offset / newLimit) * newLimit
        this.setState({ limit: newLimit, offset: newOffset }, () => { this.listFiles() })
    }


    @autobind
    onChangePattern(event) {
        event.preventDefault()
        const newPattern = event.target.value
        this.setState({ pattern: newPattern }, () => { this.listFiles() })
    }


    @autobind
    changeOffset(newOffset) {
        this.setState({ offset: newOffset }, () => { this.listFiles() })
    }

    @autobind
    renderTable() {
        return this.state.files.map((item, index) => {
            const theItem = item
            return (
                <tr key={index}>
                    <td>{index + 1}</td>
                    <td><Link to={"/api/v1/file/down/" + this.state.bucket + "/" + theItem.name} target="_blank" download>{theItem.name}</Link></td>
                    <td>{Humanize.fileSize(theItem.size)}</td>
                    <td>{moment(theItem.modtime).format('YYYY-MM-DD HH:MMZ')}</td>
                </tr>
            )
        })
    }

    render() {
        return (
            <Fragment>
                <Layout>

                    <h5><Link to="/"><i className="fas fa-folder-open"></i></Link> /{this.props.match.params[0]}
                            <i className="fas fa-sync fa-xs ml-3" onClick={this.listFiles}></i>
                    </h5>

                    <div className="row mb-2">
                        <div className="col">
                            <div className="input-group input-group-sm flex-nowrap">

                                <div className="input-group-prepend">
                                    <div className="input-group-text">{this.state.total}</div>
                                </div>

                                <input type="text" className="form-control" id="file-pattern"  value={this.state.pattern} onChange={this.onChangePattern} />

                                <div className="input-group-append">
                                    <select className="custom-select" id="page-limit" value={this.state.limit} onChange={this.onChangeLimit}>
                                        <option value="5">5</option>
                                        <option value="10">10</option>
                                        <option value="25">25</option>
                                        <option value="50">50</option>
                                        <option value="100">100</option>
                                        <option value="200">200</option>
                                        <option value="500">500</option>
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
                                <th>mtime</th>
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
        this.setState({
                limit: store.fileLimit,
                pattern: store.filePattern
            },
            () => { this.listFiles() }
        )
    }

    componentWillUnmount() {
        store.fileLimit = this.state.limit
        store.filePattern = this.state.pattern
    }

}

export default Files
