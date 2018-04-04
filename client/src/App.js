import React, {Component} from 'react';
import { BarChart,Bar,Tooltip } from 'recharts';
import 'react-datamaps';
import './App.css';
import {Servers} from "./components/Servers";
import {Timer} from "./components/Timer";
import Datamap from "./components/datamap";


class App extends Component {


    constructor() {
        super();
        this.historyProp = [{server: ""}];
        this.histogranProp = [
            0
        ]
        this.uncles = [{name: 0, cont: 0, fill: "#ccc"}];
        this.transactions = [
            {name: 0, transactions: 0, fill: "#8884d8"}
        ];
        this.last_block = 0; //tiempo demorado en obtener el bloque
        this.state = {
            transactions:this.transactions,

            times: [{name: 0, cont: 0, fill: "#ccc"}], //se llena en addTime
            avgtime: 15, //se calcula en addTime
            best_block: 0,//el ultimo bloque segun los nodos conectados al stat
            lastknow_block: 0,//el ultimo bloque cuando esta sincronizandose
            gas_price:0,
            miners:["1","2"],
            dificulty:0,
            dificulties:[],
            totalDificulty:0,
            gasLimit:0,
            gasLimitList:[],
            gasUsed:[],
            uncles:[{name:0,cont:0,fill:"#ccc"}],
            uncle_val:0,
            uncle_50:0,
            hash_rate:0,
            peers: 0,
            propagation: {Block: 0, Date: new Date()},

            servers: []
        };

         this.last_Block();
        // let urlbase = "ws://" + document.location.host + "/ws";
        let urlbase = "ws://localhost:8080/ws";// para desarrollo
        let ws = new WebSocket(urlbase);
        let self = this;
        ws.onerror = function (error) {

        };
        ws.onmessage = function(event) {
            let response = JSON.parse(event.data);
            switch(response.info_type){
                case "Syncing":
                    self.setSyncing(response.data);
                    break;
                case "Server":
                    self.setServersD(response.data);
                    self.setFullInfo(response.data);
                    break;
                case "FullInfo":
                    self.setFullInfo(response.data);
                    break;
                 default:
                     console.log("No definido");
                     // self.setStatus(response.info_type,response.data,response.block);
                    break;
            }
        }
        ws.onclose = this.closeF;
    }

    addPropagationServer = function (name, prop) {
        this.histogranProp.push(prop);
        var propagation = this.histogranProp;
        console.log(propagation);
        var data = window.d3.layout.histogram()
            .frequency(false)
            .range([0, 10000])
            .bins(10)
            (propagation);

        var freqCum = 0;
        var histogram = data.map(function (val) {
            freqCum += val.length;
            var cumPercent = (freqCum / Math.max(1, propagation.length));

            return {
                x: val.x,
                dx: val.dx,
                y: val.y,
                frequency: val.length,
                cumulative: freqCum,
                cumpercent: cumPercent
            };
        });
        console.log("HISTOGRAMMM SSSS");

        console.log(histogram);
        this.setState({
            histogram: histogram
        });
    }
    closeF(event) {
        console.log("CONEXION CERRADA");
    }
    setFullInfo(data) {
        if (data.BlockNumber > this.state.best_block) {
            this.addTime(this.last_block, data.BlockNumber);
            this.addTransactions(data.Transactions, data.BlockNumber);
            if (data.Block != null) {

                if (data.Block.Number > this.state.best_block) {

                    var newminers = [];
                    newminers.push(this.state.miners[1]);
                    newminers.push(data.Block.Miner);
                    this.setGasUsed(data.Block.GasUsed);
                    this.setGasLimit(data.Block.GasLimit);
                    this.setDificulties(data.Block.Difficulty);
                    this.addUncle(data.Block.Uncles.length, data.Block.Number);
                    var propagation = {Block: data.Block.Number, Date: new Date()};
                    this.setState({
                        miners: newminers,
                        dificulty: data.Block.Difficulty,
                        totalDificulty: data.Block.TotalDifficulty,
                        gasLimit: data.Block.GasLimit,
                        propagation: propagation
                    })
                }
            }
            this.setState({
                best_block: data.BlockNumber,
                gas_price: data.GasPrice,
                hash_rate: data.HashRate,
                peers: data.Peers,
                // last_block: 0
            });
            this.last_Block = 0;
            if (data.Sincing.IsSyncing) {
                this.setState({
                    lastknow_block: data.Sincing.HighestBlock,
                });
            }
        }
    }
    setServersD(data) {
        data.key = data.Server;

        var serversOld = this.state.servers;
        var servers = [];
        var finded = false;
        for (var i = 0; i < serversOld.length; i++) {
            var server = serversOld[i];
            if (server.Server === data.Server) {
                server = data;
                finded = true;
            }
            servers.push(server);
        }
        if (!finded) {
            servers.push(data);
        }
        this.setState({
            servers: servers
            }
        );
    }
    setSyncing(data){
        if(this.state.best_block !== data.CurrentBlock) {
            this.addTime(this.last_block, data.CurrentBlock);
            this.setState({
                best_block: data.CurrentBlock,
                lastknow_block:data.HighestBlock,
            });
            this.last_Block = 0;
        }
    }

