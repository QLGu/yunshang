$(function () {
    var ractive = new Ractive({
        el: "products",
        template: "#products_tpl",
        data: {
            products: []
        }
    });
    ractive.on({
        "load": function () {
            $.getJSON(dataURL, function (ret) {
                ractive.set("products", ret.aaData);
            });
        }
    });
    ractive.fire("load");
});