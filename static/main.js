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

function preserveQuerystring() {
	if (["crm", "inventory", "reports"].some(link => window.location.pathname.includes(link))) return;

	const links = document.querySelectorAll('a');

	links.forEach(link => {
		link.addEventListener('click', function () {
			const qs = window.location.search;
			if (qs.length > 0) {
				const url = new URL(this.href, window.location.origin);
				url.search = url.search ? `${url.search}&${qs.substring(1)}` : qs;
				this.href = url.toString();
			}
		});
	});
}

document.addEventListener("DOMContentLoaded", () => {
	// Preserve querystring
	preserveQuerystring();

	// Get user variables
	localStorage.getItem("user") ?? setUser();

	// Pop quote form when CTA's are clicked
	let quoteButtons = document.querySelectorAll(".quoteButton");

	quoteButtons.forEach((button) => {
		button.addEventListener("click", function () {
			const formModal = document.getElementById("formModalContainer");
			if (formModal) formModal.style.display = "";

			const buttonClicked = document.getElementById("button_clicked");
			if (buttonClicked) buttonClicked.value = button.getAttribute("name");

			const popUp = document.getElementById("popUpModalOverlay");
			if (popUp) popUp.style.display = "none";
		});
	});
});