const wheelnumbersAC = [0, 26, 3, 35, 12, 28, 7, 29, 18, 22, 9, 31, 14, 20, 1, 33, 16, 24, 5, 10, 23, 8, 30, 11, 36, 13, 27, 6, 34, 17, 25, 2, 21, 4, 19, 15, 32];
const redNumbers = {
    1: true, 3: true, 5: true, 7: true, 9: true, 12: true, 14: true, 16: true, 18: true, 19: true, 21: true, 23: true, 25: true, 27: true, 30: true, 32: true, 34: true, 36: true
};

let chipValues = {
	'white': 25,
	'red': 50,
	'blue': 100,
	'green': 500,
	'black': 1000,
	'purple': 5000,
	'yellow': 10000,
};

const isBlackText = ['white', 'yellow'];

let previousBets = null;

let bets = {
	isSkipped: false,
	straightUp: {
		"0": 0,
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
		"5": 0,
		"6": 0,
		"7": 0,
		"8": 0,
		"9": 0,
		"10": 0,
		"11": 0,
		"12": 0,
		"13": 0,
		"14": 0,
		"15": 0,
		"16": 0,
		"17": 0,
		"18": 0,
		"19": 0,
		"20": 0,
		"21": 0,
		"22": 0,
		"23": 0,
		"24": 0,
		"25": 0,
		"26": 0,
		"27": 0,
		"28": 0,
		"29": 0,
		"30": 0,
		"31": 0,
		"32": 0,
		"33": 0,
		"34": 0,
		"35": 0,
		"36": 0
	},
	'color': {
		'black': 0,
		'red': 0
	},
	'evenOdd': {
		'even': 0,
		'odd': 0
	},
	'lowHigh': {
		'low': 0,
		'high': 0
	},
	'dozens': {
		'1st': 0,
		'2nd': 0,
		'3rd': 0
	},
	'columns': {
		'1st': 0,
		'2nd': 0,
		'3rd': 0
	}
};

let chip_color_selected = null;
let cursor_icon = 'not-allowed';

// get all the betting board buttons
const bettingBoardBtns = document.getElementById('betting_board').querySelectorAll('button');

function buildWheel(){
	const container = document.getElementById('wheel_container');
	const wheel = document.createElement('div');
    wheel.setAttribute('id', 'wheel');
	wheel.setAttribute('class', 'wheel');

	let outerRim = document.createElement('div');
	outerRim.setAttribute('class', 'outerRim');
	wheel.append(outerRim);

	let numbers = [0, 32, 15, 19, 4, 21, 2, 25, 17, 34, 6, 27, 13, 36, 11, 30, 8, 23, 10, 5, 24, 16, 33, 1, 20, 14, 31, 9, 22, 18, 29, 7, 28, 12, 35, 3, 26];
	for(i = 0; i < numbers.length; i++){
		let a = i + 1;
		let spanClass = (numbers[i] < 10)? 'single' : 'double';
		let sect = document.createElement('div');
		sect.setAttribute('id', 'sect'+a);
		sect.setAttribute('class', 'sect');
		let span = document.createElement('span');
		span.setAttribute('class', spanClass);
		span.innerText = numbers[i];
		sect.append(span);
		let block = document.createElement('div');
		block.setAttribute('class', 'block');
		sect.append(block);
		wheel.append(sect);
	}

	let pocketsRim = document.createElement('div');
	pocketsRim.setAttribute('class', 'pocketsRim');
	wheel.append(pocketsRim);

	let ballTrack = document.createElement('div');
    ballTrack.setAttribute('id', 'ballTrack');
	ballTrack.setAttribute('class', 'ballTrack');
	let ball = document.createElement('div');
	ball.setAttribute('class', 'ball');
	ballTrack.append(ball);
	wheel.append(ballTrack);

	let pockets = document.createElement('div');
	pockets.setAttribute('class', 'pockets');
	wheel.append(pockets);

	let cone = document.createElement('div');
	cone.setAttribute('class', 'cone');
	wheel.append(cone);

	let turret = document.createElement('div');
	turret.setAttribute('class', 'turret');
	turret.setAttribute('id', 'turretContainer');
	wheel.append(turret);

	let turretHandle = document.createElement('div');
	turretHandle.setAttribute('class', 'turretHandle');
	turretHandle.setAttribute('id', 'turretHandleContainer');
	let thendOne = document.createElement('div');
	thendOne.setAttribute('class', 'thendOne');
	turretHandle.append(thendOne);
	let thendTwo = document.createElement('div');
	thendTwo.setAttribute('class', 'thendTwo');
	turretHandle.append(thendTwo);
	wheel.append(turretHandle);

	let displayContainer = document.createElement('div');
	displayContainer.setAttribute('class', 'wheelDisplay');
	let display = document.createElement('p');
	display.setAttribute('id', 'wheelDisplay');
	display.hidden = true;
	displayContainer.append(display);
	wheel.append(displayContainer);

	container.append(wheel);
}

