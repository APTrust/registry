// xhr module

import { initModals, attachModalClose } from './modal.js';
import { initToggles } from './toggle.js';

//
// observer observes elements into which we dynamically load content
// via xhr requests. When the childList or characterData of these nodes
// change, the obeserver fires our callback, which ensures that xhr
// events are attached to newly-loaded content as necessary.
//
// Note that we track which elements have xhr events attached using
// the attribute data-xhr-initialized.
//
// https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver
// https://developer.mozilla.org/en-US/docs/Web/API/MutationObserver/observe
const callback = function(mutationsList, observer) {
    initXHR()
    initModals()
	initModalPosts()
    initToggles()
}
const observer = new MutationObserver(callback);

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
        return
    }
    observer.observe(element, {childList: true, characterData: true})
	return function(htmlData) {
        element.innerHTML = htmlData
		attachModalClose(element)
		focusOnChild(element)
		//console.log(document.activeElement)
	}
}

function deleteContent(elementId) {
	return function(textResponse) {
		document.getElementById(elementId).innerHTML = ""
		alert(textResponse)
	}
}

export function loadIntoElement(method, url, elementId) {
	load(method, "html", null, url, setContent(elementId))
}

function appendToElement(method, url, elementId) {
	load(method, "html", null, url, setContent(elementId))
}

function getJSON(url, callback) {
	load("GET", "json", null, url, callback)
}

function isElement(element) {
    return element instanceof Element || element instanceof Document;
}

function focusOnChild(element) {
	let child = element.querySelector('a:not(.modal-exit)', 'input', 'select', 'button')
	if (child != null) {
		child.focus()
	}
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
        if (item.dataset.xhrInitialized != "true") {
		    item.addEventListener("click", function (event) {
                if (fn == "append") {
			        appendToElement(method, url, target)
                } else {
                    loadIntoElement(method, url, target)
                }
		    });
            item.dataset.xhrInitialized = "true"
        }
	});

    let modalItems = document.querySelectorAll("[data-modal][data-xhr-url]");
	modalItems.forEach(function(item){
		//console.log(item.dataset.xhrUrl, item.dataset.xhrTarget)
        let modal = document.getElementById(item.dataset.modal)
        let modalContentDiv = modal.querySelector('.modal-container');
        let method = item.dataset.xhrMethod || "get"
        if (item.dataset.initialized != "true") {
		    item.addEventListener("click", function (event) {
			    loadIntoElement(method, item.dataset.xhrUrl, modalContentDiv)
		    });
            item.dataset.initialized = "true"
        }
	});
}

// Post a form via XHR and load the result into modal with targetId.
// This does assume that the target is a modal.
export function modalPost(formName, modalId){
	var form = document.forms[formName]
	var data = new FormData(form)

	function resultToModal() { 
		let modal = document.getElementById(modalId)
		let modalContentDiv = modal.querySelector('.modal-container')
		modalContentDiv.innerHTML = xhr.responseText 
		document.body.classList.add("freeze");
		modal.classList.add("open")
		attachModalClose(modal)
	}

	var xhr = new XMLHttpRequest()
	xhr.open(form.method, form.action)
	xhr.onload = resultToModal
	xhr.onerror = resultToModal
	xhr.send(data); 
}

function initModalPosts() {
	let modalPostItems = document.querySelectorAll("[data-modal-post-form][data-modal-post-target]");
	modalPostItems.forEach(function(item){
		let formName = item.dataset.modalPostForm
		let modalId = item.dataset.modalPostTarget
		item.addEventListener('click', () => { modalPost(formName, modalId) })
	});
}

