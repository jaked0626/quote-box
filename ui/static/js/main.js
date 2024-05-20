"use strict";

// keep track of which tab is open
let navLinks = document.querySelectorAll("nav a");
for (let link of navLinks) {
  if (link.getAttribute("href") == window.location.pathname) {
    link.classList.add("live");
    break;
  } else {
    console.log(link.getAttribute("href"), window.location.pathname);
  }
}
