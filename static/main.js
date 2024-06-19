document.addEventListener("DOMContentLoaded", () => {
  const phoneButtons = document.querySelectorAll(".phoneNumber");

  phoneButtons.forEach(telephoneButton => {
    const phoneNumber = telephoneButton.textContent.trim();
  
    telephoneButton.textContent = formatPhoneNumber(phoneNumber);
  
    telephoneButton.addEventListener("click", () => {
      window.location.href = "tel:" + phoneNumber;
    });
  });

  function formatPhoneNumber(phoneNumber) {
    const cleaned = ("" + phoneNumber).replace(/\D/g, "");
    const match = cleaned.match(/^(\d{3})(\d{3})(\d{4})$/);
    if (match) {
      return "(" + match[1] + ") " + match[2] + " - " + match[3];
    }

    return null;
  }
});
