"use strict";

// keep track of which tab is open
let navLinks = document.querySelectorAll("nav a");
for (let link of navLinks) {
  // the following is very buggy and just plain incorrect... I would either need to pass the authenticated user's id 
  // to the template data or change the paths completely
  if (link.getAttribute("href") == "/snippet/mine" && window.location.pathname.includes("snippet/user")) {
    link.classList.add("live");
    break;
  }
  else if (link.getAttribute("href") == window.location.pathname) {
    link.classList.add("live");
    break;
  } else {
    console.log(link.getAttribute("href"), window.location.pathname);
  }
}
