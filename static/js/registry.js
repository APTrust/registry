// registry.js
import { initXHR } from './modules/xhr.js';
import { initModals } from './modules/modal.js';
import { initToggles } from './modules/toggle.js';


window.addEventListener('load', (event) => {
    initXHR()
    initModals();
    initToggles();
});
