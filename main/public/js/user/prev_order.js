$(function () {
    var ractive = new Ractive({
        el: "items",
        template: "#items_tpl",
        data: {
            "order": {},
            "items": [],
            "submitable": function(){
              return false;
            },
            "fixed": function (v) {
                return accounting.toFixed(v, 2);
            }
        }
    });

    ractive.on({
        "load": function () {
            $.getJSON(orderURL + "?code=" + orderCode, function (ret) {
                var order = ret.data;
                ractive.set("order", order);
            });
            $.getJSON(orderItemsURL + "?code=" + orderCode, function (ret) {
                var items = ret.data;
                ractive.set("items", items);
                ractive.fire("load-products");
            });
        },
        "load-products": function () {
            $.getJSON(orderProductsURL + "?code=" + orderCode, function (ret) {
                var ps = ret.data;
                var items = ractive.get("items");
                _.each(ps, function (p, i) {
                    var c = _.find(items, function (item) {
                        return item.product_id == p.id;
                    });
                    c.name = p.name;
                    c.model = p.model;
                    c["url"] = "/products/p/" + p.id;
                    c["purl"] = "/product/image?file=" + p.id + ".jpg";
                });
                ractive.set("items", items);
            });
        },
        "load-das": function () {
            $('#das').load(dasForSelect);
        }
    });

    ractive.fire("load");
    ractive.fire("load-das");

    $('#btn-refresh-das').live("click", function () {
        ractive.fire("load-das");
    });
});
