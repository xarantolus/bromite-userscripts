// Copyright © 2022 xarantolus
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see https://www.gnu.org/licenses/gpl-3.0.html.
//
// --------------------------------------------------------------------------------
//
// ==UserScript==
// @name         I don't care about cookies
// @namespace    xarantolus
// @version      {{.version}}
// @description  Removes cookie banners, based on the browser extension
// @author       xarantolus
// @match        *://*/*
// @grant        none
// @run-at       document-end
// ==/UserScript==
/// @stats {{.statistics}}


var commons = {{ .commons }};
var rules = {{ .rules }};
var javascriptFixes = {{ .javascriptFixes }};
var cookieBlockCSS = {{ .cookieBlockCSS }};

var log = function (...data) {
    console.log("[I don't care about cookies (v{{.version}})]:", ...data);
}

var injectedCookieBlockCSS = false;

var scriptFun = function () {

    // injectStyle injects the css text at the end of the body
    function injectStyle(cssStyle) {
        var style = document.createElement('style');
        style.type = 'text/css';
        style.innerHTML = cssStyle;
        document.getElementsByTagName('body')[0].appendChild(style);
    }

    // findRules finds all matching rule for this host
    function findRules(host) {
        // Try to find rules for this host. E.g. 
        // if the domain is "sub.domain.com" it will first try
        // "sub.domain.com", then "domain.com"

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

        return output;
    }

    // styleForRule returns
    function styleForRule(rule) {
        var style = "";
        if (rule["s"]) {
            // Specialized style
            style += rule["s"];
        }

        if (rule["c"]) {
            style += "\n" + commons[rule["c"]] || "";
        }

        return style;
    }

    function scriptForRule(rule) {
        if (rule["j"]) {
            log("Rule requests script", rule["j"]);
            return javascriptFixes["" + rule["j"]];
        }
        return null;
    }

    // Only do this on the first run of the function
    if (!injectedCookieBlockCSS) {
        injectStyle(cookieBlockCSS);
        injectedCookieBlockCSS = true;
        log("Injected common CSS rules");
    }

    var result = findRules(location.host);
    log("Found", result.length, "rule(s) for", location.host);
    if (rules.length == 0) {
        return;
    }

    for (var i = 0; i < result.length; i++) {
        var r = result[i];

        log("Rule:", r);

        var css = styleForRule(r);

        if (css) {
            injectStyle(css);
            log("Computed style:", css);
        } else {
            log("No style to inject")
        }

        var js = scriptForRule(r);
        if (js) {
            js();
            log("Injected script");
        }
    }
}

log("Running on inject");
scriptFun();

if (document.readyState !== 'complete') {
    window.addEventListener('load', scriptFun);
    log("Registered as 'load' event listener. This means that the injection of rules will happen twice (else cookie prompts that haven't loaded yet could be missed)");
}
