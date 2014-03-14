var TheTable = function () {
    return {
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(yunshang.extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": dataUrl,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "name", "bSortable": true },
                    { "mData": "code", "bSortable": true },
                    { "mData": "created_at", "bSortable": true, "mRender": yunshang.mRenderTime },
                    { "mData": "enabled", "bSortable": true,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">可用</span>' : '<span class="label label-warn">不可用</span>';
                        }
                    },
                    { "mData": "description", "bSortable": false },
                ],
                "aaSorting": [
                    [0, 'asc']
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
                newUrl: newUrl,
                changeStatusUrl: changeStatusUrl,
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });

            ractive.on({
                    "delete": function () {
                       // doAjaxPost(deleteUrl +"?id=" +  ractive.getSelectedData()[0].id, function(){
                       //    ractive.refreshTable();
                       // });
                    }
                }
            );

            $("#e1").select2({
                placeholder: "选择分类状态"
            });
        }
    };
}();

$(function () {
    TheTable.init();
})