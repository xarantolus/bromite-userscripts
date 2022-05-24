// ==UserScript==
// @name         Cosmetic Ad Block for Bromite
// @namespace    xarantolus
// @version      {{.version}}
// @description  Blocks annoying elements in pages, sourced from many different filter lists
// @author       xarantolus
// @match        *://*/*
// @grant        none
// @run-at       document-start
// ==/UserScript==
/// @stats {{.statistics}}
{
    let userscriptLog = function (...data) {
        console.log("[Cosmetic filters by xarantolus (v{{.version}})]:", ...data);
    }


    function injectStyle(cssStyle) {
        let style = document.createElement('style');
        style.type = 'text/css';
        style.innerHTML = cssStyle;
        document.getElementsByTagName('head')[0].appendChild(style);
    }

    let deduplicatedStrings = {{.deduplicatedStrings }};
    let injectionRules = {{.injectionRules }};
    let scriptletRules = {{.scriptletRules}};
    let rules = {{.rules }};
    let defaultRules = rules[""];

    {{ .scriptletDefinitions }}
    let scriptletLookupTable = {{ .scriptletLookup }};


    function getRules(host) {
        let domainSplit = host.split(".");

        let output = [];

        for (let i = 0; i < domainSplit.length - 1; i++) {
            let domain = domainSplit.slice(i, domainSplit.length).join(".").toLowerCase();

            userscriptLog("Checking if we got a rule for", domain);

            let rule = rules[domain];
            if (rule != null) {
                if (typeof rule === 'number') {
                    // the selector is saved at this index in the deduplicatedRules array
                    let realRule = deduplicatedStrings[rule];
                    userscriptLog("Found deduplicated rule", rule, "for domain", domain);
                    output.push({ "s": realRule });
                } else {
                    // It's a string that directly defines the selector
                    userscriptLog("Found normal rule for domain", domain);
                    output.push({ "s": rule });
                }
            }

            let injection = injectionRules[domain];
            if (injection != null) {
                if (typeof injection === 'number') {
                    let realInjection = deduplicatedStrings[injection];
                    userscriptLog("Found deduplicated injection", injection, "for domain", domain);
                    output.push({ "i": realInjection })
                } else {
                    userscriptLog("Found normal injection for domain", domain);
                    output.push({ "i": injection });
                }
            }

            let scriptlets = scriptletRules[domain];
            if (scriptlets != null) {
                userscriptLog("Found " + scriptlets.length + " scriptlet(s) for domain", domain);
                output.push({"scriptlets": scriptlets})
            }
        }

        output.push({ "s": defaultRules, isDefault: true });

        return output;
    }

    let hiddenStyle = "display:none!important;min-height:0!important;height:0!important;z-index:-99999!important;visibility:hidden!important;width:0!important;min-width:0!important;overflow:hidden!important";
    let hideRules = "{" + hiddenStyle + "}"

    let foundRules = getRules(location.host);

    userscriptLog("Found", foundRules.length, "rules to inject");

    let hiddenElementsSelector = foundRules.filter(r => r["s"] != null)
        .map(r => r["s"]).join(",") + hideRules;

    let cssInjections = foundRules.filter(r => r["i"] != null).map(r => r["i"]).join("");
    let scriptlets = foundRules.filter(r => r["scriptlets"] != null).flatMap(r => r["scriptlets"]);

    let pageSpecificSelectors = foundRules.filter(r => r["s"] != null && !r.isDefault)
        .map(r => r["s"]).join(",");

    userscriptLog("Page specific selectors:", (pageSpecificSelectors || "(none)"))

    // Source: https://stackoverflow.com/a/61747276
    function elementReady(selector) {
        return new Promise((resolve, reject) => {
            const el = document.querySelector(selector);
            if (el) { resolve(el); }
            new MutationObserver((mutationRecords, observer) => {
                // Query for elements matching the specified selector
                Array.from(document.querySelectorAll(selector)).forEach((element) => {
                    resolve(element);
                    //Once we have resolved we don't need the observer anymore.
                    observer.disconnect();
                });
            })
                .observe(document.documentElement, {
                    childList: true,
                    subtree: false // This was changed to "false" since we only need "head", a direct descendant of the document element
                });
        });
    }

    function hidePageSpecificElements(reason) {
        if (pageSpecificSelectors.length == 0) return;

        userscriptLog("Searching for elements (" + reason + ")")
        let elems = [...document.querySelectorAll(pageSpecificSelectors)];
        elems.forEach(function (elem) {
            elem.setAttribute("style", hiddenStyle);
        });
        userscriptLog("Tried hiding", elems.length, "page-specific elements");
    }

    function runScriptlets() {
        for (let idx = 0; idx < scriptlets.length; idx++) {
            const scriptletName = scriptlets[idx][0];
            const scriptletArgs = scriptlets[idx].slice(1);

            let scriptletFunction = scriptletLookupTable[scriptletName];
            if (!scriptletFunction) {
                userscriptLog("could not find scriptlet function for " + scriptletName);
                continue;
            }

            userscriptLog("Running scriptlet '" + scriptletName + "' with args " + scriptletArgs);

            // Now actually run the scriptlet function we found
            try {
                let res = scriptletFunction({
                    "name": scriptletName,
                    "args": scriptletArgs,
                    "engine": "",
                    "version": "1.0.0",
                    "verbose": true,
                    "ruleText": "(rule text not available)"
                }, scriptletArgs);
                if (res) {
                    userscriptLog("Scriptlet returned " + res);
                }
            }catch(e) {
                userscriptLog("Running scriptlet: " + e);
            }
        }
    }

    // Now we have hidden a lot of stuff using rules. However, some sites still display elements
    // because they look like <span class="ad" style="display:block">
    // This means that the !important from our css declaration above will not work on these elements (as direct styles take precedence)
    // We need to replace the style of all elements with this selector
    // When the HTML has finished parsing:
    window.addEventListener('DOMContentLoaded', function () {
        hidePageSpecificElements("DOMContentLoaded");

        setTimeout(() => hidePageSpecificElements("DOMContentLoaded + 1000ms"), 1000);
    });
    // And after the page is fully loaded, we do a bunch of checks within the first second or so.
    // If a page pops up a cookie popup after the page has loaded, this one will also defeat it
    window.addEventListener('load', function () {
        hidePageSpecificElements("load - initial");

        function to(offset) {
            let ms = offset * 500;
            setTimeout(() => hidePageSpecificElements("load + " + ms + "ms"), ms);
        };
        for (let i = 1; i <= 5; i++) {
            to(i);
        }
    })

    elementReady('head').then((_) => {
        try {
            injectStyle(hiddenElementsSelector);
            userscriptLog("Injected combined style");
        } catch(e) {
            userscriptLog("Error injecting combined style: " + e);
        }

        if (cssInjections.length > 0) {
            try {
                injectStyle(cssInjections);
                userscriptLog("Also injected additional styles (usually fixes for scrolling issues)")
            } catch (e) {
                userscriptLog("Error injecting additional styles: " + e);
            }
        }

        try {
            runScriptlets();
            userscriptLog("Ran scriptlets");
        } catch (e) {
            userscriptLog("Error running scriptlets: " + e);
        }
    });
}
