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

$(function () {
    // Ajax Global Config
    $(document).ajaxSend(function (_event, _jqXHR, ajaxOptions) {
        if (!_isDevAjax(ajaxOptions)) { // not dev ajax
            console.log("ajaxSend");
            $("#ajax-preload").show();
        }
    });
    $(document).ajaxComplete(function (_event, _XMLHttpRequest, ajaxOptions) {
        if (!_isDevAjax(ajaxOptions)) {
            console.log("ajaxComplete");
            $("#ajax-preload").hide();
        }
    });

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
});
