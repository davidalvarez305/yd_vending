const SCROLL_THRESHOLD = 50;
const TIME_THRESHOLD = 3000;
const TRIGGER_THRESHOLD = 2000;

let alreadyPoppedUp = false;

const modalOverlay = document.getElementById("modalOverlay");

function showModal() {
  fetch("/pop-up-modal")
    .then((response) => {
      if (response.ok) {
        return response.text();
      } else {
        throw new Error("Error: " + response.statusText);
      }
    })
    .then((modal) => {
      modalOverlay.outerHTML = modal;
      modalOverlay.style.display = "flex";

      const modalContent = modalOverlay.querySelector(".bg-white");

      // Prevent clicks inside the modal content from closing the modal
      modalContent.addEventListener("click", (e) => {
        e.stopPropagation();
      });

      // Hide the modal when clicking outside the modal content
      modalOverlay.addEventListener("click", () => {
        modalOverlay.style.display = "none";
        alreadyPoppedUp = true;
      });
    });
}

// Trigger after time
window.addEventListener("DOMContentLoaded", () => {
  setTimeout(() => {
    if (!alreadyPoppedUp) showModal();
  }, TRIGGER_THRESHOLD);
});

// Trigger after scroll
window.addEventListener("scroll", () => {
  const scrollY = window.scrollY || window.pageYOffset;
  const windowHeight = window.innerHeight;
  const totalHeight = document.body.clientHeight;
  const scrolledPercentage = (scrollY / (totalHeight - windowHeight)) * 100;

  if (scrolledPercentage >= SCROLL_THRESHOLD && !alreadyPoppedUp) {
    showModal();
  }
});

// Hide after time
setTimeout(() => {
  modalOverlay.style.display = "none";
  alreadyPoppedUp = true;
}, TIME_THRESHOLD);
