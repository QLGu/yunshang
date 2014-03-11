/**
 * require Ractive.js, jQuery.js
 *
 */
Ractive.delimiters = [ '[[', ']]' ];
Ractive.tripleDelimiters = [ '[[[', ']]]' ];

function __check(o) {
    if (!o) {
        alert("obj is null");
    }
    return o;
}

var DatatableToolBar = Ractive.extend({
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

function mRenderTime(data) {
    if (data == "0001-01-01T00:00:00Z") {
        return "-"
    }
    return data;
}
