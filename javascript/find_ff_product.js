// Navigates to and watches the price of a freefall item
// Arguments:
// [0] = Link to product (example: "https://www.cigarbid.com/a/cao-brazilia-lambada/3783948/")
let freefall_link = arguments[0];
if (!(freefall_link.startsWith("https://www.cigarbid.com/"))) {
    freefall_link = freefall_link.replace() // TODO
    // TODO : dynamically correct the URL if it's not formatted correctly. Return error message if unable to correct
    freefall_link = "https://www.cigarbid.com/" +  freefall_link;
}

if(freefall_link.length > 0) {
    window.location.href = freefall_link[0].href;
}