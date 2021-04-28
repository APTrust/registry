// toggle.js
//
// Attach toggle behavior to items with data-toggle attribute.
// This happens on page load and when we load new content into
// the DOM via xhr.
//
// Note that we track which elements have modal events attached using
// the attribute data-toggle-initialized.
//

export function initToggles() {
    var controllers = document.querySelectorAll("[data-toggle]");
    controllers.forEach(function (c) {
        if (c.dataset.toggleInitialized != "true") {
            c.addEventListener("click", function (event) {
                event.preventDefault();
                var target = document.getElementById(c.dataset.toggle);
                target.style.display == "block" ? target.style.display = "none" : target.style.display = "block"
            });
            c.dataset.toggleInitialized = "true"
        }
    });

}