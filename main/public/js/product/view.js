$(function () {
    var ractive = new Ractive({
        el: "images",
        template: "#images_tpl",
        data: {
            images: []
        }
    });

    ractive.on({
        "load": function () {
            $.getJSON(ProductSdImagesUrl, function (ret) {
                var images = _.map(ret.data, function (v, i) {
                    v.url = ImageSdUrl + "?file=" + v.value + "&time=" + new Date().getTime();
                    return v;
                });
                ractive.IMAGES = images;

                if (ractive.IMAGES && ractive.IMAGES.length > 0) {
                    var current = images[0];
                    current.active = "active";
                    ractive.set("current", current);

                    ractive.fire("set-images", 0);
                }
            });
        },
        "click": function (e, index) {
            _.each(ractive.IMAGES, function (v) {
                v.active = "";
            })
            var images = ractive.get("images");
            var it = images[index];
            it.active = "active";
            ractive.set("images", images);
            ractive.set("current", it);
        },
        "prev": function (e) {
            if (ractive.START > 0) {
                ractive.fire("set-images", ractive.START - 1);
            }
        },
        "next": function (e) {
            if (ractive.START < ractive.IMAGES.length - 1) {
                ractive.fire("set-images", ractive.START + 1);
            }
        },
        "set-images": function (start) {
            ractive.START = start;
            ractive.set("images", ractive.IMAGES.slice(start, start + 4));
        }
    });


    var filesRactive = new Ractive({
        el: "files",
        template: "#files_tpl",
        data: {
            files: []
        },
        lastSel: null
    });
    filesRactive.on({
        "load": function () {
            $.getJSON(MFilesUrl, function (ret) {
                var files = _.map(ret.data, function (v, i) {
                    v.url = MFileUrl + "?file=" + v.value + "&time=" + new Date().getTime();
                    return v;
                });
                filesRactive.set("files", files);
            });
        }
    });
    filesRactive.fire("load");

    var spcecRactive = new Ractive({
        el: "specs",
        template: "#specs_tpl",
        data: {
            specs: [],
            empty: function (specs) {
                return specs == null || specs.length == 0;
            }
        },
        lastSel: null
    });

    spcecRactive.on({
        "load": function () {
            $.getJSON(SpecsUrl, function (ret) {
                spcecRactive.set("specs", ret.data);
            });
        }
    });
    spcecRactive.fire("load");

    var pricesRactive = new Ractive({
        el: "prices",
        template: "#prices_tpl",
        data: {
            prices: [],
            format: function (s, e) {
                return e == 0 ? s + " 以上" : s + " - " + e;
            }
        },
        lastSel: null
    });

    spcecRactive.on({
        "load": function () {
            $.getJSON(PricesUrl, function (ret) {
                var ps = ret.data;
                var min = _.min(ps, function (p) {
                    return p.start_quantity;
                });
                _.map(ps, function (v, i) {
                    if (v.id == min.id) {
                        v["type"] = "单价";
                    } else {
                        v["type"] = "优惠价";
                    }
                });
                pricesRactive.set("prices", ps);
            });
        }
    });
    spcecRactive.fire("load");

    ractive.fire("load");
});

$(function () {
    $('#btn-collect').click(function () {
        alert($(this).data("product_id"));
    });
});

$(function () {
    $("a[tabtext]").click(function () {
        var $this = $(this);
        var $p = $this.parent();

        $(".tab li").removeClass("active");

        $("div.tab-content div.box").hide();

        $(".box-tab-" + $this.attr("tabtext")).show();

        $p.addClass("active");
    });
});

$(function () {
    YSCookie.add_viewed_product(ProductId, ProductName);
});

function _star(){
    $('div.raty').raty({
        path:"/public/libs/raty/images",
        score: function() {
            return $(this).data('score');
        },
        readOnly: function() {
           return true;
        }
    });
}


$(function(){
    var url = page.reloadURL();//+"&limit=1";
    $('div.BlockDiscuss').load(url, _star);
});

function onPage(p){
    page.set("page", p);
    var url = page.reloadURL();//+"&limit=1";;
    $('div.BlockDiscuss').load(url, _star);
}

