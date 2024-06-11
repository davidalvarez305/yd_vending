const modalOverlay = document.getElementById("modalOverlay");
const modalContent = modalOverlay.querySelector('.bg-white');
let alreadyPoppedUp = false;

const scrollPercentageThreshold = 50;

function showModal() {
  modalOverlay.style.display = "flex";

   setTimeout(function () {
    modalOverlay.style.display = 'none';
    alreadyPoppedUp = true;
}, 3000);
}

// Hide the modal when clicking outside the modal content
modalOverlay.addEventListener('click', function (event) {
  modalOverlay.style.display = 'none';
  alreadyPoppedUp = true;
});

// Prevent clicks inside the modal content from closing the modal
modalContent.addEventListener('click', function (event) {
  event.stopPropagation();
});

window.addEventListener("scroll", () => {
  const scrollY = window.scrollY || window.pageYOffset;
  const windowHeight = window.innerHeight;
  const totalHeight = document.body.clientHeight;
  const scrolledPercentage = (scrollY / (totalHeight - windowHeight)) * 100;

  if (scrolledPercentage >= scrollPercentageThreshold && !alreadyPoppedUp) {
    showModal();
  }
});