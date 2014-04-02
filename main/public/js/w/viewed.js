$(function () {
    var its = _.map(YSCookie.viewed_products(), function (pid) {
        pid.url = "/products/p/" + pid.id;
        return S(' <li><div class="prod-text"><a href="{{url}}">{{name}}</a></div></li>').template(pid).s;
    });
    $('#viewed-products').html(its.join(""));

    $('#btn-clear-viewed').click(function () {
        YSCookie.clear_viewed_products();
        $('#viewed-products').html("");
    });
});