var TableUsers = function () {
    return {
        //main function to initiate the module
        init: function () {
            var selStatus = "";
            var selCertified = "";
            var sampleTable;
            var ractive;

            function getSelectedData() {
                var oTT = TableTools.fnGetInstance('sample_1');
                return oTT.fnGetSelectedData();
            }

            function refreshTable() {
                sampleTable.fnDraw(true);
                ractive.reset();
            }

            Ractive.delimiters = [ '[[', ']]' ];
            Ractive.tripleDelimiters = [ '[[[', ']]]' ];
            ractive = new Ractive({
                el: "table_tools",
                template: "#table_tools_template",
                data: {
                    selected: false,
                    enabled: "default",
                    certified: "default",
                    disabled: function () {
                        return  this.get("selected") === false ? "disabled" : "";
                    }
                }
            });
            ractive.reset = function () {
                ractive.set("selected", false);
                ractive.set("enabled", "default");
                ractive.set("certified", "default");
            };

            ractive.on({
                    "selected": function (rowdata) {
                        ractive.set("selected", true)
                        ractive.set("enabled", rowdata.enabled);
                        ractive.set("certified", rowdata.certified);
                    },
                    "deselected": function () {
                        ractive.reset();
                    },
                    "refresh": function () {
                        refreshTable();
                    },
                    "change-status": function () {
                        var url = changeStatusUrl + "?id=" + getSelectedData()[0].id;
                        doAjaxPost(url, refreshTable);
                    },
                    "reset-password": function () {
                        var url = resetPasswordUrl + "?id=" + getSelectedData()[0].id;
                        doAjaxPost(url, refreshTable);
                    },
                    "change-certified": function () {
                        var url = changeCertifiedUrl + "?id=" + getSelectedData()[0].id;
                        doAjaxPost(url, refreshTable);
                    },
                    "view-userinfo": function () {
                        $.fancybox.open({
                            href: showUserInfosUrl + "?id=" + getSelectedData()[0].id,
                            type: 'iframe',
                            padding: 5
                        });
                    },
                    "view-loginlog": function () {
                        $.fancybox.open({
                            href: showUserLoginLogsUrl + "?id=" + getSelectedData()[0].id,
                            type: 'iframe',
                            padding: 5
                        });
                    },
                    "filter-enabled": function (event) {
                        selStatus = $(event.node).val();
                        refreshTable();
                    },
                    "filter-certified": function (event) {
                        selCertified = $(event.node).val();
                        refreshTable();
                    }
                }
            );

            sampleTable = $('#sample_1').dataTable({
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
                        ractive.fire("selected", getSelectedData()[0]);
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
                    { "mData": "last_sign_at", "bSortable": true,
                        "mRender": function (data) {
                            return (data == "0001-01-01T00:00:00Z") ? "-" : data;
                        }},
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
                    aoData.push({ name: "filter_status", value: selStatus});
                    aoData.push({ name: "filter_certified", value: selCertified});
                }
            });


            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            //$('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown

            $("#e1").select2({  placeholder: "选择用户状态"  });
            $("#e2").select2({   placeholder: "选择认证状态"  });
        }
    };
}();

$(function () {
    TableUsers.init();
})