var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(yunshang.extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": dataURL,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "target_id", "bSortable": false,
                        "mRender": function (data, t, full) {
                            if (full.target_type == 1) {
                                return '<a href="/products/p/' + data + '" target="_blank">' + full.target_name + "</a>";
                            } else if (full.target_type == 2) {
                                return '<a href="/articles/' + data + '" target="_blank">' + full.target_name + "</a>";
                            }
                        }},
                    { "mData": "scores", "bSortable": false},
                    { "mData": "content", "bSortable": false},
                    { "mData": "created_at", "bSortable": false, "mRender": yunshang.mRenderTime},
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">审核通过</span>' : '<span class="label label-warn">未审核</span>';
                        }}
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
            var UserDatatableToolBar = yunshang.GetDatatableToolBar().extend({
                selCertified: "",
                reset: function () {
                    this._super();
                },
                init: function (options) {
                    this._super(options);
                }
            });
            ractive = new UserDatatableToolBar({
                newUrl: "",
                changeStatusUrl: changeStatusUrl,
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });
            ractive.on({
                    "view-userinfo": function () {
                        $.fancybox.open({
                            href: showUserInfosUrl + "?id=" + ractive.getSelectedData()[0].user_id,
                            type: 'iframe',
                            padding: 5
                        });
                    },
                    "view-loginlog": function () {
                        $.fancybox.open({
                            href: showUserLoginLogsUrl + "?id=" + ractive.getSelectedData()[0].id,
                            type: 'iframe',
                            padding: 5
                        });
                    }
                }
            );

            $("#e1").select2({  placeholder: "选择评论状态"  });
        }
    };
}();

$(function () {
    TableUsers.init();
})