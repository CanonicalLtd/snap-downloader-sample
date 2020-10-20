import React, {Component} from 'react';
import {Button, Card, Form, Input, Row, Select} from "@canonical/react-components";
import {T} from "./Utils";

const arches = [
    {label: 'amd64', value: 'amd64'},
    {label: 'armhf', value: 'armhf'},
    {label: 'arm64', value: 'arm64'},
    ]

class SnapAdd extends Component {
    render() {
        return (
            <Row>
                <Card>
                    <Form>
                        <Input onChange={this.props.onChangeName} type="text" id="name" placeholder={T('snap-name-help')} label={T('snap-name')} value={this.props.name}/>
                        <Select onChange={this.props.onChangeArch} label={T('arch')} name="arch" defaultValue={this.props.arch} options={arches}/>
                        <Button onClick={this.props.onClick} appearance="positive">{T('add')}</Button>
                        <Button onClick={this.props.onCancel} appearance="neutral">{T('cancel')}</Button>
                    </Form>
                </Card>
            </Row>
        );
    }
}

export default SnapAdd;