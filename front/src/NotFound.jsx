import React, { Fragment, Component } from 'react'

import { Layout } from './Layout'
import { checkLogin } from './main'

export class NotFound extends Component {

    render() {
        return (
            <Fragment>
                <Layout>
                    <div className="row">
                        <h4>404 Not found</h4>
                    </div>
                </Layout>
            </Fragment>
        )
    }

    componentDidMount() {
        checkLogin()
    }

}

export default NotFound
