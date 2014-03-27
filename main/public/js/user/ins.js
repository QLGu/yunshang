$(function () {
    $("#das_table_new").click(function () {
        $.fancybox.open({
            href: addNewInvoiceUrl,
            type: 'iframe',
            padding: 5,
            afterClose: function (e) {
               window.location.reload();
            }
        });
    });

    $(".fancybox").fancybox({
            afterClose: function (e) {
                var onclose = $(".fancybox").data("on-close");
                if (onclose){
                    eval(onclose);
                }
            }
        }
    );
});
