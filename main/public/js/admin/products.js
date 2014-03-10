var TheTable = function () {
    var selStatus = "";
    var selCertified = "";
    return {
        //main function to initiate the module
        init: function () {
            if (!jQuery().dataTable) {
                return;
            }
            // begin first table
            var sampleTable = $('#sample_1').dataTable({
                "bProcessing": true,
                "bServerSide": true,
                "sAjaxSource": productsDataUrl,
                //"sDom":'T<"clear">lfrtip',
                "sDom": "T<'row-fluid'<'span3'l><'span3'r><'span6'f>>t<'row-fluid'<'span6'i><'span6'p>>",
                "oTableTools": {
                    "sSwfPath": "/public/media/swf/copy_csv_xls_pdf.swf",
                    "sRowSelect": "single",
                    "fnRowSelected": function (nodes) {
                        ractive.set("selected", true);

                        var oTT = TableTools.fnGetInstance('sample_1');
                        var aData = oTT.fnGetSelectedData();
                        var enabled = aData[0].enabled;

                        ractive.set("enabled", enabled);
                    },
                    "aButtons": [
                        "copy",
                        "print",
                        {
                            "sExtends": "collection",
                            "sButtonText": "Save",
                            "aButtons": [ "csv", "xls", "pdf" ]
                        }
                    ]
                },
                "aaSorting": [
                    [0, 'desc']
                ],
                "aoColumns": [
                    { "mData": "id", "bSortable": true, "asSorting": [ "desc", "asc" ] },
                    { "mData": "code", "bSortable": false},
                    { "mData": "name", "bSortable": false},
                    { "mData": "model", "bSortable": false},
                    { "mData": "created_at", "bSortable": false},
                    { "mData": "enabled_at", "bSortable": false},
                    { "mData": "unenabled_at", "bSortable": false},
                    { "mData": "enabled", "bSortable": false,
                        "mRender": function (data, type, full) {
                            if (data) {
                                return '<span class="label label-success">已上架</span>';
                            } else {
                                return '<span class="label label-warn">未上架</span>';
                            }
                        }
                    }
                ],
                "aLengthMenu": [
                    [2, 10, 20, 30, 50, -1],
                    ["2条", "10条", "20条", "30条", "50条", "全部"] // change per page values here
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
                },
                "fnRowCallback": function (nRow, aData, iDisplayIndex) {
                    //console.log(JSON.stringify(aData));
                    return nRow;
                }
            });

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            //$('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown


            $("#sample_editable_1_refresh").click(function () {
                sampleTable.fnDraw(true);
            });

            $("#sample_editable_1_enabled").click(function () {
                var oTT = TableTools.fnGetInstance('sample_1');
                var aData = oTT.fnGetSelectedData();
                var url = changeStatusUrl + "?id=" + aData[0].id;
                doAjaxPost(url, function () {
                    sampleTable.fnDraw(true);
                });
            });

            Ractive.delimiters = [ '[[', ']]' ];
            Ractive.tripleDelimiters = [ '[[[', ']]]' ];
            ractive = new Ractive({
                el: "table_tools",
                template: "#table_tools_template",
                data: {
                    selected: false,
                    enabled: "default"
                }
            });
            ractive.on("new",function(){
                $.fancybox.open({
                    href: newProductUrl,
                    type: 'iframe',
                    padding: 5,
                    afterClose: function (e) {
                        sampleTable.fnDraw(true);
                    }
                });
            });

            $("#e1").select2({
                placeholder: "选择产品状态"
            });

            $("#e1").on("change", function (e) {
                selStatus = e.val;
                sampleTable.fnDraw(true);
            });
        }
    };
}();

$(function () {
    TheTable.init();
})