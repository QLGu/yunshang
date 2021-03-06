$(function () {
    var ractive = new Ractive({
        el: "carts",
        template: "#carts_tpl",
        data: {
            "carts": [],
            "fixed": function (v) {
                return accounting.toFixed(v, 2);
            },
            "sum": function (carts) {
                var ps = _.map(carts, function (c) {
                    return c.pref_price * c.quantity;
                });
                return _.reduce(ps, function (sum, num) {
                    return sum + num;
                });
            },
            "fname": function (index, n) {
                return "ps[" + index + "]." + n;
            }
        }
    });

    ractive.on({
        "load": function () {
            $.getJSON(cartDataURL, function (ret) {
                var carts = ret.data;
                _.each(carts, function (v, i) {
                    v["url"] = "/products/p/" + v.product_id;
                    v["purl"] = "/product/image?file=" + v.product_id + ".jpg";
                });
                ractive.set("carts", carts);
                ractive.fire("load-prices");
            });
        },
        "load-prices": function () {
            $.getJSON(cartProductPricesURL, function (ret) {
                var ps = ret.data;
                var carts = ractive.get("carts");
                _.each(ps, function (p, i) {
                    var c = _.find(carts, function (c) {
                        return c.product_id == p.id;
                    });
                    c.price = p.price;
                    c.min = p.min_number_of_orders;
                    c.stock_number = p.stock_number;
                    c.name = p.name;
                    c.model = p.model;
                    c.enabled = p.enabled;
                    c.message = "库存：" + c.stock_number + " 起订数量:" + c.min;
                });
                ractive.set("carts", carts);
                ractive.fire("load-pref-prices");
            });
        },
        "load-pref-prices": function () {
            var carts = ractive.get("carts");
            _.each(carts, function (c, i) {
                ractive.fire("load-pref-price", c, carts);
            });
        },
        "load-pref-price": function (c, carts) {
            $.getJSON(cartProductPrefPriceURL, {productId: c.product_id, quantity: c.quantity}, function (ret) {
                c.pref_price = ret.data;
                ractive.set("carts", carts);
            });
        },
        "do-inc": function (id, q, mode) {
            if (q == "") {
                return;
            }
            var q = parseInt(q);
            if (!q) {
                alert("请输入合法的数量");
                return;
            }
            q = (mode == "inc" ? q : -q);
            $.ajax({url: cartProductIncQuantityURL,
                data: {id: id, quantity: q},
                dataType: "json",
                success: function (ret) {
                    if (ret.ok) {
                        var carts = ractive.get("carts");
                        var c = _.find(carts, function (c) {
                            return c.id == id;
                        });
                        c.quantity = ret.data.quantity;
                        ractive.set("carts", carts);
                        var c = _.find(carts, function (c) {
                            return c.id == id;
                        });
                        ractive.fire("load-pref-price", c, carts);
                    } else {
                        alert(ret.message);
                    }
                }
            });
        },
        "inc": function (event, id) {
            var q = prompt("增加数量:", "1");
            ractive.fire("do-inc", id, q, "inc");
        },
        "dec": function (event, id) {
            var q = prompt("减少数量:", "1");
            ractive.fire("do-inc", id, q, "dec");
        },
        "delete": function (event, id) {
            if(confirm($(event.node).data("confirm"))){
                $.post(deleteCartProductURL, {id: id}, function (ret) {
                    window.location.reload();
                },"json");
            }
        },
        "toggle-check": function (e) {
            $('.checkbox:checkbox').prop('checked', e.node.checked);
        },
        "clean-all": function (e) {
            doAjaxPost(cleanCartURL, function (ret) {
                window.location.reload();
            });
        },
        "move-to-collects": function (e) {
            doAjaxPost(moveToCollectsURL, function (ret) {
                window.location.reload();
            });
        }
    });

    ractive.fire("load");
});
