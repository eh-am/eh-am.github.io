function isAnchorToSamePage(anchor) {
  return anchor.hostname === window.location.hostname;
}

function getPageTransitioner() {
  return document.getElementById("page-transitioner");
}

function calculateDesiredDiamater() {
  // motivation:
  // * for the circle animation to cover the whole screen,
  // its radius needs to be the same as the hypothenuse
  // * we can only set the diamater, due to css circle contraints
  // * radius = diameter/2
  // starting with the pythagorean theorem:
  // hypotenuse^2 = width^2 + height^2
  // hypotenuse = radius = diameter/2
  // (diam/2)^2 = w^2 + h^2
  // d = 2 * sqrt(w^2 + h^2)
  return 2 * Math.sqrt(Math.pow(window.innerHeight,2) + Math.pow(window.innerWidth, 2));
}

// TODO: what happens if two links are clicked at the same time?
function addAnimationToAnchor(anchor) {
  // TODO: other events
  anchor.addEventListener("click", (event) => {
    // delay going to new url
    event.preventDefault();

    const pt = getPageTransitioner();
    const desiredDiameter = calculateDesiredDiamater();

    // start animation
    pt.style.setProperty("--pageTransitionDiameter", desiredDiameter);
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
