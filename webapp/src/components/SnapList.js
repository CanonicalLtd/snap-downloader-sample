import React, {Component} from 'react';
import api from "./api";
import {formatError, T} from "./Utils";
import {MainTable, Row, Button, Notification} from "@canonical/react-components";
import SnapAdd from "./SnapAdd";
import SnapDelete from "./SnapDelete";
import SnapActions from "./SnapActions";

class SnapList extends Component {
    constructor(props) {
        super(props)
        this.state = {
            delete: {},
            showAdd: false,
            showDelete: false,
            name: '',
            arch: 'amd64',
        }
    }

    handleAddClick = (e) => {
        e.preventDefault()
        this.setState({name:'', arch:'amd64', showAdd: true})
    }

    handleNameChange = (e) => {
        e.preventDefault()
        this.setState({name: e.target.value})
    }

    handleArchChange = (e) => {
        e.preventDefault()
        this.setState({arch: e.target.value})
    }

    handleSnapCreate = (e) => {
        e.preventDefault()
        api.snapCreate(this.state.name, this.state.arch).then(response => {
            this.props.onCreate()
            this.setState({error:'', showAdd: false, repo:''})
        })
        .catch(e => {
            console.log(formatError(e.response.data))
            this.setState({error: formatError(e.response.data)});
        })
    }

    handleCancelClick = (e) => {
        e.preventDefault()
        this.setState({showAdd: false, showDelete: false, error: ''})
    }

    handleDeleteClick = (e) => {
        e.preventDefault()
        let id = e.target.getAttribute('data-key')

        let rr = this.props.records.filter(r => {
            return r.id === id
        })

        let del = this.state.delete
        del.id = id
        del.name = rr[0].name
        this.setState({showDelete: true, delete: del})
    }

    handleDeleteDo = (e) => {
        e.preventDefault()

        this.props.onDelete(this.state.delete.id)
        this.setState({showDelete: false, delete: {}})
    }

    render() {
        let data = this.props.records.map(r => {
            return {
                columns: [
                    {content: r.name, role: 'rowheader'},
                    {content: r.arch},
                    {content: r.created},
                    {content: <SnapActions id={r.id} onDelete={this.handleDeleteClick} />}
                ],
            }
        });

        return (
            <section>
                <Row>
                    <div>
                        <h3 className="u-float-left">{T('snap-list')}</h3>
                        <Button onClick={this.handleAddClick} className="u-float-right">
                            {T('add-snap')}
                        </Button>
                    </div>
                    <p>{T('home-desc')}</p>
                    {
                        this.state.error ?
                            <Row>
                                <Notification type="negative" status={T('error') + ':'}>
                                    {this.state.error}
                                </Notification>
                            </Row>
                            : ''
                    }
                    {this.state.showAdd ?
                        <SnapAdd onClick={this.handleSnapCreate} onCancel={this.handleCancelClick}
                                 onChangeName={this.handleNameChange} onChangeArch={this.handleArchChange}
                                 name={this.state.name} arch={this.state.arch} />
                        :
                        ''
                    }
                    {this.state.showDelete ?
                        <SnapDelete onCancel={this.handleCancelClick} onConfirm={this.handleDeleteDo} message={this.state.delete.name} />
                        : ''
                    }
                    <MainTable headers={[
                        {
                            content: T('name'),
                            className: "col-medium"
                        }, {
                            content: T('arch'),
                        }, {
                            content: T('created'),
                        }, {
                            content: T('actions'),
                        }
                    ]} rows={data}/>
                </Row>
            </section>
        );
    }
}

export default SnapList;