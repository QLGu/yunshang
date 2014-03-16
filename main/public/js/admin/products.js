var TheTable = function () {
    return {
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(yunshang.extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": productsDataUrl,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": false },
                    { "mData": "name", "bSortable": false },
                    { "mData": "model", "bSortable": false },
                    { "mData": "stock_number", "bSortable": true, "mRender": function (data) {
                        return data < 10 ? '<font color="red">' + data + '</font>' : data;
                    } },
                    { "mData": "created_at", "bSortable": false, "mRender": yunshang.mRenderTime },
                    { "mData": "enabled_at", "bSortable": false, "mRender": yunshang.mRenderTime },
                    { "mData": "unenabled_at", "bSortable": false, "mRender": yunshang.mRenderTime },
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">已上架</span>' : '<span class="label label-warn">未上架</span>';
                        }
                    }
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown

            var DatatableToolBar = yunshang.GetDatatableToolBar();
            ractive = new DatatableToolBar({
                newUrl: newProductUrl,
                changeStatusUrl: changeStatusUrl,
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });

            ractive.on({
                    "preview": function () {
                        var url = previewProductUrl.substring(0, previewProductUrl.lastIndexOf("/"));
                        window.open(url + "/" + ractive.getSelectedData()[0].id, "产品预览");
                    }
                }
            );

            $("#e1").select2({
                placeholder: "选择产品状态"
            });
        }
    };
}();

$(function () {
    TheTable.init();
})