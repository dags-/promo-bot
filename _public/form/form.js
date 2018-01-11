window.addEventListener("load", function () {
    updateType("server");
    loadExisting("server");
});

function loadExisting(type) {
    var slash = window.location.href.lastIndexOf("/") + 1;
    var id = window.location.href.substr(slash);
    var url = "/api/" + type + "/" + id;
    var req = new XMLHttpRequest();

    req.onload = function() {
        if (this.readyState === 4 && this.status === 200) {
            var promo = JSON.parse(this.responseText);
            var pid = promo["id"];
            if (pid !== undefined && pid === id) {
                loadPromotion(promo);
            }
        }
    };

    req.open("GET", url);
    req.send();
}

function loadPromotion(promo) {
    for (var key in promo) {
        if (!promo.hasOwnProperty(key)) {
            continue;
        }
        if (key === "type" || key === "id") {
            continue;
        }

        var el = document.getElementById(key);
        if (el === undefined) {
            return;
        }

        if (key === "whitelist") {
            el.checked = promo[key];
            el.dispatchEvent(new Event("change"));
        } else {
            el.value = promo[key];
            el.dispatchEvent(new Event("input"));

        }
    }
}

function updateType(type) {
    var content = document.getElementById("content");
    var current = document.getElementById("form");
    var form = buildForm(type);

    if (current != null) {
        content.removeChild(current);
    }

    content.appendChild(form);
}

function buildForm(type) {
    var preview = document.getElementById("pr-post");
    var header = document.getElementById("pr-header");
    var form = makeForm("form", {id: "form", method: "post"});

    form.appendChild(selector(type));
    form.appendChild(name());
    form.appendChild(description());
    form.appendChild(icon());
    form.appendChild(image());
    form.appendChild(website());
    form.appendChild(discord());
    form.appendChild(tags());

    if (type === "server") {
        form.appendChild(ip());
        form.appendChild(whitelist());
        preview.style.borderColor = "#00d56a";
        header.innerText = "#Server";
        document.getElementById("pr-ip-field").style.display = "inline-block";
        document.getElementById("pr-whitelist-field").style.display = "inline-block";
    } else {
        document.getElementById("pr-ip-field").style.display = "none";
        document.getElementById("pr-whitelist-field").style.display = "none";
    }

    if (type === "twitch") {
        header.innerText = "#Twitch";
        preview.style.borderColor = "#0080ff";
    }

    if (type === "youtube") {
        header.innerText = "#Youtube";
        preview.style.borderColor = "#ff8080";
    }

    return form;
}

function selector(type) {
    var label = makeElement("label");
    label.innerText = "What are you promoting?";

    var input = makeElement("select", {id: "type", name: "type"});
    var options = ["server", "twitch", "youtube"];
    for (var i = 0; i < options.length; i++) {
        var name = options[i];
        var option = makeElement("option");
        option.innerText = name;
        input.appendChild(option);
        if (name === type) {
            input.selectedIndex = i;
        }
    }

    input.addEventListener("change", function() {
        var selected = this.selectedIndex;
        var type = this.options[selected].value;
        updateType(type);
    });

    return wrap("section", label, input);
}

function name() {
    var label = makeElement("label");
    label.innerText = "What is your server/twitch/youtube name?";

    var input = makeElement("input", {id: "name", name: "name", type: "text", maxlength: 120});
    input.addEventListener("input", function() {
        var target = document.getElementById("pr-name");
        target.innerText = this.value;
    });

    return wrap("section", label, input);
}

function description() {
    var label = makeElement("label");
    label.innerText = "Describe what you're all about";

    var input = makeElement("textarea", {id: "description", name: "description", maxlength: 480});
    input.addEventListener("input", function() {
        var target = document.getElementById("pr-description");
        target.innerText = this.value;
    });

    return wrap("section", label, input);
}

function icon() {
    var label = makeElement("label");
    label.innerText = "Link an icon image to accompany your promotion (optional)";

    var input = makeElement("input", {id: "icon", name: "icon", type: "text", maxlength: 240});
    input.addEventListener("input", function() {
        var large = document.getElementById("pr-icon");
        var small = document.getElementById("pr-icon-small");
        if (this.value === "") {
            large.style.display = "none";
            small.style.display = "none";
        } else {
            large.style.display = "block";
            small.style.display = "inline-block";
            large.src = this.value;
            small.src = this.value;
        }
    });

    return wrap("section", label, input);
}

