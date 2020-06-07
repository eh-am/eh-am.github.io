function isAnchorToSamePage(anchor) {
  return anchor.hostname === window.location.hostname;
}

function getPageTransitioner() {
  return document.getElementById("page-transitioner");
}

// TODO: what happens if two links are clicked at the same time?
function addAnimationToAnchor(anchor) {
  // TODO: other events
  anchor.addEventListener("click", (event) => {
    const pt = getPageTransitioner();

    // delay going to new url
    event.preventDefault();

    // start animation
    pt.classList.add("start");

    // go to link after animation ends
    const listener = pt.addEventListener("animationend", () => {
      window.location = anchor.href;
      // we don't care about this listener anymore
      pt.removeEventListener("animationEnd", listener);
    });
  });
}

document.addEventListener("DOMContentLoaded", function () {
  const anchors = document.getElementsByTagName("a");

  Array.from(anchors).filter(isAnchorToSamePage).forEach(addAnimationToAnchor);
});
