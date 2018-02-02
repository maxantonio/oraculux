import React, {Component} from 'react';
import './App.css';

class App extends Component {
    render() {
        return (
            <div className="container-fluid">
                <div className="row">
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-blue"><i class="fa fa-codepen fa-4x"></i></div>
                            <div className="info">
                                <span className="title">best block</span>
                                <span className="value tc-blue">#5,018,498</span></div>
                        </div>
                    </div>
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-blue"><i class="fa fa-code-fork fa-4x"></i></div>
                            <div className="info">
                                <span className="title">uncles &nbsp;
                                <span class="small">(current / last 50)</span></span>
                                <span className="value tc-blue">0/7</span></div>
                        </div>
                    </div>
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-red"><i class="fa fa-hourglass-o fa-4x"></i></div>
                            <div className="info">
                                <span className="title">last block</span>
                                <span className="value tc-red">23s ago</span></div>
                        </div>
                    </div>
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-yellow"><i class="fa fa-clock-o fa-4x"></i></div>
                            <div className="info">
                                <span className="title">avg block time</span>
                                <span className="value tc-yellow">14.79s</span></div>
                        </div>
                    </div>
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-orange"><i class="fa fa-fire fa-4x"></i></div>
                            <div className="info">
                                <span className="title">avg network hashrate</span>
                                <span className="value tc-orange">174 TH/s</span></div>
                        </div>
                    </div>
                    <div className="col-md-2 box">
                        <div>
                            <div className="pull-left icon tc-red"><i class="fa fa-puzzle-piece fa-4x"></i></div>
                            <div className="info">
                                <span className="title">difficulty</span>
                                <span className="value tc-red">0.00 H</span></div>
                        </div>
                    </div>
                </div>
            </div>
        );
    }
}

export default App;
