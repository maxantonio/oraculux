import React from "react";

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
                    <tr>
                        <td>{server.Server}</td>
                        <td>{server.IsMining ? "X" : "No"}</td>
                        {server.Sincing != null ? <td>{server.Sincing.IsSyncing ? "X" : "No"}</td> : <td>No</td>}
                        <td>{server.Peers}</td>
                        <td>{server.BlockNumber}</td>
                        {server.Block != null && <td>{server.Block.Miner}</td>}
                        {server.Block != null && <td>{server.Block.TotalDifficulty}</td>}
                        <td>{server.Transactions}</td>
                        {server.Block != null && <td>{server.Block.Uncles.length}</td>}
                        <td>{server.Latency}{server.Pending.Pending}</td>
                        <td>{(server.Pending.Pending + server.Pending.Queued).toString()}</td>

                    </tr>
                )}

                </tbody>
            </table>
        );
    }
}