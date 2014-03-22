function onPage(p){
    page.set("page", p);
    page.reload();
}

$(function () {
    $(".f_ct").click(function (event) {
        page.set("ctcode", $(this).data("code"));
        page.reload();
    });

    $(".f_p").click(function (event) {
        page.set("p", $(this).data("id"));
        page.reload();
    });

    $("input[name=q]").on('keyup', function (e) {
        if (event.which == 13 || event.keyCode == 13) {
            e.preventDefault();
            page.set("q", $(this).val());
            page.reload();
        }
    });

    $("#btn-reset-search").click(function () {
        page.clear();
        page.reload();
    });

     $("input[name=q]").focus();
});

