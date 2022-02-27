// injectStyle injects the css text at the end of the body
function injectStyle(cssStyle) {
    var style = document.createElement('style');
    style.type = 'text/css';
    style.innerHTML = cssStyle;
    document.getElementsByTagName('body')[0].appendChild(style);
}

var rules = {{.rules}};
var defaultRules = rules[""];


function getRules(domain) {
    var rules = defaultRules;
    
    var elem = rules[domain];
    if (elem) {
        rules += "," + elem;
    }

    return rules;
}

var hideFilter = "{display:none !important; height:0 !important; z-index:-99999 !important; visibility:hidden !important; width:0 !important; overflow:hidden !important}"

var cssRule = getRules(location.hostname) + hideFilter;

injectStyle(cssRule);
