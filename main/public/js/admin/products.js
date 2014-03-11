var TheTable = function () {
    return {
        //main function to initiate the module
        init: function () {

            var selStatus = "";
            var sampleTable;
            var ractive;

            function mRenderTime(data) {
                if (data == "0001-01-01T00:00:00Z") {
                    return "-"
                }
                return data;
            }

            function getSelectedData() {
                var oTT = TableTools.fnGetInstance('sample_1');
                return oTT.fnGetSelectedData();
            }

            function refreshTable() {
                sampleTable.fnDraw(true);
                ractive.reset();
            }

            function openWindow(id) {
                $.fancybox.open({
                    href: newProductUrl + "?id=" + (id || ""),
                    type: 'iframe',
                    width: '70%',
                    minHeight: 600,
                    padding: 5,
                    afterClose: function (e) {
                        refreshTable();
                    }
                });
            }

            Ractive.delimiters = [ '[[', ']]' ];
            Ractive.tripleDelimiters = [ '[[[', ']]]' ];
            ractive = new Ractive({
                el: "table_tools",
                template: "#table_tools_template",
                data: {
                    selected: false,
                    enabled: "default",
                    disabled: function () {
                        return  this.get("selected") === false ? "disabled" : "";
                    }
                }
            });
            ractive.reset = function () {
                ractive.set("selected", false);
                ractive.set("enabled", "default");
            };
            ractive.on({
                    "selected": function (rowdata) {
                        ractive.set("selected", true)
                        ractive.set("enabled", rowdata.enabled);
                    },
                    "deselected": function () {
                        ractive.reset();
                    },
                    "new": function () {
                        openWindow();
                    },
                    "edit": function () {
                        openWindow(getSelectedData()[0].id);
                    },
                    "fresh": function () {
                        sampleTable.fnDraw(true);
                        ractive.reset();
                    },
                    "change-status": function () {
                        var url = changeStatusUrl + "?id=" + getSelectedData()[0].id;
                        doAjaxPost(url, refreshTable);
                    },
                    "filter-enabled": function (event) {
                        selStatus = $(event.node).val();
                        refreshTable();
                    }
                }
            );

            sampleTable = $('#sample_1').dataTable({
                "bProcessing": true,
                "bServerSide": true,
                "sAjaxSource": productsDataUrl,
                "sDom": "T<'row-fluid'<'span3'l><'span3'r><'span6'f>>t<'row-fluid'<'span6'i><'span6'p>>",
                "oTableTools": {
                    "sSwfPath": "/public/media/swf/copy_csv_xls_pdf.swf",
                    "sRowSelect": "single",
                    "fnRowSelected": function (nodes) {
                        ractive.fire("selected", getSelectedData()[0]);
                    },
                    "fnRowDeselected": function (_nodes) {
                        ractive.fire("deselected");
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
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": false},
                    { "mData": "name", "bSortable": false},
                    { "mData": "model", "bSortable": false},
                    { "mData": "created_at", "bSortable": false},
                    { "mData": "enabled_at", "bSortable": false, "mRender": mRenderTime },
                    { "mData": "unenabled_at", "bSortable": false,
                        "mRender": mRenderTime},
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">已上架</span>' : '<span class="label label-warn">未上架</span>';
                        }
                    }
                ],
                "aLengthMenu": [
                    [10, 20, 30, 50, -1],
                    ["10条", "20条", "30条", "50条", "全部"] // change per page values here
                ],
                // set the initial value
                "iDisplayLength": 10,
                "sPaginationType": "bootstrap",
                "oLanguage": {
                    "sLengthMenu": "每页显示_MENU_ 记录",
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

                },
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: selStatus});
                }
            });

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            //$('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown

            $("#e1").select2({
                placeholder: "选择产品状态"
            });
        }
    };
}();

$(function () {
    TheTable.init();
})