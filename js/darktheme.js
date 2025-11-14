let current = document.documentElement.getAttribute('data-theme');
const toggleDarkMode = () => {
    current = current === 'dark' ? 'light' : 'dark';
    document.documentElement.setAttribute('data-theme', current);
    localStorage.theme = current;
}

let preferredDark = (!("theme" in localStorage) && window.matchMedia("(prefers-color-scheme: dark)").matches);
if (localStorage.theme === "dark" || preferredDark) {
    localStorage.theme = "dark";
    document.documentElement.setAttribute('data-theme', "dark")
    current = "dark";
} else {
    localStorage.theme = "light";
    document.documentElement.setAttribute('data-theme', "light")
    current = "light";
}

document.toggleDarkMode = toggleDarkMode;
