//requires
const screen = require('../screen.js');
const text = require('../text.js');

//vars
let isActive = false;
let price = false;
let activityTimer;

//connect to websocket
const WebSocket = require('ws')
let ws;
function connect() {
	ws = new WebSocket('ws://ws.pusherapp.com:80/app/de504dc5763aeef9ff52?client=js&version=2.2&protocol=5');
	ws.on('open', function() {
		//subscribe to channel
		pusherSend('pusher:subscribe', {
			'channel': 'live_trades',
			'auth': null,
			'channel_data': {}
		});
	});
	ws.on('message', function(message) {
		resetActivityCheck();

		let obj = JSON.parse(message);
		let data = JSON.parse(obj.data);
		if(data.price != undefined) {
			price = Math.round(data.price).toString();
			update();
		}
	});
	ws.on('close', function() {
	    // Reset price
		price = false;
		update();

		// Reconnect
		setTimeout(function() {
			connect();
		}, 5000);
	});
	ws.on('error', function(error) {
	    console.log(error);
	});
}
connect();

function resetActivityCheck() {
	if(activityTimer) {
		clearTimeout(activityTimer);
	}

	activityTimer = setTimeout(function() {
		pusherSend('pusher:ping');

		activityTimer = setTimeout(function() {
			ws.terminate();
		}, 10000);
	}, 20000);
}

function pusherSend(event, data) {
	let json = JSON.stringify({
		'event': event,
		'data': data
	});
	ws.send(json);
}

//update
function update() {
	//is active?
	if(!isActive) {
		return;
	}

	//update screen
	screen.clear();
	let string;
	if(price === false) {
		string = '...';
	} else {
		string = price;
	}
	text.addText(screen, string, 0);
	board.update(screen);
}

//start
function start() {
	isActive = true;
	update();
}
exports.start = start;

//stop
function stop() {
	isActive = false;
}
exports.stop = stop;

//set board
function setBoard(_board) {
	board = _board;
	screen.init(board.sizeX, board.sizeY);
}
exports.setBoard = setBoard;
