//requires
const screen = require('../screen.js');
const text = require('../text.js');
const YQL = require('yql');

//vars
let isActive = false;
let temp = false;

//config
let woeid = 12596838; //Berlin, Germany. Lookup: http://woeid.rosselliot.co.nz/
let unit = 'c'; //f: Fahrenheit, c: Celsius

//load weather
function loadWeather() {
	let query = new YQL('SELECT item.condition FROM weather.forecast WHERE woeid = ' + woeid + ' AND u = "' + unit + '"');

	query.exec(function(err, data) {
		if(data !== undefined) {
			if(data.error !== undefined) {
				console.log(data.error);
				temp = false;
			} else {
				temp = data.query.results.channel.item.condition.temp;
			}
			update();
		}
	});
}
loadWeather();

//reload weather every 5 minutes, regardless if active or not
setInterval(loadWeather, 5 * 60 * 1000);

//update
function update() {
	//is active?
	if(!isActive) {
		return;
	}

	//update screen
	screen.clear();
	let string;
	if(temp === false) {
		string = '...';
	} else {
		string = temp + '°';
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
