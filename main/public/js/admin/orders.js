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
                        var s = osJSON[data];
                        return s >= 7 ? "<font color='red'>" + s + "</font>" : s;
                    }},
                    { "mData": "amount", "bSortable": true, "mRender": function (data) {
                        return "￥" + data;
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
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
            var UserDatatableToolBar = yunshang.GetDatatableToolBar().extend({
                selCertified: "",
                reset: function () {
                    this._super();
                    this.set("locked", "default");
                    this.set("can_confirm_pay", false);
                    this.set("can_confirm_verify", false);
                    this.set("can_confirm_shiped", false);
                    this.set("can_confirm_lock", false);
                    this.set("can_back", false);
                },
                init: function (options) {
                    this._super(options);
                }
            });
            ractive = new UserDatatableToolBar({
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });
            ractive.on({
                    "selected": function () {
                        var status = ractive.getSelectedData()[0].status;
                        ractive.set("locked", status == 8);

                        if (status != 6) {
                            ractive.set("can_confirm_lock", true);
                        }

                        var paymentId = ractive.getSelectedData()[0].payment_id;
                        //待支付 && 银行转账 => 确认支付
                        if (status == 2 && paymentId == 3) {
                            ractive.set("can_confirm_pay", true);
                        }

                        //已支付           => 确认待发货
                        if (status == 3) {
                            ractive.set("can_confirm_verify", true);
                        }

                        //已确认        => 确认已发货
                        if (status == 4) {
                            ractive.set("can_confirm_shiped", true);
                        }

                        if (status == 5 || status == 6) {
                            ractive.set("can_back", true);
                        }
                    },
                    "confirm-pay": function () {
                        var msg = "确认已经收到用户此订单的转账汇款， 可以发货？";
                        if (confirm(msg) && confirm("再次确认" + msg)) {
                            var url = changePayedUrl + "?id=" + ractive.getSelectedData()[0].id;
                            doAjaxPost(url, function () {
                                ractive.refreshTable();
                            });
                        }
                    },
                    "confirm-verify": function () {
                        var msg = "确认此订单可以发货？";
                        if (confirm(msg)) {
                            var url = changeVerifyUrl + "?id=" + ractive.getSelectedData()[0].id;
                            doAjaxPost(url, function () {
                                ractive.refreshTable();
                            });
                        }
                    },
                    "confirm-ship": function () {
                        var msg = "确认此订单已发货？";
                        if (confirm(msg)) {
                            var url = changeShipedUrl + "?id=" + ractive.getSelectedData()[0].id;
                            doAjaxPost(url, function () {
                                ractive.refreshTable();
                            });
                        }
                    },
                    "change-lock": function () {
                        var url = changeLockUrl + "?id=" + ractive.getSelectedData()[0].id;
                        doAjaxPost(url, function () {
                            ractive.refreshTable();
                        });
                    },
                    "set-back": function () {
                        var url = setBackUrl + "?id=" + ractive.getSelectedData()[0].id;
                        doAjaxPost(url, function () {
                            alert("请别忘了到此订单相关产品处维护库存信息！")
                            ractive.refreshTable();
                        });
                    },
                    "view-order": function () {
                        var code = ractive.getSelectedData()[0].code;
                        var userId = ractive.getSelectedData()[0].user_id;
                        window.open( showOrderUrl + "?code=" + code + "&userId=" + userId);
                        /*$.fancybox.open({
                            href: showOrderUrl + "?code=" + code + "&userId=" + userId,
                            type: 'iframe',
                            padding: 5
                        });*/
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

            $("#e1").select2({  placeholder: "选择订单状态"  });
        }
    };
}();

$(function () {
    TableUsers.init();
})