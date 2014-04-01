var TheTable = function () {
    return {
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(yunshang.extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": dataURL,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": false },
                    { "mData": "title", "bSortable": false },
                    { "mData": "created_at", "bSortable": false, "mRender": yunshang.mRenderTime },
                    { "mData": "publish_at", "bSortable": false, "mRender": yunshang.mRenderTime },
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">已发布</span>' : '<span class="label label-warn">未发布</span>';
                        }
                    },
                    { "mData": "tags", "bSortable": false },
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                    aoData.push({ name: "filter_tag", value: ractive.selTag || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown

            var DatatableToolBar = yunshang.GetDatatableToolBar();
            ractive = new DatatableToolBar({
                newUrl: newURL,
                changeStatusUrl: changeStatusURL,
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });

            ractive.on({
                    "preview": function () {
                        var url = previewURL.substring(0, previewURL.lastIndexOf("/"));
                        window.open(url + "/" + ractive.getSelectedData()[0].id, "新闻预览");
                    },
                    "filter-tag": function (event) {
                        ractive.selTag = $(event.node).val();
                        ractive.refreshTable();
                    },
                }
            );

            $("#e1").select2({
                placeholder: "选择产品状态"
            });
            //$("#e2").select2({
          //      placeholder: "选择产品标签"
           /// });
        }
    };
}();

$(function () {
    TheTable.init();
})