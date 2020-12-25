let login_link = $('div.dropdown.boostbar-login a.boostbar-login');

if(login_link.length > 0) {
    var login_post_url = login_link[0].href;
    var login_post_data = "ReturnUrl=%252f&IsCheckout=False&Register=False&Email=" + arguments[0] + "&Password=" + arguments[1];
    $.post(login_post_url, login_post_data)
}
