import React from "react";

export class Servers extends React.Component {
    constructor() {
        super();
        this.state = {
            servers: [{Server: "Serverc constructor", Peers: 5}]
        }
    }

    getDefaultProps() {

    }

    render() {
        var serv = this.props.servers
        return (
            <table className="full-table">
                <thead>
                <th>Server</th>
                <th>Peers</th>
                </thead>
                <tbody>
                {serv.map((server) =>
                    <tr>
                        <td>{server.Server}</td>
                        <td>{server.Peers}</td>
                    </tr>
                )}

                </tbody>
            </table>
        );
    }
}