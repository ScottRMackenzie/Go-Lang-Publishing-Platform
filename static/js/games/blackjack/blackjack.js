ws_url = `ws://${DOMAIN}/ws/games/bj`;
let socket = new WebSocket(ws_url);

socket.onopen = function(event) {
    disableInput(false);
};

socket.onmessage = function(event) {
    const message = JSON.parse(event.data);

    switch (message.type) {
        case 'balance':
            updateBalance(message.score);
            break;
        case 'player_hand':
            document.getElementById('gameStateText').innerHTML = 'Your turn!';
            displayCards('player_hand', message.data);
            displayScore('player_score', message.score);
            break;
        case 'dealer_hand':
            displayCards('dealer_hand', message.data);
            displayScore('dealer_score', message.score);
            break;
        case 'dealers_turn':
            disableInput(true);
            hideInput(true);
            document.getElementById('gameStateText').innerHTML = 'Dealer\'s turn!';
            break;
        case 'result':
            resultText = message.data + " - You received back " + message.score + "!";
            document.getElementById('gameStateText').innerHTML = resultText;
            document.getElementById('newGame').hidden = false;
            break;
        case 'error':
            disableInput(false);
            document.getElementById('newGame').hidden = false;
            document.getElementById('gameState').hidden = true;
            alert(message.data);
            break;
        default:
            console.error('Unknown message type:', message.type);
    };
};

socket.onclose = function(event) {
    console.log("Disconnected from WebSocket server.");
};

document.getElementById('newGame').addEventListener('click', function() {
    const betAmount = document.getElementById('betAmount').value;
    if (socket.readyState != WebSocket.OPEN) {
        console.error('WebSocket is not open:', socket.readyState);
        return;
    }
    socket.send(JSON.stringify({ type: 'new_game', value: betAmount }));

    document.getElementById('gameState').hidden = false;
    document.getElementById('newGame').hidden = true;

    hideInput(false);
    disableInput(false);
});

document.getElementById('hitButton').addEventListener('click', function() {
    socket.send('hit');
});

document.getElementById('standButton').addEventListener('click', function() {
    socket.send('stand');
    disableInput(true);
});

function displayCards(elementID, cards) {
    cardHandDiv = document.getElementById(elementID);
    cardHandDiv.innerHTML = '';

    for (let card of cards) {
        var newCard = document.createElement('img');
        newCard.src = `/static/images/cards/${card.value}${card.suit}.png`;
        newCard.alt = `${card.value} of ${card.suit}`;
        newCard.classList.add('w-32');
        newCard.classList.add('h-auto');
        newCard.classList.add('inline');
        cardHandDiv.appendChild(newCard);
    }

    if (elementID == 'dealer_hand' && cards.length == 1) {
        var backCard = document.createElement('img');
        backCard.src = `/static/images/cards/back.png`;
        backCard.alt = `back of card`;
        backCard.classList.add('w-32');
        backCard.classList.add('h-auto');
        backCard.classList.add('inline');
        cardHandDiv.appendChild(backCard);
    }
}

function displayScore(elementID, score) {
    scoreHandDiv = document.getElementById(elementID);
    if (score == 'Bust') {
        scoreHandDiv.style.color = 'red';
    } else if (score == 'Blackjack') {
        scoreHandDiv.style.color = 'green';
    } else {
        scoreHandDiv.style.color = 'black';
    }
    scoreHandDiv.innerText = score;
}

function disableInput(isDisabled) {
    document.getElementById('hitButton').disabled = isDisabled;
    document.getElementById('standButton').disabled = isDisabled;
}

function hideInput(isHidden) {
    document.getElementById('hitButton').hidden = isHidden;
    document.getElementById('standButton').hidden = isHidden;
}

function updateBalance(balance) {
    balance = Number(balance);
    balanceDisplay = document.getElementById('balance_value');
    ogBalance = Number(balanceDisplay.innerText);
    balanceDiff = balance - ogBalance;

    if (balanceDiff < 0) {
        balanceDisplay.style.color = 'red';
    } else if (balanceDiff > 0) {
        balanceDisplay.style.color = 'green';
    } else {
        balanceDisplay.textContent = balance;
        return;
    }

    const duration = 3000; // 3 seconds
    const incrementTime = 100; // Increment every 100 milliseconds
    const totalIncrements = duration / incrementTime;
    const incrementStep = (balance - ogBalance) / totalIncrements;

    let currentIncrement = ogBalance;

    const interval = setInterval(() => {
        console.log("Current increment:", currentIncrement);
        currentIncrement += incrementStep;
        balanceDisplay.textContent = Math.round(currentIncrement);

        // Stop the interval when the new amount is reached
        if (Math.round(currentIncrement) >= balance) {
            clearInterval(interval);
            balanceDisplay.style.color = 'black';
            balanceDisplay.textContent = balance; // Ensure it ends exactly at newAmount
        }
    }, incrementTime);
}