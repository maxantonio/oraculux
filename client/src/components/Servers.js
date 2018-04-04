import React from "react";
import {Server} from "./Server";

export class Servers extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            servers: [],
            propagation: {}
        }
    }
    componentWillReceiveProps(nextProps) {
        this.setState({
            servers: nextProps.servers,
            propagation: nextProps.propagation
        });
    }

    render() {

        var serv = this.state.servers
        return (
            <table className="table bg-white table-hover">
                <thead>
                <tr>
                    <th>Server</th>
                    <th>IsMining</th>
                    <th>IsSyncing</th>
                    <th>Peers</th>
                    <th>Block</th>
                    <th>Last Block</th>
                    <th>Minner</th>
                    <th>Total Dificulty</th>
                    <th>Transactions</th>
                    <th>Uncles</th>
                    <th>Propagation</th>
                    <th>Propagation Chart</th>
                    <th>Latency</th>
                    <th>UpTime</th>
                    <th>Pending?</th>
                </tr>
                </thead>
                <tbody>
                {serv.map((server) =>
                    <Server key={server.key} info={server} propagation={this.state.propagation}
                            handlePropagation={this.props.handlePropagation}/>
                )}
                </tbody>
            </table>
        );
    }
}