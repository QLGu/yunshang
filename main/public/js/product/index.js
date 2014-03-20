$(function () {
    var data = {};

    function setP(p) {
        if (!p|| p == 0){
            data.p = "";
            return
        }
        data.p = p || "";
    }

    function setCtcode(c) {
        data.ctcode = c || "";
    }

    setP(p);
    setCtcode((ctcode));

    function reload() {
        var url = S(ProductIndexURL + "?ctcode={{ctcode}}&p={{p}}").template(data).s;
        window.location.href = url;
    }

    $(".f_ct").click(function (event) {
        setCtcode($(this).data("code"));
        reload();
    });

    $(".f_p").click(function (event) {
        setP($(this).data("id"));
        reload();
    });
});