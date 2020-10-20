import React from 'react';
import {Modal} from "@canonical/react-components";
import {T} from "./Utils";

function SnapDelete(props) {
    return (
        <Modal close={props.onCancel} title={T('confirm-delete')}>
            <p>
                {T('confirm-delete-snap-message') + props.message}
            </p>
            <hr />
            <div className="u-align--right">
                <button onClick={props.onCancel} className="u-no-margin--bottom">
                    {T('cancel')}
                </button>
                <button className="p-button--negative u-no-margin--bottom" onClick={props.onConfirm} >
                    {T('delete')}
                </button>
            </div>
        </Modal>
    );
}

export default SnapDelete;