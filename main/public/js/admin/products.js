var TheTable = function () {
    return {
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": productsDataUrl,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": false },
                    { "mData": "name", "bSortable": false },
                    { "mData": "model", "bSortable": false },
                    { "mData": "created_at", "bSortable": false, "mRender": mRenderTime },
                    { "mData": "enabled_at", "bSortable": false, "mRender": mRenderTime },
                    { "mData": "unenabled_at", "bSortable": false, "mRender": mRenderTime },
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
                        alert('TODO');
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