import React, {Component} from 'react';
import {Button, Card, Form, Input, Row, Select} from "@canonical/react-components";
import {formatError, T} from "./Utils";
import api from "./api";
import AlertBox from "./AlertBox";

const seriesOptions = [{label:'16', value:'16'}, {label:'18', value:'18'}]

class Auth extends Component {
    constructor(props) {
        super(props)
        this.state = {
            email: '',
            password: '',
            otp: '',
            store: '',
            series: '16',
        }
    }

    onChangeEmail = (e) => {
        e.preventDefault()
        this.setState({email: e.target.value})
    }
    onChangePassword = (e) => {
        e.preventDefault()
        this.setState({password: e.target.value})
    }
    onChangeOTP = (e) => {
        e.preventDefault()
        this.setState({otp: e.target.value})
    }
    onChangeStore = (e) => {
        e.preventDefault()
        this.setState({store: e.target.value})
    }
    onChangeSeries = (e) => {
        e.preventDefault()
        this.setState({series: e.target.value})
    }

    handleLogin = (e) => {
        e.preventDefault()
        api.storeLogin(this.state.email, this.state.password, this.state.otp, this.state.store, this.state.series).then(response => {
            //this.props.onClick()
            this.setState({error:''})
            window.location.href = "/"
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data), message: ''});
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
                {this.renderError()}
                <Card>
                    <Form>
                        <Input onChange={this.onChangeEmail} type="text" id="email" placeholder={T('email-help')} label={T('email')} value={this.state.email}/>
                        <Input onChange={this.onChangePassword} type="password" id="password" placeholder={T('password-help')} label={T('password')} value={this.state.password}/>
                        <Input onChange={this.onChangeOTP} type="text" id="otp" placeholder={T('otp-help')} label={T('otp')} value={this.state.otp}/>
                        <Input onChange={this.onChangeStore} type="text" id="store" placeholder={T('store-help')} label={T('store')} value={this.state.store}/>
                        <Select onChange={this.onChangeSeries} label={T('series')} name="series" defaultValue={this.state.series} options={seriesOptions}/>
                        <Button onClick={this.handleLogin} appearance="positive">{T('login')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default Auth;