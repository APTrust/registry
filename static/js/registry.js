// registry.js
import { initXHR, loadIntoElement } from './modules/xhr.js'
import { initModals } from './modules/modal.js'
import { initToggles } from './modules/toggle.js'
import { initSidebar } from './modules/sidebar.js'
import { initFiltersGrid } from './modules/filters-grid.js'
import { chartColors } from './modules/charts.js'

let APT = {}
APT.loadIntoElement = loadIntoElement
APT.chartColors = chartColors

window.addEventListener('load', (event) => {
    initXHR()
    initModals();
    initToggles();
    initSidebar();
    initFiltersGrid();
    window.APT = APT;

    // aptLoadEvent is defined in the head of the document.
    window.dispatchEvent(aptLoadEvent);
});
