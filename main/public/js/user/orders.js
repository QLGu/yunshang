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
    function onPage(p){
        page.set("page", p);
        page.reload();
    }
    $("#e1").select2({  placeholder: "选择订单状态"  });

    $('#e1').change(function(){
        page.set("filter_status", $(this).val());
        page.reload();
    });
});