import React, { Component, Fragment } from 'react'

import Clock from './Clock'
import Menu from './Menu'
import { store } from './main'

export class Layout extends Component {

    render() {
        return(
            <Fragment>
                <nav className="navbar navbar-expand-sm sticky-top navbar-dark bg-dark">

                    <button className="btn m-0 p-0" data-target="#menu" data-toggle="collapse" role="button">
                        <div className="d-block d-sm-block d-md-none mr-3">
                            <i className="fa fa-bars fa-lg"></i>
                        </div>
                    </button>

                    <div className="btn m-0 p-0">
                        <div className="d-none d-md-block mr-3">
                            <i className="fa fa-bars fa-lg"></i>
                        </div>
                    </div>

                    <div className="navbar-brand">
                        <i className="fab fa-old-republic fa-lg"></i> G2
                    </div>

                    <ul className="nav justify-content-end ml-auto mr-3">
                        <li className="nav-item">
                            <Clock/>
                        </li>
                        <li className="nav-item ml-3">
                            [{store.username}]
                        </li>
                    </ul>

                </nav>

                <div className="container-fluid">

                    <div className="row">

                        <div id="menu" className="overlay fade collapse col-5 col-sm-3 col-md-2 col-lg-2 col-xl-2 d-md-none d-lg-none d-xl-none bg-dark">
                            <div className="sticky-top">
                                <div className="px-0 py-0 my-0 mx-0">
                                    <Menu />
                                </div>
                            </div>
                        </div>

                        <div className="col-5 col-sm-3 col-md-2 col-lg-2 col-xl-2 d-none d-md-block d-lg-block d-xl-block bg-dark">
                            <div className="sticky-top sticky-offset">
                                <div className="px-0 py-0 my-0 mx-0">
                                    <Menu />
                                </div>
                            </div>
                        </div>


                        <div className="col col-11 col-sm-10 col-md-9 col-lg-8 mt-4 mx-auto">

                            {this.props.children}

                        </div>
                    </div>
                </div>

                <div className="container">
                    <div className="row">
                        <div className="col my-3">
                            <hr className="justify-content-sm-center" />
                            <div className="text-center">
                                <small>made by <a href="http://wiki.unix7.org">oleg borodin</a></small>
                            </div>
                        </div>
                    </div>
                </div>

            </Fragment>
        )
    }

}

export default Layout
