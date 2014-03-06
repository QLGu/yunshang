$(function(){
    for(var i=0; i<_errorFields.length; i++) {
        var fieldName = _errorFields[i];
        $("input[name=" + fieldName + "]").focus();
    }
});