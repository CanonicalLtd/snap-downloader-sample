import React, {Component} from 'react';
import {formatError, T} from "./Utils";
import {Row} from "@canonical/react-components";
import Auth from './Auth';


class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
        }
    }

    handleCreate = () => {
        this.props.onLogin()
    }

    render() {
        return (
            <div>
                <br />
                <section>
                    <Row>
                        {!this.props.macaroon.modified ?
                            <Auth onClick={this.handleCreate} /> :
                            <h4>Store ID: {this.props.macaroon.store}</h4>
                        }
                    </Row>
                </section>
            </div>
        );
    }
}

export default Home;