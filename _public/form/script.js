window.addEventListener('load', function() {
    document.getElementById('type').addEventListener('change', onTypeChange);
    selectType('server');
});

function onTypeChange() {
    var index = this.selectedIndex;
    var type = this.options[index].value;
    selectType(type);
}

function selectType(type) {
    var details = document.getElementById('details');
    clear(details);
    switch (type) {
        case "server":
            return showServerInput(details);
        case "twitch":
            return showTwitchInput(details);
        case "youtube":
            return showYoutubeInput(details);
    }
}

function showServerInput(div) {
    div.appendChild(makeInput("What's the server's ip address?", 'ip', 'text'));
    div.appendChild(makeInput("Does your server use a whitelist?", 'whitelist', 'checkbox'));
}

function showYoutubeInput(div) {
    // var title = makeInput("What's the name/title of your youtube channel?", 'title', 'text');
    // var address = makeInput("What's the url of your youtube channel?", 'url', 'text');
    // div.appendChild(title);
    // div.appendChild(address);
}

function showTwitchInput(div) {
    // var username = makeInput("What's your twitch username?", 'username', 'text');
    // var address = makeInput("What's the url for your twitch channel?", 'url', 'text');
    // div.appendChild(username);
    // div.appendChild(address);
}

function makeInput(title, name, type) {
    var label = document.createElement('label');
    label.innerText = title;

    var input = document.createElement('input');
    input.type = type;
    input.name = name;
    input.id = name;

    var section = document.createElement('section');
    section.appendChild(label);
    section.appendChild(input);

    return section;
}

function clear(div) {
    while (div.lastChild !== null) {
        div.removeChild(div.lastChild);
    }
}