function image() {
    var label = makeElement("label");
    label.innerText = "Link an image to accompany your promotion (optional)";

    var input = makeElement("input", {id: "image", name: "image", type: "text", maxlength: 240});
    input.addEventListener("input", function() {
        var target = document.getElementById("pr-image");
        if (this.value === "") {
            target.style.display = "none";
        } else {
            target.style.display = "block";
            target.src = this.value;
        }
    });

    return wrap("section", label, input);
}

function website() {
    var label = makeElement("label");
    label.innerText = "Got a website you'd like people to visit? (optional)";

    var input = makeElement("input", {id: "website", name: "website", type: "text", maxlength: 120});
    input.addEventListener("input", function () {
        var field = document.getElementById("pr-website-field");
        var target = document.getElementById("pr-website");
        if (this.value === "") {
            field.style.display = "none";
        } else {
            field.style.display = "inline-block";
            target.innerText = this.value;
            target.href = this.value;
        }
    });

    return wrap("section", label, input);
}

function discord() {
    var label = makeElement("label");
    label.innerText = "Got a discord you'd like people to join? (optional)";

    var input = makeElement("input", {id: "discord", name: "discord", type: "text", maxlength: 120});
    input.addEventListener("input", function () {
        var field = document.getElementById("pr-discord-field");
        var target = document.getElementById("pr-discord");
        if (this.value === "") {
            target.innerText = "#Join";
            field.style.display = "none";
        } else {
            field.style.display = "inline-block";
            target.innerText = "#Join";
            target.href = this.value;
        }
    });

    return wrap("section", label, input);
}

function tags() {
    var label = makeElement("label");
    label.innerText = "Provide some comma separated keywords (optional)";

    var input = makeElement("input", {id: "tags", name: "tags", type: "text", maxlength: 120});
    input.addEventListener("input", function () {
        var target = document.getElementById("pr-footer");
        var split = this.value.split(",");
        var text = "";
        for (var i = 0; i < split.length; i++) {
            var tag = split[i].trim();
            if (tag === "") {
                continue;
            }
            text += i > 0 ? " #" + tag.trim() : "#" + tag.trim();
        }
        if (text === "") {
            target.innerText = "#promo";
        } else {
            target.innerText = text;
        }
    });

    return wrap("section", label, input);
}

function ip() {
    var label = makeElement("label");
    label.innerText = "What's the server's ip address?";

    var input = makeElement("input", {id: "ip", name: "ip", type: "text", maxlength: 120});
    input.addEventListener("input", function () {
        var target = document.getElementById("pr-ip");
        if (this.value === "") {
            target.innerText = "required!";
        } else {
            target.innerText = this.value;
        }
    });

    return wrap("section", label, input);
}

function whitelist() {
    var label = makeElement("label");
    label.innerText = "Does your server use a whitelist?";

    var input = makeElement("input", {id: "whitelist", name: "whitelist", type: "checkbox"});
    input.addEventListener("change", function () {
        var target = document.getElementById("pr-whitelist");
        if (this.checked) {
            target.innerText = "Yes";
        } else {
            target.innerText = "No";
        }
    });

    return wrap("section", label, input);
}

function wrap(type) {
    var section = document.createElement(type);
    for (var i = 1; i < arguments.length; i++) {
        section.appendChild(arguments[i]);
    }
    return section;
}

function makeForm(type, attributes) {
    var el = document.createElement(type);

    for (var k in attributes) {
        if (attributes.hasOwnProperty(k)) {
            el.setAttribute(k, attributes[k]);
        }
    }

    for (var i = 2; i < arguments.length; i++) {
        el.classList.add(arguments[i]);
    }

    return el;
}

function makeElement(type, attributes) {
    var el = document.createElement(type);

    for (var k in attributes) {
        if (attributes.hasOwnProperty(k)) {
            el.setAttribute(k, attributes[k]);
        }
    }

    for (var i = 2; i < arguments.length; i++) {
        el.classList.add(arguments[i]);
    }

    if (attributes !== undefined && attributes["id"] !== undefined) {
        var current = document.getElementById(attributes["id"]);
        if (current !== null) {
            el.innerText = current.innerText;
            el.value = current.value;
        }
    }

    return el;
}