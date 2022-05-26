// filters-grid.js
//
// Handles the clicking of "More Filters" in the filters grids
//
//

export function initFiltersGrid() {
    var gridFilters = document.querySelector('.filters-grid')
    var gridFiltersAll = document.getElementById('gridFiltersAll');
    var gridFiltersToggle = document.querySelector('.filters-grid .filter-toggle');
    var gridFiltersMore = document.querySelector('.filters-grid .more-filters')
    var gridFiltersLess = document.querySelector('.filters-grid .less-filters')

    gridFiltersToggle.addEventListener('click', function(event) {
        event.preventDefault();
        gridFilters.classList.toggle('is-open')
        gridFiltersAll.classList.toggle('is-sr-only');
        gridFiltersMore.classList.toggle('is-hidden');
        gridFiltersLess.classList.toggle('is-hidden');
    });
}
