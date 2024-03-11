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
    // This is a hack, but we have to wait for children to load,
    // and there's no way to attach an onload event to an item
    // that isn't present yet. This does not work without the timeout.
    // Ideally, we use a mutation observer for this. We'll get back
    // to this when we actually have some time.
    window.setTimeout(function() { 
        console.log("Attaching modal close listener.")
        attachModalClose(modal) 
    }, 350)
}

export function attachModalClose(modal) {
    var exits = modal.querySelectorAll(".modal-exit");
    console.log(`Found ${exits.length} close buttons`)
    exits.forEach(function (exit) {
        exit.addEventListener("click", function (event) {
            event.preventDefault();
            document.body.classList.remove("freeze");
            modal.classList.remove("open");
        });
        console.log("Added close listener to one button")
    });
}
