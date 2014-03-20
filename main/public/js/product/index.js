$(function () {
    var data = {};

    function setP(p) {
        if (!p || p == 0) {
            data.p = "";
            return
        }
        data.p = p || "";
    }

    function setCtcode(c) {
        data.ctcode = c || "";
    }

    function setQ(q) {
        data.q = q || "";
    }

    function reload() {
        var url = S(ProductIndexURL + "?ctcode={{ctcode}}&p={{p}}&q={{q}}").template(data).s;
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

    $("input[name=q]").on('keyup', function (e) {
        if (event.which == 13 || event.keyCode == 13) {
            e.preventDefault();
            setQ($(this).val());
            reload();
        }
    });

    setP(p);
    setCtcode((ctcode));
    setQ(q);
    $("input[name=q]").focus();
});

