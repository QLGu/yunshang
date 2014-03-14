$(function () {
    $("#provider").select2({
        placeholder: "选择制造商/供应商",
        minimumInputLength: 1,
        ajax: { // instead of writing the function to execute the request we use Select2's convenient helper
            url: providersDataUrl,
            dataType: 'json',
            data: function (term, page) {
                return {
                    q: term, // search term
                    page_limit: 10
                };
            },
            results: function (ret, page) { // parse the results into the format expected by Select2.
                // since we are using custom formatting functions we do not need to alter remote JSON data
                return {results: ret.data};
            }
        },
        initSelection: function (element, callback) {
            // the input tag has a value attribute preloaded that points to a preselected movie's id
            // this function resolves that id attribute to an object that select2 can render
            // using its formatResult renderer - that way the movie name is shown preselected
            var id = $(element).val();
            if (id !== "") {
                $.ajax(providerDataUrl + "?id=" + id, {
                    data: {
                        apikey: "ju6z9mjyajq2djue3gbvv26t"
                    },
                    dataType: "json"
                }).done(function (ret) {
                    callback(ret.data);
                });
            }
        },
        formatResult: function (a) {
            return a.id + " " + a.name;
        },// omitted for brevity, see the source of this page
        formatSelection: function (a) {
            if (a.id != 0) {
                return a.id + " " + a.name;
            }
            return "";
        },  // omitted for brevity, see the source of this page
        dropdownCssClass: "bigdrop", // apply css that makes the dropdown taller
        escapeMarkup: function (m) {
            return m;
        } // we do not want to escape markup since we are displaying html in results
    });

    $("#category").select2({
        placeholder: "选择商品分类",
        minimumInputLength: 1,
        ajax: { // instead of writing the function to execute the request we use Select2's convenient helper
            url: categoriesDataUrl,
            dataType: 'json',
            data: function (term, page) {
                return {
                    q: term, // search term
                    page_limit: 10
                };
            },
            results: function (ret, page) { // parse the results into the format expected by Select2.
                // since we are using custom formatting functions we do not need to alter remote JSON data
                return {results: ret.data};
            }
        },
        initSelection: function (element, callback) {
            // the input tag has a value attribute preloaded that points to a preselected movie's id
            // this function resolves that id attribute to an object that select2 can render
            // using its formatResult renderer - that way the movie name is shown preselected
            var id = $(element).val();
            if (id !== "") {
                $.ajax(categoryDataUrl + "?id=" + id, {
                    data: {
                        apikey: "ju6z9mjyajq2djue3gbvv26t"
                    },
                    dataType: "json"
                }).done(function (ret) {
                    callback(ret.data);
                });
            }
        },
        formatResult: function (a) {
            return a.id + " " + a.name;
        },// omitted for brevity, see the source of this page
        formatSelection: function (a) {
            if (a.id != 0) {
                return a.id + " " + a.name;
            }
            return "";
        },  // omitted for brevity, see the source of this page
        dropdownCssClass: "bigdrop", // apply css that makes the dropdown taller
        escapeMarkup: function (m) {
            return m;
        } // we do not want to escape markup since we are displaying html in results
    });

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
        },
        "delete": function (e) {
            var $it = $(e.node);
            var id = $it.data("id");
            alert(id)
            doAjaxPost(DeleteSdImageUrl + "?id=" + id, function () {
                var images = ractive.get("images");
                ractive.set("images", _.filter(images, function (image) {
                    //alert(image.id)
                    return  image.id != id;
                }));
            });
        }
    });
    ractive.fire("load");

    $('#sdImageForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            alert(ret.message);
            ractive.fire("load");
            $('.MultiFile-remove').click();
        }
    });
});