function spinWheel(winningSpin){
    let wheel = document.getElementById('wheel');
    let ballTrack = document.getElementById('ballTrack');

	let degree = 0;
	for(i = 0; i < wheelnumbersAC.length; i++){
		if(wheelnumbersAC[i] == winningSpin){
			degree = (i * 9.73) + 362;
		}
	}

	let numOfDots = 0;
	const spinningMessageTimeout = setInterval(function(){
		numOfDots = (numOfDots + 1) % 4;
		let spinningMessage = 'Spinning' + '.'.repeat(numOfDots);
		setMessageDisplay(spinningMessage, 'white');
	}, 500);

	wheel.style.cssText = 'animation: wheelRotate 5s linear infinite;';
	ballTrack.style.cssText = 'animation: ballRotate 1s linear infinite;';

	setTimeout(function(){
		ballTrack.style.cssText = 'animation: ballRotate 2s linear infinite;';
		style = document.createElement('style');
		style.type = 'text/css';
		style.innerText = '@keyframes ballStop {from {transform: rotate(0deg);}to{transform: rotate(-'+degree+'deg);}}';
		document.head.appendChild(style);
	}, 2000);
	setTimeout(function(){
		ballTrack.style.cssText = 'animation: ballStop 3s linear;';
	}, 6000);
	setTimeout(function(){
		ballTrack.style.cssText = 'transform: rotate(-'+degree+'deg);';
		clearInterval(spinningMessageTimeout);
	}, 9000);
	setTimeout(function(){
		wheel.style.cssText = '';
		style.remove();

		hideTurret(true);
		insertIntoResultsHistory(winningSpin);
		result = `${winningSpin}${getColorSymbol(getColor(winningSpin))}`
		showOnWheelDisplay(result);
	}, 10000);
}

bettingBoardBtns.forEach(button => {
	button.addEventListener('mouseenter', () => {
	  // Change the cursor icon to a custom one
	  button.style.cursor = cursor_icon;
	});
  
	button.addEventListener('mouseleave', () => {
	  // Reset the cursor icon to the default
	  button.style.cursor = 'auto';
	});

	button.addEventListener('mouseup', changeBets);
});

function shortenNumber(num) {
	let divsor = 1;
	let symbol = '';
    if (num >= 1_000_000_000) {
		divsor = 1_000_000_000;
		symbol = 'B';
    } else if (num >= 1_000_000) {
		divsor = 1_000_000;
		symbol = 'M';
    } else if (num >= 1_000) {
		divsor = 1_000;
		symbol = 'k';
    } else {
        return num.toString();
    }

	if (num % divsor === 0) {
		return (num / divsor) + symbol;
	}

	let shortenedNum = (num / divsor).toFixed(2);
	if (shortenedNum[shortenedNum.length - 1] === '0') {
		shortenedNum = shortenedNum.slice(0, -1);
	}
	return shortenedNum + symbol;
}

function unshortenNumber(num) {
	if (num.includes('B')) {
		return parseFloat(num.split('B')[0]) * 1_000_000_000;
	} else if (num.includes('M')) {
		return parseFloat(num.split('M')[0]) * 1_000_000;
	} else if (num.includes('k')) {
		return parseFloat(num.split('k')[0]) * 1_000;
	} else {
		return parseFloat(num);
	}
}

function insertIntoResultsHistory(result) {
	for (let i = 8; i >= 2; i--) {
		const resultElement = document.getElementById(`result_${i-1}`);
		if (resultElement.innerText === '') {
			continue;
		}

		const newResult = document.getElementById(`result_${i}`);
		newResult.innerText = resultElement.innerText;
	}

	const latestResult = document.getElementById('result_1');
	latestResult.innerText = result + getColorSymbol(getColor(result));

	latestResult.style.cssText = 'animation: flash 1s linear infinite;';

	setTimeout(function(){
		latestResult.style.cssText = '';
	}, 2000);
}

