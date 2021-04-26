// xhr module

function opts(method, contentType, body) {
	let type = contentType == 'json' ? 'application/json' : 'text/html'
	return {
		method: method,
		body: body,
		headers: { "Content-Type": type },
		mode: 'cors'
	}
}

function load(method, contentType, body, url, callback) {
	let options = opts(method, contentType, body)
	fetch(url, options).then(function (response) {
		if (response.ok) {
			if (contentType == 'json') {
				return response.json();
			}
			return response.text();
		} else {
			return Promise.reject(response);
		}
	}).then(function (data) {
		callback(data);
	}).catch(function (err) {
		console.error(err);
	});
}

// elementOrId is a DOM element or the id of a DOM element
function setContent(elementOrId) {
    let element = elementOrId
    if (!isElement(elementOrId)) {
        element = document.getElementById(elementOrId)
    }
    if (!isElement(element)) {
        console.error("XHR target is not an element.", element)
    }
	return function(htmlData) { element.innerHTML = htmlData }
}

function deleteContent(elementId) {
	return function(textResponse) {
		document.getElementById(elementId).innerHTML = ""
		alert(textResponse)
	}
}

function loadIntoElement(method, url, elementId) {
	load(method, "html", null, url, setContent(elementId))
}

function appendToElement(method, url, elementId) {
	load(method, "html", null, url, setContent(elementId))
}

function getJSON(url, callback) {
	load("GET", "json", null, url, callback)
}

function isElement(element) {
    return element instanceof Element || element instanceof HTMLDocument;
}

export function initXHR() {
	let xhrItems = document.querySelectorAll("[data-xhr-url][data-xhr-target]");
	xhrItems.forEach(function(item){
        // method can be get, put, post, delete
        // action can be replace or append
        let method = item.dataset.xhrMethod || "get"
        let fn = item.dataset.xhrAction || "replace"
        let url = item.dataset.xhrUrl
        let target = item.dataset.xhrTarget
		item.addEventListener("click", function (event) {
            if (fn == "append") {
			    appendToElement(method, url, target)
            } else {
                loadIntoElement(method, url, target)
            }
		});
	});

    let modalItems = document.querySelectorAll("[data-modal][data-xhr-url]");
	modalItems.forEach(function(item){
		//console.log(item.dataset.xhrUrl, item.dataset.xhrTarget)
        let modal = document.getElementById(item.dataset.modal)
        let modalContentDiv = modal.querySelector('.modal-content');
		item.addEventListener("click", function (event) {
			getHTML(item.dataset.xhrUrl, modalContentDiv)
		});
	});
}
