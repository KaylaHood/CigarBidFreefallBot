// Navigates to and watches the price of a freefall item
// Arguments:
// [0] = Link to product (example: "https://www.cigarbid.com/a/cao-brazilia-lambada/3783948/")
let freefall_link = arguments[0];
if (!(freefall_link.match(/^https:\/\/www\.cigarbid\.com\/a\/[^\/]+\/[\d]+\/$/))) {
    return ("URL didn't match regex, this was the url provided: " + freefall_link);
}
window.location.href = freefall_link;