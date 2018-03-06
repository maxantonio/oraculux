import React from "react";
import {Timer} from "./Timer";
import {BarChart, Bar, Tooltip} from 'recharts';

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
            UpTime: 0,
            Pending: {Pending: 0, Queued: 0},
            propagation: {Block: 0, Date: new Date()},
            History: [{Time: 10, Fill: this.getTimeFill(10)}],
            timeProp: 0,
        }
    }


    getTimeFill(time) {
        if (time <= 10) {
            return "#5fc46a";
        } else if (time > 10 && time <= 20) {
            return "#ffd75a";
        } else if (time > 20 && time < 30) {
            return "#ff8812";
        } else {
            return "#eb4b4b";
        }
    }



    componentWillReceiveProps(nextProps) {

        this.setState(nextProps.info);
        if (this.state.propagation.Block !== nextProps.propagation.Block) {//si cambia el bloque general
            this.setState({propagation: nextProps.propagation});//actualizamos la informacion de la propagacion
        }
        if (nextProps.info.BlockNumber > this.state.BlockNumber) {//si el bloque que entra nuevo es
            if (nextProps.info.BlockNumber == nextProps.propagation.Block) {
                var propagation = new Date() - nextProps.propagation.Date;
                var newHistory = new Array()
                var start = 0;
                if (this.state.History.length > 39) {
                    start = 1
                }
                for (var i = start; i < this.state.History.length; i++) {
                    newHistory.push(this.state.History[i]);
                }
                let obj = {Time: propagation, Fill: this.getTimeFill(propagation)}
                newHistory.push(obj);
                this.setState({timeProp: propagation, History: newHistory});
            }
        }


    }
    render() {

        const color_sinc = ((this.state.BlockNumber === this.state.propagation.Block) ? "tc-green" : '') +
            ((this.state.BlockNumber < this.state.propagation.Block) ? "tc-orange" : '') +
            ((this.state.BlockNumber > this.state.propagation.Block) ? "tc-blue" : '');
        return (

            <tr className={color_sinc}>
                        <td>{this.state.Server}</td>
                        <td>{this.state.IsMining ? "X" : "No"}</td>
                        {this.state.Sincing != null ? <td>{this.state.Sincing.IsSyncing ? "X" : "No"}</td> :
                            <td>No</td>}
                        <td>{this.state.Peers}</td>
                        <td>{this.state.BlockNumber}</td>
                <Timer best_block={this.state.BlockNumber} div={false}/>
                        {this.state.Block != null && <td>{this.state.Block.Miner}</td>}
                        {this.state.Block != null && <td>{this.state.Block.TotalDifficulty}</td>}
                        <td>{this.state.Transactions}</td>
                        {this.state.Block != null && <td>{this.state.Block.Uncles.length}</td>}
                <td>{this.state.timeProp}</td>
                <td>
                    <BarChart width={60} height={20} data={this.state.History} bind>
                        <Bar dataKey='Time' fillKey='Fill'/>
                        <Tooltip/>
                    </BarChart>
                </td>
                        <td>{this.state.Latency}</td>
                <td>{this.state.UpTime}%</td>
                        {this.state.Pending != null &&
                        <td>{(this.state.Pending.Pending + this.state.Pending.Queued).toString()}</td>
                        }
                <td><BarChart/></td>
                    </tr>
        );
    }
}