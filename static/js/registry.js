// registry.js
import { initXHR } from './modules/xhr.js';
import { initModals } from './modules/modal.js';

window.addEventListener('load', (event) => {
    console.log("Initializing...")
    initXHR()
    initModals();
});
