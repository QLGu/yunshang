$(function () {
    $(".fancybox").fancybox({
            afterClose: function (e) {
                var onclose = $(".fancybox").data("on-close");
                if (onclose) {
                    eval(onclose);
                }
            }
        }
    );
});