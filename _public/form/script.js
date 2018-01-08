window.addEventListener('load', function () {
    newForm('server');
});

function newForm(type) {
    var content = document.getElementById("content");
    var oldForm = document.getElementById("form");

    var form = document.createElement("form");
    form.setAttribute("id", "form");

    addSelector(form, type);
    addBody(form);

    switch (type) {
        case "server":
            addServer(form);
            break;
        case "twitch":
            break;
        case "youtube":
            break;
    }

    form.appendChild(input({type: "submit", method: "post"}));

    if (oldForm !== null) {
        content.removeChild(oldForm);
    }

    content.appendChild(form);
}

function addSelector(form, type) {
    var label = document.createElement("label");
    label.innerText = "What are you promoting?";

    var select = document.createElement("select");
    setProperties(select, {id: "type", name: "type"});

    var selected = 0;
    var options = ["server", "twitch", "youtube"];
    for (var i = 0; i < options.length; i++) {
        var name = options[i];
        var option = document.createElement("option");
        option.innerText = name;
        select.appendChild(option);
        if (name === type) {
            selected = i;
        }
    }

    select.selectedIndex = selected;
    select.addEventListener('change', function () {
        var type = this.options[this.selectedIndex].value;
        newForm(type);
    });

    var section = document.createElement("section");
    section.appendChild(label);
    section.appendChild(select);

    form.appendChild(section);
}

function addBody(form) {
    form.appendChild(input({
        name: "name",
        id: "name",
        type: "text",
        maxlength: 120,
        desc: "What is your server/twitch/youtube name?"
    }));
    form.appendChild(textArea({
        name: "description",
        id: "description",
        type: "textarea",
        maxlength: 480,
        desc: "Describe what you're all about"
    }));
    form.appendChild(input({
        name: "icon",
        id: "icon",
        type: "text",
        maxlength: 240,
        desc: "Link an icon image to accompany your promotion (optional)"
    }));
    form.appendChild(input({
        name: "media",
        id: "media",
        type: "text",
        maxlength: 240,
        desc: "Link an image to accompany your promotion (optional)"
    }));
    form.appendChild(input({
        name: "link",
        id: "link",
        type: "text",
        maxlength: 240,
        desc: "Got a website or discord you'd like people to visit? (optional)"
    }));
}

function addServer(form) {
    form.appendChild(input({
        name: "ip",
        id: "ip",
        type: "text",
        maxlength: "120",
        desc: "What\"s the server\"s ip address?"
    }));
    form.appendChild(input({
        name: "whitelist",
        id: "whitelist",
        type: "checkbox",
        desc: "Does your server use a whitelist?"
    }));
}

function input(props) {
    var section = document.createElement("section");

    if (props['desc'] !== undefined) {
        var label = document.createElement("label");
        label.innerText = props["desc"];
        section.appendChild(label);
    }

    var input = document.createElement("input");
    setProperties(input, props);

    var current = document.getElementById(props['id']);
    if (current !== null) {
        input.innerText = current.innerText;
        input.value = current.value;
    }

    section.appendChild(input);

    return section;
}

function textArea(props) {
    var label = document.createElement("label");
    label.innerText = props["desc"];

    var text = document.createElement("textarea");
    setProperties(text, props);

    var section = document.createElement("section");
    section.appendChild(label);
    section.appendChild(text);

    return section;
}

function setProperties(el, props) {
    for (var k in props) {
        if (k === "desc") {
            continue;
        }
        if (props.hasOwnProperty(k)) {
            var v = props[k];
            el.setAttribute(k, v);
        }
    }
}

function clear(div) {
    while (div.lastChild !== null) {
        div.removeChild(div.lastChild);
    }
}