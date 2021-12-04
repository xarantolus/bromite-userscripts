// ==UserScript==
// @name         Generalized Cookie Prompt Accept/Block
// @namespace    xarantolus
// @version      0.0.2
// @description  Automatically accepts many kinds of cookie prompts 
// @author       xarantolus
// @match        *://*/*
// @grant        none
// @run-at       document-end
// ==/UserScript==

var log = function (...data) {
    console.log("[Generalized Cookie Prompt Accept/Block]:", ...data);
}

var scriptFun = function () {
    var scrollFixInjected = false;
    function injectScrollFix() {
        if (scrollFixInjected) {
            return;
        }
        scrollFixInjected = true;

        var style = document.createElement('style');
        style.type = 'text/css';
        style.innerHTML = `body,html {overflow:auto !important}`;
        document.getElementsByTagName('head')[0].appendChild(style);
    }

    function querySelectorAllPages(selector) {
        var output = [...document.querySelectorAll(selector)];

        var iframes = [...document.querySelectorAll("iframe")];

        // TODO: Sometimes when accessing iframes, an exception is thrown
        try {
            iframes.forEach(frame => {
                output.concat(frame.contentWindow.document.body.querySelectorAll(selector));
            })
        } catch (ex) {
            log("Cannot access iframes:", ex);
         }

        return output;
    }

    function DOMRegexClassApply(regex, callback) {
        for (let i of querySelectorAllPages('*')) {
            if (regex.test(i.className)) {
                callback(i);
            }
        }
    }

    var acceptKeywords = [
        // English
        "accept", "agree", "continue", "allow",

        // German
        "akzeptieren", "zustimmen", "weiter", "erlauben",
    ]

    function isAcceptButton(elem) {
        var btntxt = elem.innerText.toLowerCase();

        return acceptKeywords.some(kw => btntxt.includes(kw));
    }



    function acceptOrBlock(element) {
        var buttons = [...element.querySelectorAll("button")];
        var acceptButton = buttons.find(isAcceptButton);
        if (acceptButton) {
            // Click accept if possible
            acceptButton.click();
            log("Clicked an accept button within an cookie/consent/gdpr element");
            return;
        }

        // Remove the element, but also inject a common fix for scrolling issues
        element.remove();
        log("Removed a cookie/consent/gdpr element");
        injectScrollFix();
    }



    function removeElements() {
        querySelectorAllPages("button").filter(x => isAcceptButton(x)).forEach(x => {
            x.click();
            log("Clicked an accept button", x);
        });

        // e.g. https://www.nytimes.com/
        DOMRegexClassApply(/gdpr/, acceptOrBlock);
        DOMRegexClassApply(/consent/, acceptOrBlock);
        DOMRegexClassApply(/cookie/, acceptOrBlock);
    }


    removeElements();

    var observer = new MutationObserver(removeElements);
    observer.observe(this.document.documentElement, {
        childList: true,
        subtree: true
    });

    setInterval(removeElements, 2500);

    log("Finished initialization");
}

if (document.readyState == 'complete') {
    scriptFun();
    log("Ran on completed document");
} else {
    window.addEventListener('load', scriptFun);
    log("Registered as 'load' event listener");
}

