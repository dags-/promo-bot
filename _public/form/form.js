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

        var colors = ["#00d56a", "#0080ff", "#ff8080"];
        var preview = document.getElementById("pr-post");
        preview.style.borderColor = colors[this.selectedIndex];

        var header = document.getElementById("pr-header");
        header.innerText = "#" + type.charAt(0).toUpperCase() + type.slice(1);
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
    }, function () {
        var target = document.getElementById("pr-" + this.id);
        target.innerText = this.value;
    }));

    form.appendChild(makeTextArea({
        name: "description",
        id: "description",
        type: "textarea",
        maxlength: 480,
        desc: "Describe what you're all about"
    }, function () {
        var target = document.getElementById("pr-" + this.id);
        target.innerText = this.value;
    }));

    form.appendChild(makeInput({
        name: "icon",
        id: "icon",
        type: "text",
        maxlength: 240,
        desc: "Link an icon image to accompany your promotion (optional)"
    }, function () {
        var icon = document.getElementById("pr-icon");
        var smallIcon = document.getElementById("pr-icon-small");
        if (this.value === "") {
            icon.style.display = "none";
            smallIcon.style.display = "none";
        } else {
            icon.style.display = "block";
            smallIcon.style.display = "inline-block";
            icon.src = this.value;
            smallIcon.src = this.value;
        }
    }));

    form.appendChild(makeInput({
        name: "image",
        id: "image",
        type: "text",
        maxlength: 240,
        desc: "Link an image to accompany your promotion (optional)"
    }, function () {
        var image = document.getElementById("pr-image");
        if (this.value === "") {
            image.style.display = "none";
        } else {
            image.style.display = "block";
            image.src = this.value;
        }
    }));

    form.appendChild(makeInput({
        name: "website",
        id: "website",
        type: "text",
        maxlength: 240,
        desc: "Got a website you'd like people to visit? (optional)"
    }, function() {
        var field = document.getElementById("pr-website-field");
        var site = document.getElementById("pr-website");
        if (this.value === "") {
            field.style.display = "none";
        } else {
            field.style.display = "inline-block";
            site.innerText = this.value;
            site.href = this.value;
        }
    }));

    form.appendChild(makeInput({
        name: "discord",
        id: "discord",
        type: "text",
        maxlength: 120,
        desc: "Got a discord you'd like people to join? (optional)"
    }, function() {
        var field = document.getElementById("pr-discord-field");
        var site = document.getElementById("pr-discord");
        if (this.value === "") {
            field.style.display = "none";
        } else {
            field.style.display = "inline-block";
            site.href = this.value;
        }
    }));

    switch (type) {
        case "server":
            document.getElementById("pr-ip-field").style.display = "inline-block";
            document.getElementById("pr-whitelist-field").style.display = "inline-block";
            addServer(form);
            break;
        case "twitch":
        case "youtube":
            document.getElementById("pr-ip-field").style.display = "none";
            document.getElementById("pr-whitelist-field").style.display = "none";
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
        desc: "What's the server's ip address?"
    }, function() {
        var ip = document.getElementById("pr-ip");
        if (this.value === "") {
            ip.innerText = "required";
        } else {
            ip.innerText = this.value;
        }
    }));

    form.appendChild(makeInput({
        name: "whitelist",
        id: "whitelist",
        type: "checkbox",
        desc: "Does your server use a whitelist?"
    }, function () {
        var whitelist = document.getElementById("pr-whitelist");
        console.log(this.checked);
        if (this.checked) {
            whitelist.innerText = "Yes";
        } else {
            whitelist.innerText = "No";
        }
    }));
}

function makeInput(props, listener) {
    var section = document.createElement("section");

    if (props["desc"] !== undefined) {
        var label = document.createElement("label");
        label.innerText = props["desc"];
        section.appendChild(label);
    }

    var input = makeElement("input", props);
    var event = props["type"] === "checkbox" ? "change" : "input";
    input.addEventListener(event, listener);

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

function makeTextArea(props, listener) {
    var label = document.createElement("label");
    label.innerText = props["desc"];

    var text = makeElement("textarea", props);
    text.addEventListener("input", listener);

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