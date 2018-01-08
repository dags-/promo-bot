window.addEventListener("load", function () {
    newForm("server");
});

function newForm(type) {
    var content = document.getElementById("content");
    var oldForm = document.getElementById("form");
    var form = makeElement("form", {id: "form", method: "post"});

    addSelector(form, type);
    addBody(form, type);
    addSubmit(form);

    if (oldForm !== null) {
        content.removeChild(oldForm);
    }

    content.appendChild(form);
}

function addSelector(form, type) {
    var label = document.createElement("label");
    label.innerText = "What are you promoting?";

    var select = makeElement("select", {id: "type", name: "type"});

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
    select.addEventListener("change", function () {
        var type = this.options[this.selectedIndex].value;
        newForm(type);
    });

    var section = document.createElement("section");
    section.appendChild(label);
    section.appendChild(select);

    form.appendChild(section);
}

function addBody(form, type) {
    form.appendChild(makeInput({
        name: "name",
        id: "name",
        type: "text",
        maxlength: 120,
        desc: "What is your server/twitch/youtube name?"
    }));
    form.appendChild(makeTextArea({
        name: "description",
        id: "description",
        type: "textarea",
        maxlength: 480,
        desc: "Describe what you're all about"
    }));
    form.appendChild(makeInput({
        name: "icon",
        id: "icon",
        type: "text",
        maxlength: 240,
        desc: "Link an icon image to accompany your promotion (optional)"
    }));
    form.appendChild(makeInput({
        name: "media",
        id: "media",
        type: "text",
        maxlength: 240,
        desc: "Link an image to accompany your promotion (optional)"
    }));
    form.appendChild(makeInput({
        name: "link",
        id: "link",
        type: "text",
        maxlength: 240,
        desc: "Got a website or discord you'd like people to visit? (optional)"
    }));

    switch (type) {
        case "server":
            addServer(form);
            break;
        case "twitch":
            break;
        case "youtube":
            break;
    }
}

function addSubmit(form) {
    var submit = makeInput({type: "submit"});
    form.appendChild(submit);
}

function addServer(form) {
    form.appendChild(makeInput({
        name: "ip",
        id: "ip",
        type: "text",
        maxlength: "120",
        desc: "What\"s the server\"s ip address?"
    }));
    form.appendChild(makeInput({
        name: "whitelist",
        id: "whitelist",
        type: "checkbox",
        desc: "Does your server use a whitelist?"
    }));
}

function makeInput(props) {
    var section = document.createElement("section");

    if (props["desc"] !== undefined) {
        var label = document.createElement("label");
        label.innerText = props["desc"];
        section.appendChild(label);
    }

    var input = makeElement("input", props);
    var current = document.getElementById(props['id']);
    if (current !== null) {
        input.innerText = current.innerText;
        input.value = current.value;
    }

    section.appendChild(input);

    return section;
}

function makeElement(type, props) {
    var el = document.createElement(type);
    setProperties(el, props);
    return el;
}

function makeTextArea(props) {
    var label = document.createElement("label");
    label.innerText = props["desc"];

    var text = makeElement("textarea", props);

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