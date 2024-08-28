function handleCloseAlertModal() {
  const closeModalButton = document.querySelectorAll(".closeModal");
  const newModal = document.getElementById("alertModal");

  closeModalButton.forEach((button) => {
    let children = button.children;
    Array.from(children).forEach((child) => {
      child.setAttribute("name", button.name);
    });

    button.addEventListener("click", () => (newModal.style.display = "none"));
  });
}

function setUser() {
  const user = {
    landingPage: window.location.href,
    referrer: document.referrer,
  };

  localStorage.setItem("user", JSON.stringify(user));
}

document.addEventListener("DOMContentLoaded", () => localStorage.getItem("user") ?? setUser());