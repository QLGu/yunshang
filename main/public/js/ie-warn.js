if ($.browser.msie) {
    if ($.browser.version < 8.0) {
        if (!$.cookie('ie-warn')) {
            alert("检测您在使用低版本的IE浏览器访问本站。 为了提升您的使用体验，请使用IE8或以上版本，谢谢支持！");
            $.cookie('ie-warn', true, {expires: 1, path: '/' });
        }
    }
}