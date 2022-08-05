// modals
//
// Attach events to modal controllers. This happens on page load
// and when we load new content into the DOM via xhr.
//
// Note that we track which elements have modal events attached using
// the attribute data-modal-initialized.
//

export function initModals() {
    var modalControllers = document.querySelectorAll("[data-modal]");
    modalControllers.forEach(function (c) {
        if (c.dataset.modalInitialized != "true") {
            c.addEventListener("click", function (event) {
                event.preventDefault();
                var modal = document.getElementById(c.dataset.modal);
                var body = document.body;
                body.classList.add("freeze");
                modal.classList.add("open");
                var exits = document.querySelectorAll(".modal-exit");
                exits.forEach(function (exit) {
                    exit.addEventListener("click", function (event) {
                        event.preventDefault();
                        body.classList.remove("freeze");
                        modal.classList.remove("open");
                    });
                });
            });
            c.dataset.modalInitialized = "true"
        }
    });
}
