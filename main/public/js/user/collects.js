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
                    { "mData": "product_id", "bSortable": true, "mRender": function (data, t, row) {
                        var str = '<a href="/products/p/{{id}}" target="_blank"><img src="/product/image?file={{id}}.jpg" style="width: 60px;height: 60px">{{name}}</a>';
                        return S(str).template({id: data, name: row["name"]}).s
                    } },
                    {"mData": "price", "mRender": function (price, t, row) {
                        var currentPrice = row["current_price"];
                        return ( price != currentPrice ? price + "(<b>现价:" + currentPrice + "</b>)" : price);
                    }},
                    { "mData": "created_at", "bSortable": true, "mRender": yunshang.mRenderTime }
                ],
                "fnServerParams": function (aoData) {
                    aoData.push({ name: "filter_status", value: ractive.selStatus || ""});
                }
            }));

            $('#sample_1_wrapper .dataTables_filter input').addClass("m-wrap medium"); // modify table search input
            $('#sample_1_wrapper .dataTables_length select').addClass("m-wrap small"); // modify table per page dropdown
            $('#sample_1_wrapper .dataTables_length select').select2(); // initialzie select2 dropdown

            var DatatableToolBar = yunshang.GetDatatableToolBar();
            ractive = new DatatableToolBar({
                table: {
                    instance: sampleTable,
                    id: "sample_1"
                }
            });

            ractive.on({
                    "delete": function () {
                        doAjaxPost(deleteUrl + "?id=" + ractive.getSelectedData()[0].id, function () {
                            ractive.refreshTable();
                        });
                    },
                    "add-cart":function() {
                        doAjaxPost(addToCartURL + "?productId=" + ractive.getSelectedData()[0].product_id, function () {
                            ractive.refreshTable();
                        });
                    }
                }
            );
        }
    };
}();

$(function () {
    TheTable.init();
})