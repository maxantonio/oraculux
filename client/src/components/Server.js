import React from "react";

export class Server extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            Server: "",
            IsMining: "",
            Sincing: {IsSyncing: false,},
            Peers: 0,
            BlockNumber: 0,
            Block: {Uncles: [], Miner: "Not defined", TotalDifficulty: 0},
            Transactions: 0,
            Latency: 0,
            Pending: {Pending: 0, Queued: 0}
        }
    }
    componentWillReceiveProps(nextProps) {
        this.setState(nextProps.info);

    }

    render() {

        return (
                    <tr>
                        <td>{this.state.Server}</td>
                        <td>{this.state.IsMining ? "X" : "No"}</td>
                        {this.state.Sincing != null ? <td>{this.state.Sincing.IsSyncing ? "X" : "No"}</td> :
                            <td>No</td>}
                        <td>{this.state.Peers}</td>
                        <td>{this.state.BlockNumber}</td>
                        {this.state.Block != null && <td>{this.state.Block.Miner}</td>}
                        {this.state.Block != null && <td>{this.state.Block.TotalDifficulty}</td>}
                        <td>{this.state.Transactions}</td>
                        {this.state.Block != null && <td>{this.state.Block.Uncles.length}</td>}
                        <td>{this.state.Latency}</td>
                        {this.state.Pending != null &&
                        <td>{(this.state.Pending.Pending + this.state.Pending.Queued).toString()}</td>}
                    </tr>
        );
    }
}