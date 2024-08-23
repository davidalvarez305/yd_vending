const form = document.getElementById("get-a-quote-form");
const alertModal = document.getElementById("alertModal");

form.addEventListener("submit", e => {
    e.preventDefault();
    const isValid = validateForm();

    if (!isValid) return;

    const marketing = Object.fromEntries(new URLSearchParams(window.location.search));
    const data = new FormData(e.target);

    for (const [key, value] of Object.entries(marketing)) {
        data.append(key, value);
    }

    fetch("/quote", {
        method: "POST",
        credentials: "include",
        body: data
    })
        .then((response) => {
            if (response.ok) {
                const token = response.headers.get('X-Csrf-Token');

                if (token) {
                    const csrf_token = document.getElementById('csrf_token');

                    if (!csrf_token) return;

                    csrf_token.value = token;
                }

                return response.text();
            } else {
                return response.text().then((err) => {
                    throw new Error(err);
                });
            }
        })
        .then(html => {
            alertModal.outerHTML = html;
            const modal = document.getElementById("alertModal");
            handleCloseAlertModal();

            form.reset();
        })
        .catch(errorHTML => {
            alertModal.outerHTML = errorHTML
            handleCloseAlertModal();
        });
});

function validatePhoneNumber(phoneNumberInput) {
    let isValid = true;
    const phoneNumber = phoneNumberInput.trim();

    const phonePattern = /^[0-9]{10}$/;

    if (!phonePattern.test(phoneNumber)) {
        isValid = false;
    }

    return isValid;
}

function validateForm() {
    let isValid = true;

    // Validate first name
    const firstNameInput = document.getElementById("first_name");
    if (!firstNameInput.value.trim()) {
        isValid = false;
        firstNameInput.classList.add("border-red-500");
    } else {
        firstNameInput.classList.remove("border-red-500");
    }

    // Validate last name
    const lastNameInput = document.getElementById("last_name");
    if (!lastNameInput.value.trim()) {
        isValid = false;
        lastNameInput.classList.add("border-red-500");
    } else {
        lastNameInput.classList.remove("border-red-500");
    }

    // Validate phone number
    const phoneNumberInput = document.getElementById("phone_number");
    if (!validatePhoneNumber(phoneNumberInput.value)) {
        isValid = false;
        phoneNumberInput.classList.add("border-red-500");
    } else {
        phoneNumberInput.classList.remove("border-red-500");
    }

    // Validate machine type
    const machineTypeSelect = document.getElementById("machine_type");
    if (machineTypeSelect.value === "") {
        isValid = false;
        machineTypeSelect.classList.add("border-red-500");
    } else {
        machineTypeSelect.classList.remove("border-red-500");
    }

    // Validate location type
    const locationTypeSelect = document.getElementById("location_type");
    if (locationTypeSelect.value === "") {
        isValid = false;
        locationTypeSelect.classList.add("border-red-500");
    } else {
        locationTypeSelect.classList.remove("border-red-500");
    }

    // Validate city
    const citySelect = document.getElementById("city");
    if (citySelect.value === "") {
        isValid = false;
        citySelect.classList.add("border-red-500");
    } else {
        citySelect.classList.remove("border-red-500");
    }

    if (!isValid) {
        alert("Please fill in all required fields.");
    }

    return isValid;
}