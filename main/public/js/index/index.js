/**
 * require(jQuery).
 */
$(function () {
    var imagesRactive = new Ractive({
        el: "images",
        template: "#images_tpl",
        data: {
            images: []
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
                    ractive.set("curr", ractive.images[0]);
                });
            },
            "selected": function (e, index) {
                ractive.set("curr", ractive.images[index]);
            }
        });
        ractive.fire("load");
    })();
});

