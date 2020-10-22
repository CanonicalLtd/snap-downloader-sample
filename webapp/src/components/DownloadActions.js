import React from 'react';
import {T} from "./Utils";
import {Link} from "@canonical/react-components";


function DownloadActions(props) {
    return (
        <div>
            <Link href={'/v1/downloads/' + props.name + '/' + props.snap} title={T("download-snap")}>
                <img className="action" src="/static/images/download.svg" alt={T("download-snap")} />
            </Link>
            <Link href={'/v1/downloads/' + props.name + '/' + props.assertion} title={T("download-assertion")}>
                <img className="action" src="/static/images/download.svg" alt={T("download-assertion")} />
            </Link>
        </div>
    );
}

export default DownloadActions;