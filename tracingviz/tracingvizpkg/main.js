const { createElement, useState } = React;
const render = ReactDOM.render;
const jsx = htm.bind(createElement);

function EventsTable(props) {
    return jsx`
        <table className="events-table">
            <tbody>
                <tr>
                    <td className="event-name">${props.name}</td>
                    <td className="event-since-start-of-run">${props.timeSinceStartOfRun}</td>
                    <td title="${props.percentageThrough}%" className="event-percentage-through-cell">
                        <span style="border-left: 1px solid blue; margin-left:${props.percentageThrough}%"></span>
                    </td>
                </tr>
            </tbody>
        </table>
    `;
}

function Row(props) {
    return jsx`
        <tr style="display: none;">
            <td colspan="5">
                ${props.showEvents && jsx`<${EventsTable} />`}
            </td>
        </tr>
    `;
}

function nanosToMs(nanos) {
    return Math.round(nanos / (1000 * 1000));
}

function TraceRow(props) {
    const {trace} = props;

    const [expanded, setExpanded] = useState(false);

    const traceDurationNanos = trace.endTimeNanos - trace.startTimeNanos;
    return jsx`
        <tr>
            <td>
                <button type="button" className="expand-events-row" onClick=${() => setExpanded(!expanded)}>${expanded ? "Collapse" : "Expand"}</button>
            </td>
            <td>${trace.name}</td>
            <td>${nanosToMs(traceDurationNanos)}ms</td>
        </tr>
    
    
        ${expanded && jsx`<tr>
            <td colSpan="5">
                <table className="events-table">
                    <thead>
                        <tr>
                            <th>Name</th>
                            <th>Duration</th>
                            <th>Timeline</th>
                        </tr>
                    </thead>
                    <tbody>
                        ${trace.spans.map((span, idx) => {
                            const spanDuration = span.endTimeNanos - span.startTimeNanos;
                            const percentageThroughFromStart = (span.startTimeNanos - trace.startTimeNanos) * 100 / traceDurationNanos;
                            const style = {
                                borderLeft: '1px solid blue',
                                marginLeft: `${percentageThroughFromStart}%`,
                                minWidth: '1px',
                                width: `${(spanDuration * 100) / (traceDurationNanos)}%`,
                                height: '20px', // TODO: better measurement
                                backgroundColor: 'blue',
                            };
                            return jsx`
                            <tr key=${idx}>
                                <td className="event-name">${span.name}</td>
                                <td className="event-since-start-of-run">${nanosToMs(spanDuration)}ms</td>
                                <td title="${span.name}: ${nanosToMs(spanDuration)}ms" className="event-percentage-through-cell">
                                    <div style="${style}"></div>
                                </td>
                            </tr>
                        `})}
                    </tbody>
                </table>
            </td>
        </tr>`}
    `;
}

function Table(props) {
    return jsx`
        <table cellSpacing="10px" className="traces-table">
            <thead>
                <tr>
                    <th></th>
                    <th>Trace</th>
                    <th>Duration</th>
                </tr>
            </thead>
            <tbody>
            ${props.traces.map((trace, idx) => jsx`
                <${TraceRow} trace=${trace} key=${idx} />
            `)}
            </tbody>
        </table>
    `;
}

render(jsx`<${Table} traces=${window.traces}/>`, document.getElementById("App"));
