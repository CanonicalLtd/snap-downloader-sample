import React from 'react';
import {T} from "./Utils";
import {MainTable, Row} from "@canonical/react-components";
import DownloadActions from "./DownloadActions";

const DownloadList = (props) => {
    let data = props.records.map(r => {
        return {
            columns: [
                {content: r.name, role: 'rowheader'},
                {content: r.arch},
                {content: r.revision},
                {content: <DownloadActions name={r.name} snap={r.filename} assertion={r.assertion} />},
            ],
        }
    });

    return (
        <Row>
            <div>
                <h3 className="u-float-left">{T('download-list')}</h3>
            </div>

            {data.length === 0 ? <p>{T('no-downloads')}</p> :
            <MainTable headers={[
                {
                    content: T('name'),
                    className: "col-medium"
                }, {
                    content: T('arch'),
                }, {
                    content: T('revision'),
                }, {
                    content: T('actions'),
                }
            ]} rows={data}/>}
        </Row>
    );
};

export default DownloadList;