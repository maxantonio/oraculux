import React from "react";

export class Timer extends React.Component {
    constructor(props) {
        super(props);
        this.state = {
            last_block: 0,
            best_block: 0,
            div: props.div,
        }
        this.last_Block();
    }

    componentWillReceiveProps(nextProps) {
        if (nextProps.best_block > this.state.best_block) {
            this.setState({
                last_block: 0,
                best_block: nextProps.best_block,
                div: nextProps.div
            })
        }
        this.setState(nextProps.info);
    }

    last_Block() {
        setInterval(() => {
            this.setState({
                last_block: this.state.last_block + 1
            })
        }, 1000);
    }

    render() {
        const last_block = ((this.state.last_block <= 12) ? 'tc-green' : '') +
            ((this.state.last_block >= 13 && this.state.last_block <= 19) ? 'tc-yellow' : '') +
            ((this.state.last_block >= 20 && this.state.last_block <= 29) ? 'tc-orange' : '') +
            ((this.state.last_block >= 30) ? 'tc-red' : '');
        if (this.state.div) {
            return (
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
            );
        } else {
            return (
                <td className={last_block}>
                    <span className={"value " + last_block}>{this.state.last_block}s ago</span>
                </td>
            );
        }
    }
}