    compare(a, b) {
        if (a.Time < b.Time)
            return -1;
        if (a.Time > b.Time)
            return 1;
        return 0;
    }
    setDificulties(gass) {
        var temp = [];
        var inicio = 0;
        if (this.state.dificulties.length > 49)
            inicio = 1;
        for (var i = inicio; i < this.state.dificulties.length; i++) {
            temp.push(this.state.dificulties[i]);
        }
        temp.push({cont: gass, fill: this.getTimeFill()});
        this.setState({
            dificulties: temp
        });
    }
    setGasUsed(gass){
        var temp = [];
        var inicio = 0;
        if(this.state.gasUsed.length>49)
            inicio = 1;
        for (var i = inicio; i < this.state.gasUsed.length; i++) {
            temp.push(this.state.gasUsed[i]);
        }
        temp.push({cont:gass,fill:this.getTimeFill()});
        this.setState({
                gasUsed:temp
            });
    }
    setGasLimit(gass){
        var temp = [];
        var inicio = 0;
        if(this.state.gasLimitList.length>49)
            inicio = 1;
        for (var i = inicio; i < this.state.gasLimitList.length; i++) {
            temp.push(this.state.gasLimitList[i]);
        }
        temp.push({cont:gass,fill:this.getTimeFill()});
        this.setState({
            gasLimitList:temp
        });
    }

