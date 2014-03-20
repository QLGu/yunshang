$(function () {
    var data = {p: p, ctcode: ctcode};

    function reload() {
        var url = S(ProductIndexURL + "?ctcode={{ctcode}}&p={{p}}").template(data).s;
        window.location.href = url;
    }

    $(".f_ct").click(function (event) {
        data.ctcode = $(this).data("code");

        reload();
    });

    $(".f_p").click(function (event) {
        data.p = $(this).data("id");

        reload();
    });

});