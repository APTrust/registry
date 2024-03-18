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

        // The following is for aria keyboard navigation. We don't want
        // assistive technologies to tab into a hidden sub-menu. Setting
        // both aria-expanded and display=block/hidden gives us the 
        // correct behavior. Tabbing skips closed menu items but works
        // on open ones. Part of the general accessibility card at
        // https://trello.com/c/CEQ5jAe1
        var subnavIsOpen = sidebarMore.classList.contains("is-hidden")
        if (subnavIsOpen) {
            sidebarToggle.setAttribute("aria-expanded", "true")
            sidebarSubnav.style.display = "block"
        } else {
            sidebarToggle.setAttribute("aria-expanded", "false")
            sidebarSubnav.style.display = "none"
        }
    });
}
