$(function () {
    $("a.log-btn2").click(function () {
        $('#regForm').submit();
    });

    $("form input[name=validateCode]").on('keyup', function (e) {
        if (event.which == 13 || event.keyCode == 13) {
            e.preventDefault();
            $('#regForm').submit();
        }
    });

    if(typeof(_errorFields) == 'undefined'){
        $("#regForm input[name=email]").focus();
    }
});