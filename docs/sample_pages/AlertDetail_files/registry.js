// registry.js
import { initXHR, loadIntoElement } from './modules/xhr.js';
import { initModals } from './modules/modal.js';
import { initToggles } from './modules/toggle.js';

window.APT = {}
window.APT.loadIntoElement = loadIntoElement

window.addEventListener('load', (event) => {
    initXHR()
    initModals();
    initToggles();
});
