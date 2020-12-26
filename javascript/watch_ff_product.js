// get freefall price element
let freefall_price = document.querySelector('span.lot-ff span.price-amount')
if(freefall_price.length > 0) {
    var observer = new MutationObserver(function(mutations) {
        mutations.forEach(function(mutation) {
            if (mutation.type == "attributes") {
                console.log("attributes changed")
            }
        });
    });

    observer.observe(freefall_price[0], {
        attributes: true //configure it to listen to attribute changes
    });
}

jQuery.connection['auctionClientHub'].client.lotTick = function(r) {
    console.log("I did this: lotId: " + r.l + ", price: " + r.p + ", dutchInstanceId: " + r.d + ", status: " + r.s);
    1 || WTF.trigger("lot-tick", jQuery.extend({
        lot: {
            lotId: r.l,
            price: r.p,
            dutchInstanceId: r.d,
            status: r.s
        }
    }, "auctionClientHub"))
}

WTF.auctioneer.on("lot-tick", function(i, r) {
    console.log("I did this: lotId: " + r.lot.lotId + ", price: " + r.lot.price);
    var u = jQuery(".lot-tick-" + r.lot.lotId + " .price-amount").each(function() {
        switch (r.lot.status) {
        case "t":
            WTF.ui.dollars(this, r.lot.price);
            break;
        case "r":
            jQuery(this).html("Resetting&hellip;")
        }
    }).length;
    u === 0 && WTF.auctioneer.removeLot(r.lot.lotId)
})