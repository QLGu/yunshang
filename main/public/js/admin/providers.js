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
                    { "mData": "short_name", "bSortable": true },
                    { "mData": "created_at", "bSortable": true, "mRender": yunshang.mRenderTime },
                    { "mData": "enabled", "bSortable": true,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">可用</span>' : '<span class="label label-warn">不可用</span>';
                        }
                    },
                    { "mData": "tags", "bSortable": true,
                        "mRender": function (data, type, full) {
                            return data == "推荐" ? '<span class="label label-success">推荐</span>' : '<span class="label label-warn">-</span>';
                        }
                    }
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                    aoData.push({ name: "filter_tags", value: ractive.selTags || ""});
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
                    "preview": function () {
                        var url = previewUrl.substring(0, previewUrl.lastIndexOf("/"));
                        window.open(url + "/" + ractive.getSelectedData()[0].id, "");
                    },
                    "delete": function () {
                        doAjaxPost(deleteUrl + "?id=" + ractive.getSelectedData()[0].id, function () {
                            ractive.refreshTable();
                        });
                    },
                    "filter-tags": function (event) {
                        ractive.selTags= $(event.node).val();
                        ractive.refreshTable();
                    }
                }
            );

            $("#e1").select2({
                placeholder: "选择制造商状态"
            });
            $("#e2").select2({
                placeholder: "选择是否推荐"
            });
        }
    };
}();

$(function () {
    TheTable.init();
})