$(function () {
    $("#provider").select2({
        placeholder: "选择制造商/供应商",
       // minimumInputLength: 1,
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
            return  a.name;
        },// omitted for brevity, see the source of this page
        formatSelection: function (a) {
            if (a.id != 0) {
                return  a.name;
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
       // minimumInputLength: 1,
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
            var level = S(a.code).count("-");
            return S('&nbsp;').times(level * 4).s + a.code + " " +  a.name;
        },// omitted for brevity, see the source of this page
        formatSelection: function (a) {
            if (a.id != 0) {
                return a.name;
            }
            return "";
        },  // omitted for brevity, see the source of this page
        dropdownCssClass: "bigdrop", // apply css that makes the dropdown taller
        escapeMarkup: function (m) {
            return m;
        } // we do not want to escape markup since we are displaying html in results
    });

    $(".btn-tag").click(function (e) {
        var $this = $(this);
        var tag = "#" + $this.text();
        var tags = $("input[name='p.Tags']").val();
        if (tags.indexOf(tag) == -1) {
            $("input[name='p.Tags']").val(tags + " " + tag);
        }
    });
    //////////////////////////////////////////////////////////////////////////////////////

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
                $.getJSON(ProductSdImagesUrl, function (ret) {
                    var images = _.map(ret.data, function (v, i) {
                        v.url = ImageSdUrl + "?file=" + v.value + "&time=" + new Date().getTime();
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
                $it.attr("style", "width:110px;height:110px;");
                ractive.lastSel = $it;
                ractive.set("sel", $it.data("id"));
            },
            "deselected": function (e) {
                ractive.lastSel && ractive.lastSel.attr("style", "width:100px;height:100px;");
                ractive.fire("clear");
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeleteSdImageUrl + "?id=" + id, function () {
                    var images = ractive.get("images");
                    images.splice(index, 1);
                    ractive.fire("clear");
                });
            }
        });
        ractive.fire("load");
    })();


    var imagesPicsRactive = new Ractive({
        el: "images_pics",
        template: "#images_pics_tpl",
        data: {
            images: []
        },
        lastSel: null
    });

    (function () {
        var ractive = imagesPicsRactive;
        ractive.on({
            "load": function () {
                $.getJSON(ProductImagesPicsUrl, function (ret) {
                    var images = _.map(ret.data, function (v, i) {
                        v.url = ImagePicUrl + "?file=" + v.value + "&time=" + new Date().getTime();
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
                $it.attr("style", "width:60%;height:60%;");
                ractive.lastSel = $it;
                ractive.set("sel", $it.data("id"));
            },
            "deselected": function (e) {
                ractive.lastSel && ractive.lastSel.attr("style", "width:50%;height:50%;");
                ractive.fire("clear");
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", 0);
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeleteImagePicUrl + "?id=" + id, function () {
                    var images = ractive.get("images");
                    images.splice(index, 1);
                    ractive.fire("clear");
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
            imagesPicsRactive.fire("load");

            $('.MultiFile-remove').click();
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
    (function () {
        var ractive = filesRactive;
        ractive.on({
            "load": function () {
                $.getJSON(MFilesUrl, function (ret) {
                    var files = _.map(ret.data, function (v, i) {
                        v.url = MFileUrl + "?file=" + v.value + "&time=" + new Date().getTime();
                        return v;
                    });
                    ractive.set("files", files);
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
            },
            "deselected": function (e) {
                ractive.fire("clear");
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeleteMFileUrl + "?id=" + id, function () {
                    var files = ractive.get("files");
                    files.splice(index, 1);
                    ractive.fire("clear");
                });
            }
        });
        ractive.fire("load");

        $('#mFileForm').ajaxForm({
            dataType: 'json',
            success: function (ret) {
                alert(ret.message);
                ractive.fire("load");
                $('.MultiFile-remove').click();
            }
        });
    })();

    var spcecRactive = new Ractive({
        el: "specs",
        template: "#specs_tpl",
        data: {
            specs: []
        },
        lastSel: null
    });

    (function () {
        var ractive = spcecRactive;
        ractive.on({
            "load": function () {
                $.getJSON(SpecsUrl, function (ret) {
                    ractive.set("specs", ret.data);
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

                var specs = ractive.get("specs");
                var spec = _.find(specs, function (spec) {
                    return spec.id == $it.data("id");
                });

                ractive.set("sel", $it.data("id"));
                $("#div_id").show();

                $("#specForm input[name=name]").val(spec.name);
                $("#specForm input[name=id]").val(spec.id);
                $("#specForm input[name=value]").val(spec.value);
            },
            "deselected": function (e) {
                ractive.fire("clear");
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeleteSpecUrl + "?id=" + id, function () {
                    var specs = ractive.get("specs");
                    specs.splice(index, 1);
                    ractive.fire("clear");
                });
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
                $("#specForm input[name=name]").val("");
                $("#specForm input[name=id]").val("");
                $("#specForm input[name=value]").val("");
                $("#div_id").hide();
            }
        });
        ractive.fire("load");

        $('#specForm').ajaxForm({
            dataType: 'json',
            success: function (ret) {
                alert(ret.message);

                ractive.fire("load");
                ractive.fire("clear");
            }
        });
    })();

    ////////////////////////////////////////////////
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

    (function () {
        var ractive = pricesRactive;
        ractive.on({
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
                    ractive.set("prices", ps);
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

                var prices = ractive.get("prices");
                var price = _.find(prices, function (price) {
                    return price.id == $it.data("id");
                });

                ractive.set("sel", $it.data("id"));
                $("#price_div_id").show();

                $("#priceForm").show();
                $("#priceForm input[name=id]").val(price.id);
                $("#priceForm input[name=price]").val(price.price);
                $("#priceForm input[name=start_quantity]").val(price.start_quantity);
                $("#priceForm input[name=end_quantity]").val(price.end_quantity);
            },
            "deselected": function (e) {
                ractive.fire("clear");
            },
            "delete": function (e, index) {
                var $it = $(e.node);
                var id = $it.data("id");
                doAjaxPost(DeletePriceUrl + "?id=" + id, function () {
                    var prices = ractive.get("prices");
                    prices.splice(index, 1);
                    ractive.fire("clear");
                });
            },
            "clear": function () {
                ractive.lastSel = null;
                ractive.set("sel", null);
                $("#priceForm").hide();
                $("#priceForm input[type=text]").val("");
                $("#price_div_id").hide();
            }
        });
        ractive.fire("load");

        $('#priceSplitForm').ajaxForm({
            dataType: 'json',
            beforeSubmit: function () {
                return (confirm("确认执行定价操作？对已有定价会进行清除操作！"));
            },
            success: function (ret) {
                alert(ret.message);
                if (ret.ok) {
                    ractive.fire("load");
                    ractive.fire("clear");
                } else {
                    $('input[name=' + ret.data + ']').focus();
                }
            }
        })

        $('#priceForm').ajaxForm({
            dataType: 'json',
            success: function (ret) {
                alert(ret.message);
                if (ret.ok) {
                    ractive.fire("load");
                    ractive.fire("clear");
                } else {
                    $('input[name=' + ret.data + ']').focus();
                }
            }
        });
    })();

    //////////////////////////////////////////////


    var logsRactive = new Ractive({
        el: "logs",
        template: "#logs_tpl",
        data: {
            logs: [],
            stockNumber: stockNumber,
            format: function (d) {
                return yunshang.mRenderTime(d);
            }
        },
        lastSel: null
    });

    (function () {
        var ractive = logsRactive;
        ractive.on({
            "load": function () {
                $.getJSON(StockLogsUrl, function (ret) {
                    ractive.set("logs", ret.data);
                });
            }
        });
        ractive.fire("load");

        $('#stockLogForm').ajaxForm({
            dataType: 'json',
            beforeSubmit: function () {
                return confirm("确认进行入库操作？");
            },
            success: function (ret) {
                alert(ret.message);
                ractive.set("stockNumber", ret.data)

                ractive.fire("load");
            }
        });
    })();
});

//
$(function () {
    $('#updateLogoForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            alert(ret.message);
            if (ret.ok) {
                $img = $('img.plogo');
                $img.attr("src", $img.attr("src") + "&time=" + new Date().getTime());
            }
            $('.MultiFile-remove').click();
        }
    });
    $(".fancybox").fancybox();
});

$(function () {
    window.UEDITOR_CONFIG.imageManagerUrl = ImagePicsListUrl;
    window.UEDITOR_CONFIG.imageManagerPath = ImagePicUeditorUrl;
    window.UEDITOR_CONFIG.imageUrl = ImagePicUeditorUploadUrl;
    window.UEDITOR_CONFIG.imagePath = ImagePicUeditorUrl;
    window.UEDITOR_CONFIG.imageFieldName = "file";
    window.UEDITOR_CONFIG.savePath = [ '默认' ];

    var editor = UE.getEditor('editor-container');
    if (pId) {
        setTimeout(function () {
            editor.setContent(detail);
        }, 1000);
    }
    $("#btnSaveDetail").click(function () {
        doAjaxPost(saveDetailUrl, {content: editor.getContent()}, function (ret) {

        });
    });
});