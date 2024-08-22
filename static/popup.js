const SCROLL_THRESHOLD = 50;
const TRIGGER_THRESHOLD = 30000;

let alreadyPoppedUp = false;

function showModal() {
  fetch("/partials/pop-up-modal")
    .then((response) => {
      if (response.ok) {
        return response.text();
      } else {
        throw new Error("Error: " + response.statusText);
      }
    })
    .then((modalHTML) => {
      const modalOverlay = document.getElementById("modalOverlay");
      modalOverlay.outerHTML = modalHTML;
      modalOverlay.style.display = "flex";

      const modalContent = document.getElementById("modalContent");

      // Prevent clicks inside the modal content from closing the modal
      modalContent.addEventListener("click", (e) => {
        e.stopPropagation();
      });

      // Hide the modal when clicking outside the modal content
      const modal = document.getElementById("modalOverlay");
      modal.addEventListener("click", () => {
        modal.outerHTML = "";
        alreadyPoppedUp = true;
      });

      // Clicks redirect
      applyButtonlogic();
    })
    .catch(console.error);
}

// Trigger after time
window.addEventListener("load", () => {
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
