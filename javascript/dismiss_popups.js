// disable those stupid OneSignal notifications!
// just remove the whole dang reference
window.OneSignal = [];

waitForElementToDisplay('div#sweepstakes-popup .close', function(){
    let sweepstakes_popup_close = $('div#sweepstakes-popup .close');
    if(sweepstakes_popup_close.length > 0) {
        sweepstakes_popup_close[0].click();
    }
}, 500, 5000);

waitForElementToDisplay('div#onesignal-slidedown-container button#onesignal-slidedown-allow-button', function(){
    let onesignal_accept_button = $('div#onesignal-slidedown-container button#onesignal-slidedown-allow-button');
    if(onesignal_accept_button.length > 0) {
        onesignal_accept_button.click();
    }
}, 500, 5000);

function waitForElementToDisplay(selector, callback, checkFrequencyInMs, timeoutInMs) {
  var startTimeInMs = Date.now();
  (function loopSearch() {
    if (document.querySelector(selector) != null) {
      callback();
      return;
    }
    else {
      setTimeout(function () {
        if (timeoutInMs && Date.now() - startTimeInMs > timeoutInMs)
          return;
        loopSearch();
      }, checkFrequencyInMs);
    }
  })();
}