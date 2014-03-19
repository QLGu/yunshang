/**
 * require(jQuery).
 */
$(function () {
    var imagesRactive = new Ractive({
        el: "images",
        template: "#images_tpl",
        data: {
            images: [],
            curr: null
        }
    });
    (function () {
        var ractive = imagesRactive;
        ractive.on({
            "load": function () {
                $.getJSON(dataURL, function (ret) {
                    var images = _.map(ret.data, function (v, i) {
                        v.url = ImageUrl + "?file=" + v.value + "&time=" + new Date().getTime();
                        return v;
                    });
                    ractive.images = images;
                    ractive.set("images", ractive.images);

                    ractive.fire("selected", "", 0);
                });
            },
            "selected": function (e, index) {
                ractive.set("curr", ractive.images[index]);
                ractive.index = index;
            }
        });
        ractive.fire("load");

        setInterval(function () {
            var len = ractive.images.length;
            var index = (ractive.index + 1) % len;
            ractive.fire("selected", "", index);
        }, 1000 * 8);
    })();
});

