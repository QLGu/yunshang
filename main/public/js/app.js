var _$captchaImgSrc = "";
function freshCaptcha() {
    var $captchaImg = $("#captchaImg");
    if (!_$captchaImgSrc) {
        _$captchaImgSrc = $captchaImg.attr("src");
    }
    $captchaImg.attr("src", _$captchaImgSrc + "?time=" + new Date().getTime());
}

$(function () {
    function doAjaxLogin() {
        $login = $('input[name=login]');
        $password = $('input[name=password]');
        $.post(loginURL,
            {login: $login.val(), password: $password.val()},
            function (ret) {
                if (ret.ok) {
                    window.location.href = "/";
                } else {
                    window.location.href = "/passport/login";
                }
            },
            "json");
    }

    $('.log-btn').click(doAjaxLogin);

    $("div.log-form input[name=password]").on('keyup', function (e) {
        if (event.which == 13 || event.keyCode == 13) {
            e.preventDefault();
            doAjaxLogin();
        }
    });
});

$(function () {
    $.fn.lrkf && $("#lrkfwarp").lrkf({
        root: '/public/libs/lrkf/',
        skin: 'lrkf_green1',
        kfTop: '186',
        defShow: false,
        qqs: KF.qqs,
        tel: [
            {'name': '售前', 'tel':KF.sales_phone},
            {'name': '售后', 'tel':  KF.after_sales_phone}
        ],
        qrCode: '/public/img/qrcode.jpg',
        foot: "8:00-17:00"
    });
});
