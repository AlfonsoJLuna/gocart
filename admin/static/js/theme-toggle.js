const root = document.documentElement;

function applyTheme(theme) {
    root.setAttribute("data-theme", theme);
}

function toggleTheme() {
    const current = root.getAttribute("data-theme") || "light";
    const next = current === "dark" ? "light" : "dark";
    localStorage.setItem("gocart-theme", next);
    applyTheme(next);
    const btn = document.getElementById("theme-btn");
    if (btn) btn.textContent = next === "dark" ? "Theme: Dark" : "Theme: Light";
}

document.addEventListener("DOMContentLoaded", function() {
    const btn = document.getElementById("theme-btn");
    if (btn) {
        btn.addEventListener("click", toggleTheme);
        btn.textContent = (localStorage.getItem("gocart-theme") || "light") === "dark" ? "Theme: Dark" : "Theme: Light";
    }
});
