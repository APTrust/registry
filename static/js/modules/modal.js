export function initModals() {
    // Get all buttons/links that have a data-modal attibute
    var modals = document.querySelectorAll("[data-modal]");

    // Attach callbacks to the modal controllers.
    modals.forEach(function (modalController) {
        modalController.addEventListener("click", function (event) {
            event.preventDefault();
            var modal = document.getElementById(modalController.dataset.modal);
            modal.classList.add("open");
            var exits = modal.querySelectorAll(".modal-exit");
            exits.forEach(function (exit) {
                exit.addEventListener("click", function (event) {
                    event.preventDefault();
                    modal.classList.remove("open");
                });
            });
        });
    });
}
