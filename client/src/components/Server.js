import React from "react";

export class Servers extends React.Component {
    constructor(props) {
        super(props);
        console.log("creando el server")
        this.state = {
            servers: []
        }
    }

    getDefaultProps() {
        console.log("si se llama default props")
    }

    componentWillReceiveProps(nextProps) {
        console.log("SE LLAMA WILLRECEIVprops ")
        // Can use shallowEquals() helper here to avoid comparing every prop
        // if (this.props.servers !== nextProps.servers) {
        this.setState({
            servers: nextProps.servers
        });
        // }
    }

    render() {
        var serv = this.state.servers
        return (
            <table className="full-table">
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
                        <td>{server.Latency}</td>
                        <td>{server.Penging}</td>

                    </tr>
                )}

                </tbody>
            </table>
        );
    }
}