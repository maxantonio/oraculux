import React, {Component} from 'react';
import { BarChart,Bar } from 'recharts';
import './App.css';

class App extends Component {

    constructor() {
        super();
        this.uncles = [{block:0,cont:4,fill:"#ccc"}]
        this.transactions = [
            {block: 0, uv: 4, fill: "#8884d8"},
            {block: 0, uv: 3, fill: "#cc3f2c"},
            {block: 0, uv: 5, fill: "#8884d8"},
            {block: 0, uv: 2, fill: "#8884d8"},
            {block: 0, uv: 5, fill: "#8884d8"},
            {block: 0, uv: 1, fill: "#8884d8"},
            {block: 0, uv: 1, fill: "#8884d8"},
            {block: 0, uv: 1, fill: "#8884d8"},
            {block: 0, uv: 1, fill: "#8884d8"},
            {block: 0, uv: 1, fill: "#8884d8"},
        ];
        this.state = {
            transactions:[{block: 0, uv: 11, fill: "#8884d8"}],
            last_block: 0,
            best_block: 0,
            lastknow_block:0,
            gas_price:0,
            uncles:[{block:0,cont:0,fill:"#ccc"}],
            uncle_val:0,
            uncle_50:0,
            hash_rate:0,
            peers:0
        };

         this.last_Block();
        // this.best_Block();
//" + document.location.host + "
        var ws = new WebSocket("ws://localhost:6060/ws");
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
                    self.setStatus(response.info_type,response.data,response.block);
                    break;
            }
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
    setStatus(type,data,block){
        switch(type){
            case "Uncles":
               this.addUncle(data,block);
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
                this.addTransactions(data,block);
                break;
            default:
                console.log(type+" no esta definido todavia",data)
        }
    }
    getTransFill(count){
        if(count<5) {
            return "#8884d8"
        }else if(count>5&&count<10){
            return "#0084cc"
        }else{
            return "#008400"
        }
    }

    addUncle(count,block){

        if(this.uncles.length===0||block!==this.uncles[this.uncles.length-1].block) {
            var nuevo = {block: block, cont: count, fill: this.getTransFill(count)}
            this.uncles.push(nuevo);
            if(this.uncles.length>50)
                this.uncles.shift();
            var temp = [];
            for (var j = 0; j < this.uncles.length; j++) {
                temp.push(this.uncles[j]);
            }
            var result = 0
            for(var i in temp) {
                console.log(i);
                result += temp[i].cont;
            }
            console.log("RESULT DEL REDUXE",result);
            this.setState({ uncles:temp,uncle_50:result,uncle_val:nuevo.cont});
        }
    }
    addTransactions(count,block){
        //AQUI ES NECESARIO
        if(block!==this.transactions[this.transactions.length-1].block) {
            var trans = [];
            for (var j = 1; j < this.transactions.length; j++) {
                trans.push(this.transactions[j]);
            }
            var nuevo = {block: block, uv: count, fill: this.getTransFill(count)}
            trans.push(nuevo);
            this.transactions = trans
            this.setState({transactions: trans});
        }
        // this.transactions = count;
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
                                <span className="value tc-blue">{this.state.uncle_val}/{this.state.uncle_50}</span>
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
                        <span className="title2">Transactions</span>
                        <BarChart width={280} height={80} data={this.state.transactions} bind>
                            <Bar dataKey='uv' fillKey='fill'/>
                        </BarChart>
                    </div>
                    <div className="box">
                        <span className="title2">Uncles</span>
                        <BarChart width={280} height={80} data={this.state.uncles} bind>
                            <Bar dataKey='cont' fillKey='fill'/>
                        </BarChart>
                    </div>
                </div>

            </div>
        );
    }
}

export default App;