function setMessageDisplay(message, color) {
	const display = document.getElementById('message_display');
	display.innerText = message;
	display.style.cssText = 'color: '+color+';';

	if (color === 'green') {
		display.style.cssText += 'background-color: #fff;';
		display.style.cssText += 'animation: flash 1s linear 2;';
	}
}

function setTotalBetDisplay(totalValue) {
	const display = document.getElementById('totalBetDisplay');
	display.innerText = 'Amount Bet: '+shortenNumber(totalValue)+'c';
}

function showOnWheelDisplay(result) {
	const display = document.getElementById('wheelDisplay');
	display.innerText = result;
	display.hidden = false;
	display.style.cssText = 'opacity: 0; animation: fadeIn 1s linear;';

	setTimeout(function(){
		display.style.cssText = 'opacity: 1;';
	}, 1000);

	setTimeout(function(){
		display.style.cssText = 'opacity: 1; animation: fadeIn 1s linear reverse;';

		setTimeout(function(){
			display.innerText = '';
			hideTurret(false);
		}, 1000);
	}, 3000);
}

function displayValuesOfChips() {
	for (const color in chipValues) {
		const chipValue = chipValues[color];
		const chipText = document.getElementById(color+'_chip_label');
		chipText.innerText = chipValue;
	}
}

function selectBet(color) {
	bet_selected = chipValues[color];
	if (bet_selected === undefined) {
		console.error('Invalid chip value');
		return;
	}

	// get all the bet images under parent with id 'chips'
	const betImages = document.getElementById('chips').querySelectorAll('img');
	betImages.forEach(betImg => {
		betImg.style.cssText = '';
	});

	const betImg = document.getElementById(color+'_chip');
	betImg.style.cssText = 'border: 2px solid #ffF;';

	cursor_icon = `url(/static//images/roulette/cursor_chip_${color}.png), auto`;
	chip_color_selected = color;
}

function changeBets(event) {
	// get the bet amount from the selected chip
	let betAmount = chipValues[chip_color_selected];
	if (betAmount === undefined) {
		chipsLabel = document.getElementById('label_chips');
		chipsLabel.style.cssText = 'color: red; font-size: 32px;';

		setTimeout(function(){
			chipsLabel.style.cssText = 'color: white;';
		}, 2000);
		return;
	}

	// if right click, then subtract the bet amount
	if (event.button === 2) {
		betAmount = -betAmount;
	}

	const buttonPressed = event.target;
	const parentDiv = buttonPressed.parentElement;
	const betTextElement = parentDiv.querySelector('p');
	const imgElement = parentDiv.querySelector('img');
	const buttonId = buttonPressed.id;


	const betTypes = buttonId.split('_');
	if (betTypes[0] !== 'bet') {
		return;
	}

	const betAlreadyPlaced = bets[betTypes[1]][betTypes[2]];

	// if the type is not found, then the bet is not valid
	if (betAlreadyPlaced === undefined || betAlreadyPlaced === null) {
		return;
	}

	imgElement.src = `/static/images/roulette/${chip_color_selected}_chip_blank.png`;
	if (isBlackText.includes(chip_color_selected)) {
		betTextElement.className = 'text-black';
	} else {
		betTextElement.className = 'text-white';
	}

	if (betAlreadyPlaced + betAmount <= 0) {
		betAmount = -betAlreadyPlaced;
		betTextElement.hidden = true;
		imgElement.hidden = true;
	} else {
		betTextElement.hidden = false;
		imgElement.hidden = false;
	}

	if (!event.isPrevious) {
		bets[betTypes[1]][betTypes[2]] += betAmount;
	}

	betTextElement.innerText = shortenNumber(bets[betTypes[1]][betTypes[2]])+'c';

	const totalValue = calculateTotalValueOfBets();
	setTotalBetDisplay(totalValue);
}

