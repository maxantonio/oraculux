import React from "react";
import {Server} from "./Server";

export class Servers extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            servers: []
        }
    }

    componentWillReceiveProps(nextProps) {
        this.setState({
            servers: nextProps.servers
        });
    }

    render() {

        var serv = this.state.servers
        return (
            <table className="table bg-white table-hover">
                <thead>
                <th>Server</th>
                <th>IsMining</th>
                <th>IsSyncing</th>
                <th>Peers</th>
                <th>Last Block</th>
                <th>Minner</th>
                <th>Total Dificulty</th>
                <th>Transactions</th>
                <th>Uncles</th>
                <th>Latency</th>
                <th>Pending?</th>
                </thead>
                <tbody>
                {serv.map((server) =>
                    <Server info={server}/>
                )}
                </tbody>
            </table>
        );
    }
}