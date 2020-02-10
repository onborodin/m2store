import React, { Component, Fragment } from 'react'
import autobind from 'autobind-decorator'

export class Pager extends Component {

    @autobind
    countPages() {
        return Math.floor(this.props.total / this.props.limit) + 1
    }

    @autobind
    currentPage() {
        return Math.floor(this.props.offset / this.props.limit)
    }


    @autobind
    renderLine() {
        const countPages = this.countPages()
        const currentPage = this.currentPage()
        const seq = Array.from(Array(countPages).keys())

        return seq.map((pageN, index) => {
            const offset = pageN * Number(this.props.limit)
            const limit = Number(this.props.limit)
            const total = Number(this.props.total)

            function show() {
                const down = offset + 1
                var up = offset + limit
                if (up > total) { up = total }
                return down + "-" + up
            }

            if (pageN == 0 && pageN != currentPage ) {
                return (
                     <li key={index} className="page-item" onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>{show()}</small>
                            </div>
                     </li>
                )
            }

            if (pageN == countPages - 1 && pageN != currentPage ) {
                return (
                     <li key={index} className="page-item"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>{show()}</small>
                            </div>
                     </li>
                )
            }

            if (pageN == currentPage ) {
                return (
                     <li key={index} className="page-item active"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>{show()}</small>
                            </div>
                     </li>
                )
            }

            if (pageN > (currentPage - 3) && pageN < currentPage) {
                return (
                     <li key={index} className="page-item"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>{show()}</small>
                            </div>
                     </li>
                )
            }

            if (pageN < (currentPage + 3) && pageN > currentPage) {
                return (
                     <li key={index} className="page-item"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>{show()}</small>
                            </div>
                     </li>
                )
            }

            if (pageN == (currentPage - 3)) {
                return (
                     <li key={index} className="page-item"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>...</small>
                            </div>
                     </li>
                )
            }

            if (pageN == (currentPage + 3)) {
                return (
                     <li key={index} className="page-item"  onClick={() => this.props.callback(offset)}>
                            <div className="page-link">
                                <small>...</small>
                            </div>
                     </li>
                )
            }
        })
    }

    @autobind
    render() {
        return (
            <nav>
                <ul className="pagination pagination-sm">
                    {this.renderLine()}
                </ul>
            </nav>
        )
    }

}

export default Pager
