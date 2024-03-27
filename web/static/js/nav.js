let navOpen = false
handleDrawer = () => {
    if (!navOpen) {
        navDrawer.classList.remove("-translate-x-full")
        navDrawer.classList.add("translate-x-0")
    } else {
        navDrawer.classList.remove("translate-x-0")
        navDrawer.classList.add("-translate-x-full")
    }
    navOpen = !navOpen
}

let navButton = document.getElementById("nav-button")
let navDrawer = document.getElementById("nav-drawer")
if (window.innerWidth > 768) {
    navDrawer.classList.remove("h-0")
    navDrawer.classList.add("h-auto")
}
navButton.addEventListener("click", handleDrawer)
navDrawer.addEventListener("click", handleDrawer)
