const telephoneButton = document.getElementById("telephone-button");

telephoneButton.addEventListener("click", () => {
  const phoneNumber = telephoneButton.textContent.trim();
  window.location.href = "tel:" + phoneNumber;
});
