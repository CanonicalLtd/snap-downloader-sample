import React, {Component} from 'react';
import {formatError, T} from "./Utils";
import {Row, Card, Form, Input, Button} from '@canonical/react-components'
import api from "./api";
import AlertBox from "./AlertBox";

class Settings extends Component {
    constructor(props) {
        super(props)
        this.state = {
            lastRun: '',
            interval: 300,
            error: '',
        }
    }

    componentDidMount() {
        this.getDataInterval()
        this.getDataLastRun()
    }

    getDataLastRun() {
        api.settingsLastRun().then(response => {
            this.setState({lastRun: response.data.record})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    getDataInterval() {
        api.settingsInterval().then(response => {
            this.setState({interval: response.data.record})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    onChangeInterval = (e) => {
        e.preventDefault()
        this.setState({interval: parseInt(e.target.value, 10)})
    }

    handleSetInterval = (e) => {
        e.preventDefault()
        api.settingsSetInterval(this.state.interval).then(response => {
            this.setState({error: ''});
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    renderError() {
        if (this.state.error) {
            return (
                <AlertBox message={this.state.error} />
            );
        }
    }

    render() {
        return (
            <Row>
                <h2>{T('settings')}</h2>

                {this.renderError()}
                <Card>
                    <h4>{T('watch-daemon')}</h4>
                    <p>{T('last-run')}: {this.state.lastRun}</p>

                    <Form>
                        <Input onChange={this.onChangeInterval} type="number" min="5" max="86400" id="interval" placeholder={T('interval-help')} label={T('interval')} value={this.state.interval}/>
                        <Button onClick={this.handleSetInterval} appearance="positive">{T('save')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default Settings;
