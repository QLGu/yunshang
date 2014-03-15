$(function () {
    var ractive = new Ractive({
        el: "images",
        template: "#images_tpl",
        data: {
            images: []
        },
        lastSel: null
    });
    ractive.on({
        "load": function () {
            $.getJSON(ProductSdImagesUrl, function (ret) {
                var images = _.map(ret.data, function (v, i) {
                    v.url = ImageSdUrl + "?file=" + v.value + "&time=" + new Date().getTime();
                    return v;
                });
                ractive.set("images", images);
            });
        },
        "selected": function (e) {
            ractive.lastSel && ractive.lastSel.attr("style", "width:100px;height:100px;");

            var $it = $(e.node);
            $it.attr("style", "width:110px;height:110px;");
            ractive.lastSel = $it;
            ractive.set("sel", $it.data("id"));
        }
    });
    ractive.fire("load");


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
            specs: []
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
});