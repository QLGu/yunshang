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
        qqs: [
            {'name': '售前陈R', 'qq': '2252410803'},
            {'name': '售前谢R', 'qq': '2930355581'},
            {'name': '售后王S', 'qq': '2252410803'},
            {'name': '售后李S', 'qq': '2252410803'}
        ],
        tel: [
            {'name': '售前', 'tel': '400-0686-198'},
            {'name': '售后', 'tel': '0755-2759786'}
        ],
        qrCode: '/public/libs/lrkf/qrcode.png',
        foot: "8:00-17:00"
    });
});
