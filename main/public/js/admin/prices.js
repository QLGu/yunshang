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
                    { "mData": "model", "bSortable": true },
                    { "mData": "quantity", "bSortable": true },
                    { "mData": "created_at", "bSortable": false, "mRender": yunshang.mRenderTime},
                    { "mData": "replies", "bSortable": true }
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
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });
            ractive.on({
                    "selected": function () {

                    },
                    "view-inquiry": function () {
                        $.fancybox.open({
                            href: newReplyURL + "?id=" + ractive.getSelectedData()[0].id,
                            type: 'iframe',
                            padding: 5,
                            afterClose: function (e) {
                                ractive.refreshTable();
                            }
                        });
                    },
                    "view-reply": function () {
                        $.fancybox.open({
                            href: newReplyURL + "?id=" + ractive.getSelectedData()[0].id,
                            type: 'iframe',
                            padding: 5,
                            afterClose: function (e) {
                                ractive.refreshTable();
                            }
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

            $("#e1").select2({  placeholder: "选择回复状态"  });
        }
    };
}();

$(function () {
    TableUsers.init();
})