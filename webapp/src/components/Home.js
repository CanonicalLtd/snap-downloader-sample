import React, {Component} from 'react';
import {T} from "./Utils";
import {Row} from "@canonical/react-components";


class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
        }
    }

    render() {
        return (
            <div>
                <Row>
                    <h3>{T('home')}</h3>
                    <p>{T('home-desc')}</p>
                </Row>
                <section>
                    <Row>
                        <h5>Store ID: {this.props.macaroon['Snap-Device-Store']} ({this.props.macaroon['Modified']})</h5>
                    </Row>
                </section>
            </div>
        );
    }
}

export default Home;