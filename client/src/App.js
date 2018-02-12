import React, {Component} from 'react';
import { LineChart, Line } from 'recharts';
import './App.css';

class App extends Component {
    constructor() {
        super();
        this.transactions = 0;
        this.state = {
            last_block: 0,
            best_block: 0,
            lastknow_block:0,
            gas_price:0,
            uncles:0,
            hash_rate:0,
            peers:0
        };

         this.last_Block();
        // this.best_Block();

        var ws = new WebSocket("ws://" + document.location.host + "/ws");
         console.log('ws', ws)
        // setTimeout(() => {
        //     ws.send('eth2');
        // }, 5000)
        //
        var self = this;
        ws.onmessage = function(event) {
             var response = JSON.parse(event.data);

            switch(response.info_type){
                case "Syncing":
                    self.setSyncing(response.data);
                    break;
                 default:
                    self.setStatus(response.info_type,response.data);
                    break;
            }
            console.log("respondio el server el dato de "+response.info_type, response.data);
        }
    }

    setSyncing(data){
        console.log("ESCRIBIENDO CON SYNCING")
        if(this.state.best_block !== data.CurrentBlock) {
            this.setState({
                best_block: data.CurrentBlock,
                lastknow_block:data.HighestBlock,
                last_block: 0
            });
        }
    }
    setStatus(type,data){
        switch(type){
            case "Uncles":
                this.setState({ uncles:data});
                break;
            case "GasPrice":
                this.setState({ gas_price:data});
                break;
            case "Peers":
                this.setState({ peers:data});
                break;
            case "Hashrate":
                this.setState({ hash_rate:data});
                break;
            case "Transactions":
                this.addTransactions(data);
                break;
            default:
                console.log(type+" no esta definido todavia",data)
        }
    }
    addTransactions(count){
        //AQUI ES NECESARIO
        this.transactions = count;
    }
    best_Block() {
        setInterval(() => {
            this.setState({
                best_block: Math.trunc(Math.random() * (5019999 - 5019000) + 5019000)
            });
        }, 5000);
    }

    last_Block() {
        setInterval(() => {
           this.setState({
                last_block: this.state.last_block + 1
            });
        }, 1000);
    }

    numberWithCommas = (x) => {
        return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    }


    render() {
        const last_block = ((this.state.last_block <= 12) ? 'tc-green' : '') +
            ((this.state.last_block >= 13 && this.state.last_block <= 19) ? 'tc-yellow' : '') +
            ((this.state.last_block >= 20 && this.state.last_block <= 29) ? 'tc-orange' : '') +
            ((this.state.last_block >= 30) ? 'tc-red' : '');

        return (
            <div className="container-fluid">
                <div className="row">
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-blue"><i className="fa fa-codepen fa-4x"></i></div>
                            <div className="info">
                                <span className="title">best block <span className="small">(last block:#{this.numberWithCommas(this.state.lastknow_block)} )</span></span>
                                <span className="value tc-blue">#{this.numberWithCommas(this.state.best_block)}</span>
                            </div>
                        </div>
                    </div>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-blue"><i className="fa fa-code-fork fa-4x"></i></div>
                            <div className="info">
                                <span className="title">uncles &nbsp;
                                    <span className="small">(current block / last 50)</span></span>
                                <span className="value tc-blue">{this.state.uncles}</span>
                            </div>
                        </div>
                    </div>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-red">
                                <i className={"fa fa-hourglass-o fa-4x " + last_block}></i></div>
                            <div className="info">
                                <span className="title">last block</span>
                                <span className={"value " + last_block}>{this.state.last_block}s ago</span>
                            </div>
                        </div>
                    </div>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-yellow"><i className="fa fa-clock-o fa-4x"></i></div>
                            <div className="info">
                                <span className="title">avg block time</span>
                                <span className="value tc-yellow">14.79s</span>
                            </div>
                        </div>
                    </div>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-orange"><i className="fa fa-fire fa-4x"></i></div>
                            <div className="info">
                                <span className="title">network hashrate</span>
                                <span className="value tc-orange">{this.state.hash_rate} TH/s</span>
                            </div>
                        </div>
                    </div>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-red"><i className="fa fa-puzzle-piece fa-4x"></i></div>
                            <div className="info">
                                <span className="title">difficulty</span>
                                <span className="value tc-red">0.00 H</span>
                            </div>
                        </div>
                    </div>
                </div>
                <div className="row">
                    <div className="box-small">
                        <div>
                            <i className="fa fa-laptop ml-2 tc-green"></i>
                            <span className="title ml-3">Peers</span>
                            <span className="pull-right tc-green mr-2">{this.state.peers}</span>
                        </div>
                    </div>
                    <div className="box-small">
                        <div>
                            <i className="fa fa-tag ml-2 tc-blue fa-flip-horizontal"></i>
                            <span className="title ml-3">gas price</span>
                            <span className="pull-right tc-blue mr-2">{this.state.gas_price}</span>
                        </div>
                    </div>
                    <div className="box-small">
                        <div>
                            <i className="fa fa-tag ml-2 tc-blue fa-flip-horizontal"></i>
                            <span className="title ml-3">gas limit</span>
                            <span className="pull-right tc-blue mr-2">8000029 gas</span>
                        </div>
                    </div>
                    <div className="box-small">
                        <div>
                            <i className="fa fa-tag ml-2 tc-red fa-flip-horizontal"></i>
                            <span className="title ml-3">page latency</span>
                            <span className="pull-right tc-red mr-2">7401 ms</span>
                        </div>
                    </div>
                    <div className="box-small">
                        <div>
                            <i className="fa fa-tag ml-2 tc-green fa-flip-horizontal"></i>
                            <span className="title ml-3">uptime</span>
                            <span className="pull-right tc-green mr-2">100%</span>
                        </div>
                    </div>
                    <div className="box-small">
                    </div>
                </div>
                <div className="row">
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-red">
                                <i className={"fa fa-hourglass-o fa-4x "}></i></div>
                            <div className="info">
                                <span className="title">Transactions</span>
                                <span className="value">{this.transactions}</span>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default App;
