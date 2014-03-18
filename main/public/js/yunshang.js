/**
 * require Ractive.js, jQuery.js, moment.js
 *
 */
Ractive.delimiters = [ '[[', ']]' ];
Ractive.tripleDelimiters = [ '[[[', ']]]' ];

var yunshang = (function () {
    function __check(o) {
        if (!o) {
            alert("obj is null");
        }
        return o;
    }

    function GetDatatableToolBar() {
        return Ractive.extend({
            el: "table_tools",
            template: "#table_tools_template", // this will be applied to all DatatableToolBar instances
            selStatus: "",

            // initialisation code
            init: function (options) {
                var self = this;
                this.reset();

                this.table = options.table;
                if (!this.table) {
                    alert('Please init options.table');
                }
                this.newUrl = options.newUrl;
                this.changeStatusUrl = options.changeStatusUrl;

                this.on({
                    "selected": function () {
                        self.set("selected", true)
                        self.set("enabled", self.getSelectedData()[0].enabled);
                    },
                    "deselected": function () {
                        self.reset();
                    },
                    "refresh": function () {
                        this.refreshTable();
                    },
                    "filter-enabled": function (event) {
                        self.selStatus = $(event.node).val();
                        self.refreshTable();
                    },
                    "new": function () {
                        self.openWindow();
                    },
                    "edit": function () {
                        self.openWindow(self.getSelectedData()[0].id);
                    },
                    "change-status": function () {
                        var url = changeStatusUrl + "?id=" + self.getSelectedData()[0].id;
                        doAjaxPost(url, function () {
                            self.refreshTable();
                        });
                    }
                });
            },
            reset: function () {
                this.set("selected", false);
                this.set("enabled", "default");
            },
            openWindow: function (id) {
                var self = this;
                $.fancybox.open({
                    href: self.newUrl + "?id=" + (id || ""),
                    type: 'iframe',
                    width: '70%',
                    minHeight: 600,
                    padding: 5,
                    afterClose: function (e) {
                        self.refreshTable();
                    }
                });
            },
            refreshTable: function () {
                this.tableInstance().fnDraw(true);
                this.reset();
            },
            getSelectedData: function () {
                var oTT = TableTools.fnGetInstance(this.table.id);
                return oTT.fnGetSelectedData();
            },
            tableInstance: function () {
                return __check(this.table.instance);
            }
        });
    }

    function mRenderTime(data) {
        if (data == "0001-01-01T00:00:00Z") {
            return "-"
        }
        var m = moment(data.substring(0, "2014-03-12T10:39:16".length));
        return m.format("YYYY-MM-DD HH:mm:ss");
    }

    function extendDefaultOptions(options) {
        var __defaultDataTableOptions = {
            "bProcessing": true,
            "bServerSide": true,
            "sDom": "T<'row-fluid'<'span3'l><'span3'r><'span6'f>>t<'row-fluid'<'span6'i><'span6'p>>",
            "oTableTools": {
                "sSwfPath": "/public/media/swf/copy_csv_xls_pdf.swf",
                "sRowSelect": "single",
                "fnRowSelected": function (nodes) {
                    options.ractive().fire("selected");
                },
                "fnRowDeselected": function (_nodes) {
                    options.ractive().fire("deselected");
                },
                "aButtons": [
                    "copy",
                    "print", {
                        "sExtends": "collection",
                        "sButtonText": "Save",
                        "aButtons": [ "csv", "xls", "pdf" ] } ]
            },
            "aaSorting": [
                [0, 'desc']
            ],
            "aLengthMenu": [
                [10, 20, 30, 50, -1],
                ["10条", "20条", "30条", "50条", "全部"] // change per page values here
            ],
            // set the initial value
            "iDisplayLength": 10,
            "sPaginationType": "bootstrap",
            "oLanguage": {
                "sLengthMenu": "每页 _MENU_ 记录",
                "sInfo": "共计 _TOTAL_ 条， 显示_START_ 到 _END_ 条",
                "sInfoEmpty": "",
                "sEmptyTable": "查询不到数据",
                "sSearch": "搜 索:",
                "oPaginate": {
                    "sFirst": "首页",
                    "sPrevious": "前一页",
                    "sNext": "后一页",
                    "sLast": "末页"
                },
                "sInfoEmtpy": "没有数据",
                "sProcessing": "正在加载数据...",

            }
        }
        return $.extend(options, __defaultDataTableOptions);
    }

    return {
        GetDatatableToolBar: GetDatatableToolBar,
        mRenderTime: mRenderTime,
        extendDefaultOptions: extendDefaultOptions
    };
})();

function _isDevAjax(ajaxOptions) {
    return ajaxOptions["url"].indexOf("/@reveltang_dev/__source_changed") != -1;
}

function reloadWindow() {
    window.location.reload();
}

function doAjaxPost(url, data, after) {
    if (_.isFunction(data)) {
        after = data;
        data = {};
    }
    if (confirm("确认执行此操作？")) {
        $.post(url, data, function (result) {
            if (result.ok) {
                alert("操作成功：" + result.message);
                if (after) {
                    after(result)
                }
            } else {
                alert("操作失败：" + result.message);
            }
        }, "json");
    }
}

function isUserLogined() {
    //UserLogined from page
    return UserLogined === 'true';
}

$(function () {
    // Ajax Global Config
    $(document).ajaxSend(function (_event, _jqXHR, ajaxOptions) {
        if (!_isDevAjax(ajaxOptions)) { // not dev ajax
            if ($("#loadingbar").length === 0) {
                $("body").append("<div id='loadingbar'></div>")
                $("#loadingbar").addClass("waiting").append($("<dt/><dd/>"));

                $("#loadingbar").width((50 + Math.random() * 30) + "%");
            }
        }
    });
    $.ajaxPrefilter(function (opt, origOpt, jqxhr) {
        jqxhr.always(function () {
            $("#loadingbar").width("101%").delay(200).fadeOut(400, function () {
                $(this).remove();
            });
            //$("[data-plugin]").plugin();
        });
    });

    // Ajaxable Links
    $("a[data-url], button[data-url]").click(function () {
        var $this = $(this);
        if ($this.data("logined")) {
            if (!isUserLogined()) {
                alert("您还未登录，请您登录后执行此操作！");
                return;
            }
        }
        $this.attr("disabled", true);
        if (confirm($this.data("confirm-text") || "确认执行此操作？")) {
            $.post($(this).data("url"),function (result) {
                if (result.ok) {
                    alert("操作完成：" + result.message);
                    var after = $this.data("after");
                    if (after) {
                        eval(after);
                    }
                } else {
                    alert("操作未完成：" + result.message);
                }
            }, "json").always(function () {
                $this.attr("disabled", false);
            });
        }
    });
});
