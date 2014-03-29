var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable(yunshang.extendDefaultOptions({
                "ractive": function () {
                    return ractive;
                },
                "sAjaxSource": ordersDataURL,
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": true},
                    { "mData": "status", "bSortable": false, "mRender": function (data) {
                        return osJSON[data];
                    }},
                    { "mData": "payment_id", "bSortable": false, "mRender": function (data) {
                        return pmJSON[data];
                    }},
                    { "mData": "shipping_id", "bSortable": false, "mRender": function (data) {
                        return spJSON[data];
                    }},
                    { "mData": "submit_at", "bSortable": false, "mRender": yunshang.mRenderTime},
                    { "mData": "pay_at", "bSortable": false, "mRender": yunshang.mRenderTime}
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                    aoData.push({ name: "filter_certified", value: ractive.selCertified || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
            var UserDatatableToolBar = yunshang.GetDatatableToolBar().extend({
                selCertified: "",
                reset: function () {
                    this._super();
                    this.set("certified", "default");
                    this.set("locked", "default");
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
                    "view-order": function () {
                        var code = ractive.getSelectedData()[0].code;
                        var userId = ractive.getSelectedData()[0].user_id;
                        $.fancybox.open({
                            href: showOrderUrl + "?code=" + code + "&userId=" + userId,
                            type: 'iframe',
                            padding: 5
                        });
                    },
                    "view-user": function () {
                        $.fancybox.open({
                            href: showUserInfosUrl + "?id=" + ractive.getSelectedData()[0].user_id,
                            type: 'iframe',
                            padding: 5
                        });
                    },
                    "filter-status": function (event) {
                        ractive.selStatus = $(event.node).val();
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