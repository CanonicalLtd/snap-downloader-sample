import React from 'react';
import {T} from "./Utils";
import {Link} from "@canonical/react-components";

function SnapActions(props) {
    return (
        <div>
            <Link href="" title={T("delete")} onClick={props.onDelete}>
                <img className="action" src="/static/images/delete.svg" alt={T("delete")} data-key={props.id}/>
            </Link>
        </div>
    );
}

export default SnapActions;