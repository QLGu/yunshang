$(function () {
    var imagesRactive = new Ractive({
        el: "images",
        template: "#images_tpl",
        data: {
            images: []
        },
        lastSel: null
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
                    ractive.set("images", images);
                });
            },
            "click-row": function (e) {
                var $it = $(e.node);
                if ($it.data("id") == (ractive.lastSel && ractive.lastSel.data("id"))) {
                    ractive.fire("deselected", e);
                } else {
                    ractive.fire("selected", e);
                }
            },
            "selected": function (e) {
                var $it = $(e.node);


                ractive.lastSel = $it;
                ractive.set("sel", $it.data("id"));

                $("#imageLinkForm").show();
            },
            "deselected": function (e) {
                ractive.fire("clear");
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
                $("#imageLinkForm").hide();
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeleteImageUrl + "?id=" + id, function () {
                    var images = ractive.get("images");
                    images.splice(index, 1);
                    ractive.fire("clear");
                });
            },
            "first": function (e, id) {
                doAjaxPost(SetFirstAdImageUrl + "?id=" + id, function (ret) {
                    ractive.fire("clear");
                    ractive.fire("load");
                });
            },
            "set-link": function (e, id) {
                var $it = $(e.node);
                var link = $it.val();
                $.post(SetAdImageLinkURL + "?id=" + id, {link: link}, function (ret) {

                });
            }
        });
        ractive.fire("load");
    })();
    $('#imageForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            alert(ret.message);

            imagesRactive.fire("load");

            $('.MultiFile-remove').click();
        }
    });
});