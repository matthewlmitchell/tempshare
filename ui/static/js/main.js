var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

const switchTheme = document.querySelector("#switch");

switchTheme.addEventListener("click", () => {
	document.body.classList.toggle("dark-mode");

	if(document.body.classList.contains("dark-mode")) {
		localStorage.setItem("dark-mode", "enabled");
	} else {
		localStorage.setItem("dark-mode", "disabled");
	}
});

if(localStorage.getItem("dark-mode") == "enabled") {
	document.body.classList.toggle("dark-mode");
	document.getElementById("switch").checked = true;
}