    //utiles para obtener colores
    getTimeFill(time){
        if(time<=10) {
            return "#5fc46a";
        }else if(time>10&&time<=20){
            return "#ffd75a";
        }else if(time>20&&time<30){
            return "#ff8812";
        }else{
            return "#eb4b4b";
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

    //agregando valores alos graficos
    addTime(time,block){
//cambio
        var nuevo = {name: block, cont: time, fill: this.getTimeFill(time)};
        var temp = [];
        var total =0;
        for (var j = 0; j < this.state.times.length; j++) {
            total+=this.state.times[j].cont;
            temp.push(this.state.times[j]);
        }
        var avg = total/temp.length;
        if( this.state.times.length>49)
            temp.shift();
        temp.push(nuevo);
        this.setState({
                times: temp,
                avgtime:avg
            });
        this.last_block = 0
    }
    addUncle(count,block){
        if (this.uncles.length === 0 || block > this.uncles[this.uncles.length - 1].name) {
            var nuevo = {name: block, cont: count, fill: this.getTransFill(count)};
            this.uncles.push(nuevo);
            if(this.uncles.length>50)
                this.uncles.shift();
            var temp = [];
            for (var j = 0; j < this.uncles.length; j++) {
                temp.push(this.uncles[j]);
            }
            var result = 0;
            for(var i in temp) {
                result += temp[i].cont;
            }
            this.setState({ uncles:temp,uncle_50:result,uncle_val:nuevo.cont});
        }
    }
    addTransactions(count,block){
        //AQUI ES NECESARIO
        if (block > this.transactions[this.transactions.length - 1].name) {
            var trans = [];
            if(this.transactions.length>40)
                this.transactions.shift();
            for (var j = 0; j < this.transactions.length; j++) {
                trans.push(this.transactions[j]);
            }
            var nuevo = {name: block, transactions: count, fill: this.getTransFill(count)};
            trans.push(nuevo);
            this.transactions = trans;
            this.setState({transactions: trans});
        }
        // this.transactions = count;
    }

    last_Block() {
        setInterval(() => {
            this.last_block = this.last_block + 1;
        }, 1000);
    }
    numberWithCommas = (x) => {
        return x.toString().replace(/\B(?=(\d{3})+(?!\d))/g, ",");
    };
    render() {
        const last_block = ((this.last_block <= 12) ? 'tc-green' : '') +
            ((this.last_block >= 13 && this.last_block <= 19) ? 'tc-yellow' : '') +
            ((this.last_block >= 20 && this.last_block <= 29) ? 'tc-orange' : '') +
            ((this.last_block >= 30) ? 'tc-red' : '');


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
                    <Timer best_block={this.state.best_block} div={true}/>
                    <div className="box">
                        <div>
                            <div className="pull-left icon tc-yellow"><i className="fa fa-clock-o fa-4x"></i></div>
                            <div className="info">
                                <span className="title">avg block time</span>
                                <span className="value tc-yellow">{this.state.avgtime}</span>
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
                                <span className="value tc-red dificulty">{this.state.dificulty} H</span>
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
                            <span className="pull-right tc-blue mr-2">{this.state.gasLimit} gas</span>
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
                    <div className="rowInside">
                        <div className="row">
                            <div className="box">
                                <span className="title2">Transactions</span>
                                <BarChart width={280} height={80} data={this.state.transactions} bind >
                                    <Bar dataKey='transactions' fillKey='fill'/>
                                    <Tooltip/>
                                </BarChart>
                            </div>
                            <div className="box">
                                <span className="title2">Uncles</span>
                                <BarChart width={280} height={80} data={this.state.uncles} bind>
                                    <Bar dataKey='cont' fillKey='fill'  />
                                    <Tooltip/>
                                </BarChart>
                            </div>
                            <div className="box">
                                <span className="title2">Block Time</span>
                                <BarChart width={280} height={80} data={this.state.times} bind>
                                    <Bar dataKey='cont' fillKey='fill'  />
                                    <Tooltip/>
                                </BarChart>
                            </div>
                            <div className="box">
                                <span className="title2">DIFFICULTY</span>
                                <BarChart width={280} height={80} data={this.state.dificulties} bind>
                                    <Bar dataKey='cont' fillKey='fill'  />
                                    <Tooltip/>
                                </BarChart>
                            </div>
                        </div>
                        <div className="row">
                            <div className="box">
                                <span className="title2">GAS SPENDING</span>
                                <BarChart width={280} height={80} data={this.state.gasUsed} bind >
                                    <Bar dataKey='cont' fillKey='fill'/>
                                    <Tooltip/>
                                </BarChart>
                            </div>
                            <div className="box">
                                <span className="title2">GAS LIMIT</span>
                                <BarChart width={280} height={80} data={this.state.gasLimitList} bind>
                                    <Bar dataKey='cont' fillKey='fill'  />
                                    <Tooltip/>
                                </BarChart>
                            </div>

                            <div className="box" >
                                <span className="title2">last blocks miners</span>
                                <div className="small-title-miner ng-binding minners">{this.state.miners[0]}</div>
                                <div blocks="14">
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                </div>
                                <div className="small-title-miner ng-binding minners">{this.state.miners[1]}</div>
                                <div blocks="14">
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                    <div className="block bg-info"></div>
                                </div>
                            </div>
                            <div className="box">
                                <span className="title2">Histogram</span>
                                <BarChart width={280} height={80} data={this.state.histogram} bind>
                                    <Bar dataKey='y'/>
                                    <Tooltip/>
                                </BarChart>
                            </div>
                        </div>
                    </div>
                    <div className="mapa-content">
                        <Datamap/>
                    </div>
                </div>
                <div className="row">
                    <Servers servers={this.state.servers} propagation={this.state.propagation}
                             handlePropagation={this.addPropagationServer.bind(this)}/>
                </div>
            </div>

        );
    }
}

export default App;
