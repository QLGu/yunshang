$(function () {
    $(".btn-tag").click(function (e) {
        var $this = $(this);
        var tag = "#" + $this.text();
        var tags = $("input[name='roles']").val();
        if (tags.indexOf(tag) == -1) {
            $("input[name='roles']").val(tags + " " + tag);
        }
    });
});
