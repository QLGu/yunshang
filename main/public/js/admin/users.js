var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": "/admin/users/data",
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "login_name", "bSortable": false,
                        "mRender": function (data) {
                            return (data == "admin") ? '<span class="label label-warning">admin</span>' : data;
                        }},
                    { "mData": "email", "bSortable": false},
                    { "mData": "real_name", "bSortable": false},
                    { "mData": "from", "bSortable": false},
                    { "mData": "scores", "bSortable": true},
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">可用</span>' : '<span class="label label-warn">不可用</span>';
                        }},
                    { "mData": "certified", "bSortable": true,
                        "mRender": function (data, type, full) {
                            return data ? '<span class="label label-success">已认证</span>' : '<span class="label label-warn">未认证</span>';
                        }},
                    { "mData": "last_sign_at", "bSortable": true, "mRender": mRenderTime},
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                    aoData.push({ name: "filter_certified", value: ractive.selCertified || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
            var UserDatatableToolBar = DatatableToolBar.extend({
                selCertified: "",
                reset: function () {
                    this._super();
                    this.set("certified", "default");
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
                    "selected": function () {
                        ractive.set("certified", ractive.getSelectedData()[0].certified);
                    },
                    "reset-password": function () {
                        var url = resetPasswordUrl + "?id=" + ractive.getSelectedData()[0].id;
                        doAjaxPost(url, function () {
                            ractive.refreshTable();
                        });
                    },
                    "change-certified": function () {
                        var url = changeCertifiedUrl + "?id=" + ractive.getSelectedData()[0].id;
                        doAjaxPost(url, function () {
                            ractive.refreshTable();
                        });
                    },
                    "view-userinfo": function () {
                        $.fancybox.open({
                            href: showUserInfosUrl + "?id=" + ractive.getSelectedData()[0].id,
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
                    },
                    "filter-certified": function (event) {
                        ractive.selCertified = $(event.node).val();
                        ractive.refreshTable();
                    }
                }
            );

            $("#e1").select2({  placeholder: "选择用户状态"  });
            $("#e2").select2({   placeholder: "选择认证状态"  });
        }
    };
}();

$(function () {
    TableUsers.init();
})