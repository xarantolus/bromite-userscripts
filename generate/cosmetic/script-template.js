// ==UserScript==
// @name         Cosmetic Ad Block for Bromite
// @namespace    xarantolus
// @version      {{.version}}
// @description  Blocks annoying elements in pages, sourced from many different filter lists
// @author       xarantolus
// @match        *://*/*
// @grant        none
// @run-at       document-end
// ==/UserScript==

var log = function (...data) {
    console.log("[Cosmetic filters by xarantolus (v{{.version}})]:", ...data);
}


// injectStyle injects the css text at the end of the body
function injectStyle(cssStyle) {
    var style = document.createElement('style');
    style.type = 'text/css';
    style.innerHTML = cssStyle;
    document.getElementsByTagName('body')[0].appendChild(style);
}

var rules = {{.rules}};
var defaultRules = rules[""];


function getRules(host) {
    var domainSplit = host.split(".");

    var output = [];

    for (let i = 0; i < domainSplit.length - 1; i++) {
        var domain = domainSplit.slice(i, domainSplit.length).join(".").toLowerCase();

        log("Checking if we got a rule for", domain);

        var rule = rules[domain];
        if (rule) {
            log("Found a rule for domain", domain);
            output.push(rule);
        }
    }

    output.push(defaultRules);

    return output.join(",");
}

var hideFilter = "{display:none !important; height:0 !important; z-index:-99999 !important; visibility:hidden !important; width:0 !important; overflow:hidden !important}"

var cssRule = getRules(location.host) + hideFilter;

log("Injecting style...")
injectStyle(cssRule);
