var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            var ractive = {};
            var sampleTable = $('#sample_1').dataTable({
                "bProcessing": true,
                "bServerSide": true,
                "sAjaxSource": "/admin/users/data",
                //"sDom":'T<"clear">lfrtip',
                "sDom": "T<'row-fluid'<'span3'l><'span3'r><'span6'f>>t<'row-fluid'<'span6'i><'span6'p>>",
                "oTableTools": {
                    "sSwfPath": "/public/media/swf/copy_csv_xls_pdf.swf",
                    //"sSelectedClass": "highlight",
                    "sRowSelect": "single",
                    "fnRowSelected": function (nodes) {
                        ractive.fire("selected");
                    },
                    "fnRowDeselected": function (_nodes) {
                        ractive.fire("deselected");
                    },
                    "aButtons": [
                        "copy",
                        "print", {
                            "sExtends": "collection",
                            "sButtonText": "Save",
                            "aButtons": [ "csv", "xls", "pdf" ] } ]
                },
                "aaSorting": [
                    [0, 'desc']
                ],
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
                "aLengthMenu": [
                    [10, 20, 30, 50, -1],
                    ["10条", "20条", "30条", "50条", "全部"] // change per page values here
                ],
                // set the initial value
                "iDisplayLength": 10,
                "sPaginationType": "bootstrap",
                "oLanguage": {
                    "sLengthMenu": "每页显示_MENU_ 记录",
                    "sInfo": "共计 _TOTAL_ 条， 显示_START_ 到 _END_ 条",
                    "sInfoEmpty": "",
                    "sEmptyTable": "查询不到数据",
                    "sSearch": "搜 索:",
                    "oPaginate": {
                        "sFirst": "首页",
                        "sPrevious": "前一页",
                        "sNext": "后一页",
                        "sLast": "末页"
                    },
                    "sInfoEmtpy": "没有数据",
                    "sProcessing": "正在加载数据...",

                },
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                    aoData.push({ name: "filter_certified", value: ractive.selCertified || ""});
                }
            });

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            //$('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown
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