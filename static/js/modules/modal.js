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
        c.removeEventListener("click", modalOpen)
        c.addEventListener("click", modalOpen)
    });
}

function modalOpen(event) {
    let modal = document.getElementById(event.currentTarget.dataset.modal);
    event.preventDefault();
    document.body.classList.add("freeze");
    modal.classList.add("open");
}

// This is called in xhr.js after content is loaded into modal via ajax request.
export function attachModalClose(modal) {
    var parent = modal.parentElement
    var exits = modal.querySelectorAll(".modal-exit");
    console.log(`Found ${exits.length} close buttons`)
    exits.forEach(function (exit) {
        exit.addEventListener("click", function (event) {
            event.preventDefault();
            document.body.classList.remove("freeze");
            parent.classList.remove("open");
        });
        console.log("Added close listener to one button")
    });
}