function repeatBets() {
	if (previousBets === null) {
		return;
	}

	bets = previousBets; // JSON.parse(JSON.stringify(previousBets));
	for (const betType in bets) {
		for (const bet in bets[betType]) {
			if (bets[betType][bet] > 0) {
				const buttonId = `bet_${betType}_${bet}`;
				const button = document.getElementById(buttonId);
				changeBets({target: button, isPrevious: true});
			}
		}
	}
}

function clearBets() {
	// get all elements with class 'button-container' or 'button-container-no-rotate'
	const buttonContainers = document.querySelectorAll('.button-container, .button-container-no-rotate', 'w');
	buttonContainers.forEach(buttonContainer => {
		const imgElement = buttonContainer.querySelector('img');
		const betTextElement = buttonContainer.querySelector('p');

		betTextElement.hidden = true;
		imgElement.hidden = true;
	});

	for (const betType in bets) {
		for (const bet in bets[betType]) {
			bets[betType][bet] = 0;
		}
	}
}

function calculateTotalValueOfBets() {
	let totalValue = 0;
	for (const betType in bets) {
		for (const bet in bets[betType]) {
			totalValue += bets[betType][bet];
		}
	}

	return totalValue;
}

function hideTurret(hide) {
	const turret = document.getElementById('turretContainer');
	const turretHandle = document.getElementById('turretHandleContainer');

	turret.hidden = hide;
	turretHandle.hidden = hide;
}

function getColor(number) {
	if (number === 0) {
		return 'green';
	} 

	if (redNumbers[number]) {
		return 'red';
	} else {
		return 'black';
	}
}

function getColorSymbol(color) {
	if (color === 'red') {
		return '🔴';
	} else if (color === 'black') {
		return '⚫';
	} else {
		return '🟢';
	}
}

buildWheel();

// wait 1 second before spinning
// setTimeout(function(){
//     // random number between 0 and 36
//     randNum = Math.floor(Math.random() * 36);
//     console.log(randNum);
//     spinWheel(randNum);
// }, 100);

spinButton = document.getElementById('spin');
spinButton.addEventListener('click', async function() {
	allButtons = document.querySelectorAll('button');
	allButtons.forEach(button => {
		button.disabled = true;
	});
	const oldCursor = cursor_icon;
	cursor_icon = 'not-allowed';
	fetchResults();
	
	setTimeout(function(){
		allButtons.forEach(button => {
			button.disabled = false;
		});
		cursor_icon = oldCursor;
	}, 10000);
});

async function fetchResults() {
	console.log(bets);
	if (calculateTotalValueOfBets !== 0) {
		previousBets = await JSON.parse(JSON.stringify(bets));
	}

	const url = 'http://' + DOMAIN + '/api/v1/games/rl';
	const req = await fetch(url, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
		},
		body: JSON.stringify(bets),
	});

	const data = await req.json();
	console.log(data);

	if (data['error']) {
		setMessageDisplay(data['error']);
		return;
	}

	const winningNumber = Number(data['winningNumber']);
	spinWheel(winningNumber);
	const msg = `Bet: ${data['totalBets']} - Returned: ${data['amountReturned']}`;
	// if amount returned is 0 and amount bet is not 0, then display red message, if amount returned is 0, white, else green
	const msgColor = (data['amountReturned'] === 0 && data['totalBets'] !== 0)? 'red' : (data['amountReturned'] === 0)? 'white' : 'green';

	setTimeout(function(){
		setMessageDisplay(msg, msgColor);
		updateBalance(data['newBalance'], data['amountReturned']);

		if (data['amountReturned'] === 0) {
			clearBets();
		}

		setTimeout(function(){
			setMessageDisplay('Place your bets!', 'white');
			const totalValue = calculateTotalValueOfBets();
			setTotalBetDisplay(totalValue);
		}, 3000);
	}, 10000);
}

function updateBalance(newBalance, returnAmount) {
	const balanceElement = document.getElementById('balance_value');
	if (returnAmount === 0) {
		balanceElement.style.cssText = 'color: red;';
	}
	if (returnAmount > 0) {
		balanceElement.style.cssText = 'color: green;';
	}
	balanceElement.innerText = '+'+shortenNumber(returnAmount)+'c';

	setTimeout(function(){
		balanceElement.style.cssText = 'color: black;';
		balanceElement.innerText = shortenNumber(newBalance)+'c';
	}, 2000);
}

function skipSpin() {
	clearBets();
	bets.isSkipped = true;
	fetchResults();
}