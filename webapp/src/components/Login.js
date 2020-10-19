import React, {Component} from 'react';
import {T} from "./Utils";
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
                <Row>
                    <h3>{T('login-title')}</h3>
                    <p>{T('login-desc')}</p>
                </Row>
                <section>
                    <Row>
                        <Auth onClick={this.handleCreate} />
                    </Row>
                </section>
            </div>
        );
    }
}

export default Home;