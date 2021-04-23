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

function setContent(elementId) {
	return function(htmlData) { document.getElementById(elementId).innerHTML = htmlData }
}

function deleteContent(elementId) {
	return function(textResponse) {
		document.getElementById(elementId).innerHTML = ""
		alert(textResponse)
	}
}

function getHTML(url, elementId) {
	load("GET", "html", null, url, setContent(elementId))
}

function getJSON(url, callback) {
	load("GET", "json", null, url, callback)
}

export function initXHR() {
	let xhrItems = document.querySelectorAll("[data-xhr-url][data-xhr-target]");
	xhrItems.forEach(function(item){
		console.log(item.dataset.xhrUrl, item.dataset.xhrTarget)
		item.addEventListener("click", function (event) {
			getHTML(item.dataset.xhrUrl, item.dataset.xhrTarget)
		});
	});
}
