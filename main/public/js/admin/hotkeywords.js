$(function () {
    var keywordsRactive = new Ractive({
        el: "keywords",
        template: "#keywords_tpl",
        data: {
            keywords: []
        },
        lastSel: null
    });
    (function () {
        var ractive = keywordsRactive;
        ractive.on({
            "load": function () {
                $.getJSON(dataURL, function (ret) {
                    ractive.set("keywords", ret.data);
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

                ractive.fire("deselected");

                ractive.lastSel = $it;

                var keywords = ractive.get("keywords");
                var keyword = _.find(keywords, function (keyword) {
                    return keyword.id == $it.data("id");
                });

                ractive.set("sel", $it.data("id"));
                $("#div_id").show();


                $("#keywordForm input[name=id]").val(keyword.id);
                $("#keywordForm input[name=value]").val(keyword.value);
            },
            "deselected": function (e) {
                ractive.fire("clear");
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(deleteURL + "?id=" + id, function () {
                    var keywords = ractive.get("keywords");
                    keywords.splice(index, 1);
                    ractive.fire("clear");
                });
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
                $("#keywordForm input[name=id]").val("");
                $("#keywordForm input[name=value]").val("");
                $("#div_id").hide();
            },
            "first": function (e, id) {
                doAjaxPost(setFirstURL + "?id=" + id, function (ret) {
                    ractive.fire("clear");
                    ractive.fire("load");
                });
            }
        });
        ractive.fire("load");
    })();
    $('#keywordForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            alert(ret.message);

            keywordsRactive.fire("load");
        }
    });
});