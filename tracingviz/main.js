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

function Table(props) {
    return jsx`
        <table cellSpacing="10px">
            <thead>
                <tr>
                    <th></th>
                    <th>Run</th>
                    <th>Summary</th>
                    <th>Start Time</th>
                    <th>Duration</th>
                </tr>
            </thead>
            <tbody>
            ${props.spans.map((span, idx) => {
                console.log(span)
                return jsx`
                <tr>
                    <td>
                        <button type="button" className="expand-events-row">Expand</button>
                    </td>
                    <td>{{.Name}}</td>
                    <td>{{.Summary}}</td>
                    <td>{{.StartTime}}</td>
                    <td>{{.Duration}}</td>
                </tr>
                <tr>
                    <td colspan="5">
                        <table className="events-table">
                            <tbody>
                                ${span.events.map((event, idx) => {
                                    const style = {
                                        borderLeft: '1px solid blue',
                                        marginLeft: `${event.percentageThrough}%`,
                                    };
                                    return jsx`
                                    <tr>
                                        <td className="event-name">${event.name}</td>
                                        <td className="event-since-start-of-run">${event.timeSinceStartOfRun}</td>
                                        <td title="${event.percentageThrough}%" className="event-percentage-through-cell">
                                            <span style=${style}"></span>
                                        </td>
                                    </tr>
                                `})}
                            </tbody>
                        </table>
                    </td>
                </tr>
            `})}
            </tbody>
        </table>
    `;
}

function Text(props) {
    console.log('props', props)
    return jsx`<p>aaa ${props.text}</p>`;
}

function ClickCounter() {
    const [count, setCount] = useState(0);
    
    return jsx`
    <div>
        <button onClick=${() => setCount(count + 1)}>
        Clicked ${count} times
        </button>
        <${Text} text=${'count::' + count}>
    </div>
    `;
}

// render(jsx`<${ClickCounter}/>`, document.getElementById("App"));
render(jsx`<${Table} spans=${[]}/>`, document.getElementById("App"));
