import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';
import App from './App';
import registerServiceWorker from './registerServiceWorker';

ReactDOM.render(<App />, document.getElementById('root'));
registerServiceWorker();
var perfData = window.performance.timing;
var pageLoadTime = perfData.loadEventEnd - perfData.navigationStart;
var connectTime = perfData.responseEnd - perfData.requestStart;
var renderTime = perfData.domComplete - perfData.domLoading;
// console.log(pageLoadTime);
// console.log(connectTime);
// console.log(renderTime);
// console.log(perfData);



