function onPage(p) {
    page.set("page", p);
    page.reload();
}

$(function () {
    $(".f_ct").live("click", function (event) {
        page.set("ctcode", $(this).data("code"));
        page.reload();
    });

    $(".f_p").live("click", function (event) {
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

    $('#btn-search').click(function () {
        var q = $("input[name=q]").val();
        page.set("q", q);
        page.reload();
    });


    $('a.box-link-expand').click(function () {
        var $this = $(this);
        if (hide_filters) {
            $this.html("收起刷选");
            hide_filters = "";
            $('div.prod-inner').show();
        } else {
            $this.html("显示刷选");
            hide_filters = "yes";
            $('div.prod-inner').hide();
        }
    });

    $('#btn-more-ps').click(function () {
        $('#curr-ps').hide();
        $('#more-ps').load(AllProvidersURL);
    });

    $('#btn-more-cs').click(function () {
        $('#curr-cs').hide();
        $('#more-cs').load(AllCategoriesURL);
    });

    $("input[name=q]").focus();
});

