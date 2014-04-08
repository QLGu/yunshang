function afterLoad() {

    $('#commentForm').ajaxForm({
        dataType: 'json',
        success: function (ret) {
            alert(ret.message);
        }
    });
}

function _loadComments() {
    var url = page.reloadURL();//+"&limit=1";
    $('div.BlockDiscuss').load(url, afterLoad);
}

$(function () {
    _loadComments();
});

function onPage(p) {
    page.set("page", p);
    _loadComments();
}
