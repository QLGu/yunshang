$(function () {
    function addToList(id, value, text) {
        return $('#' + id).append('<option value="' + value + '">' + text + '</option>');
    }

    var cityHolder = 0;
    var cityHolder2 = 0;

    cityHolder = $('#holder').ChinaCitySelect({'prov': '#state', 'city': '#city', 'dist': '#district', url: '/public/libs/city/city.json'});
    cityHolder2 = $('#holder').ChinaCitySelect({url: '/public/libs/city/city.json'});

    setTimeout(function () {
        _p && $('#state').val(_p).change();
        _c && $('#city').val(_c).change();
        _a && $('#district').val(_a).change();
    }, 1000);
});