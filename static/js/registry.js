// registry.js
import { initXHR } from './modules/xhr.js';
import { initModals } from './modules/modal.js';

function initToggles() {
    var controllers = document.querySelectorAll("[data-toggle]");
    controllers.forEach(function (c) {
        c.addEventListener("click", function (event) {
            event.preventDefault();
            var target = document.getElementById(c.dataset.toggle);
            target.style.display == "block" ? target.style.display = "none" : target.style.display = "block"
        });
    });

}

window.addEventListener('load', (event) => {
    initXHR()
    initModals();
    initToggles();
});
