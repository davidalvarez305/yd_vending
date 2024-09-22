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

function preserveQuerystring() {
	if (window.location.pathname.includes("crm")) return;

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

document.addEventListener("DOMContentLoaded", () => preserveQuerystring());