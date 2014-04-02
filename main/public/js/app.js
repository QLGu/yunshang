Ractive.delimiters = [ '[[', ']]' ];
Ractive.tripleDelimiters = [ '[[[', ']]]' ];

var _$captchaImgSrc = "";
function freshCaptcha() {
    var $captchaImg = $("#captchaImg");
    if (!_$captchaImgSrc) {
        _$captchaImgSrc = $captchaImg.attr("src");
    }
    $captchaImg.attr("src", _$captchaImgSrc + "?time=" + new Date().getTime());
}

function _isDevAjax(ajaxOptions) {
    return ajaxOptions["url"].indexOf("/@reveltang_dev/__source_changed") != -1;
}

function reloadWindow() {
    window.location.reload();
}

function doAjaxPost(url, data, after) {
    if (_.isFunction(data)) {
        after = data;
        data = {};
    }
    if (confirm("确认执行此操作？")) {
        $.post(url, data, function (result) {
            if (result.ok) {
                alert("操作成功：" + result.message);
                if (after) {
                    after(result)
                }
            } else {
                alert("操作失败：" + result.message);
            }
        }, "json");
    }
}

$(function () {
    // Ajax Global Config
    $(document).ajaxSend(function (_event, _jqXHR, ajaxOptions) {
        if (!_isDevAjax(ajaxOptions)) { // not dev ajax
            if ($("#loadingbar").length === 0) {
                $("body").append("<div id='loadingbar'></div>")
                $("#loadingbar").addClass("waiting").append($("<dt/><dd/>"));

                $("#loadingbar").width((50 + Math.random() * 30) + "%");
            }
        }
    });
    $.ajaxPrefilter(function (opt, origOpt, jqxhr) {
        jqxhr.always(function () {
            $("#loadingbar").width("101%").delay(200).fadeOut(400, function () {
                $(this).remove();
            });
            //$("[data-plugin]").plugin();
        });
    });
    /*
     $(document).ajaxComplete(function (_event, _XMLHttpRequest, ajaxOptions) {
     if (!_isDevAjax(ajaxOptions)) {
     $("#loadingbar").width("101%").delay(200).fadeOut(400, function () {
     $(this).remove();
     });
     }
     });*/

    // Ajaxable Links
    $("a[data-url]").click(function () {
        if (confirm("确认执行此操作？")) {
            var $this = $(this)
            $.post($(this).data("url"), function (result) {
                if (result.ok) {
                    alert("操作成功：" + result.message);
                    var after = $this.data("after");
                    if (after) {
                        eval(after);
                    }
                } else {
                    alert("操作失败：" + result.message);
                }
            }, "json");
        }
    });

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
