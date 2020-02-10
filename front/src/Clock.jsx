import React, { Component, Fragment } from 'react'

export class Clock extends Component {

    constructor(props) {
        super(props)
        this.state = { date: new Date() }
    }

    tick() {
        this.setState({
                date: new Date()
        })
    }

    componentDidMount() {
        this.timer = setInterval(
            () => this.tick(),
            1000
        )
    }

    componentWillUnmount() {
        clearInterval(this.timer)
    }

    render() {
        return (
            <Fragment>
                { this.state.date.toLocaleTimeString() }
            </Fragment>
        )
    }
}

export default Clock

