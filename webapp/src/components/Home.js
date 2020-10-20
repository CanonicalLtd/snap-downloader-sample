import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import {Row, Notification} from "@canonical/react-components";
import SnapList from "./SnapList";


class Home extends Component {
    constructor(props) {
        super(props)
        this.state = {
            snaps: [{id:'abc', name:'test-snap', arch:'amd64'}],
            delete: {},
        }
    }

    componentDidMount() {
        this.getSnaps()
    }

    getSnaps() {
        api.snapsList().then(response => {
            this.setState({snaps: response.data.records})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    handleSnapCreateClick = () => {
        this.getSnaps()
    }

    handleSnapDelete = (snapId) => {
        api.snapDelete(snapId).then(response => {
            this.getSnaps()
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    render() {
        return (
            <div>
                <Row>
                    <h5>Store ID: {this.props.macaroon['Snap-Device-Store']} (logged-in: {this.props.macaroon['Modified']})</h5>
                </Row>
                {
                    this.state.error ?
                        <Row>
                            <Notification type="negative" status={T('error') + ':'}>
                                {this.state.error}
                            </Notification>
                        </Row>
                        : ''
                }
                <section>
                    <Row>
                        <SnapList records={this.state.snaps} onCreate={this.handleSnapCreateClick} onDelete={this.handleDeleteClick} onDelete={this.handleSnapDelete} />
                    </Row>
                </section>
            </div>
        );
    }
}

export default Home;