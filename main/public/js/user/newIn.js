$(function () {
    $(":radio[value=2]").click(function () {
        $("#more").show();
        $("#more input").addAttr("required");
    });

    $(":radio[value=1]").click(function () {
        $("#more").hide();
        $("#more input").removeAttr("required");
    });

    if ($(":radio[value=2]").prop("checked")) {
        $("#more").show();
        $("#more input").addAttr("required");
    }

    if ($(":radio[value=1]").prop("checked")) {
        $("#more input").removeAttr("required");
    }
});