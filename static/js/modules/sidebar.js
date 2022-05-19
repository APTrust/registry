// sidebar.js
//
// Handles the clicking of "More Options" in the global sidebar
//
//

export function initSidebar() {
    var sidebarToggle = document.getElementById('sidebarToggle');
    var sidebarSubnav = document.getElementById('sidebarSubnav');
    var sidebarMore = document.querySelector('.more-options')
    var sidebarLess = document.querySelector('.less-options')

    sidebarToggle.addEventListener('click', function(event) {
        event.preventDefault();
        sidebarSubnav.classList.toggle('is-sr-only');
        sidebarMore.classList.toggle('is-hidden');
        sidebarLess.classList.toggle('is-hidden');
    });
}
