export function initModals() {
    var modalControllers = document.querySelectorAll("[data-modal]");
    modalControllers.forEach(function (modalController) {
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
