var _$captchaImgSrc = "";
function freshCaptcha() {
    var $captchaImg = $("#captchaImg");
    if (!_$captchaImgSrc) {
        _$captchaImgSrc = $captchaImg.attr("src");
    }
    $captchaImg.attr("src", _$captchaImgSrc + "?time=" + new Date().getTime());
}

// ajaxable
$(function () {
    $("a[data-url]").click(function () {
        if (confirm("确认执行此操作？")) {
            $.post($(this).data("url"), function (result) {
                alert(JSON.stringify(result));
            }, "json");
        }
    });
})