// ==UserScript==
// @name         Ad Block: Twitter
// @namespace    xarantolus
// @version      0.0.6
// @description  Removes sponsored tweets on Twitter
// @author       xarantolus
// @match        *://twitter.com/*
// @match        *://*.twitter.com/*
// @grant        none
// @run-at       document-end
// ==/UserScript==

var log = function (...data) {
    console.log("[Ad Block: Twitter]:", ...data);
}

var scriptFun = function () {
    // Array source: https://github.com/jodylecompte/twitter-adblock-chrome (MIT Licensed)
    // The .filter at the end is there because of a Bromite bug (https://github.com/bromite/bromite/issues/792#issuecomment-974766145)
    // that replaces multi-byte characters with empty characters. An empty entry would lead to all trends being filtered out, which is not what we want 
    var sponsoredTranslations = [
        "مُروَّج", // Arabic - العربية
        "مُروَّج", // Arabic (Feminine) - العربية (مؤنث)
        "বিজ্ঞাপিত", // Bangla - বাংলা
        "Sustatua", // Basque (beta) - euskara
        "Promoted", // British English
        "Спонсорирано", // Bulgarian - български
        "Patrocinat", // Catalan - català
        "Sponzorirano", // Croatian - hrvatski
        "Sponzorováno", // Czech - čeština
        "Promoveret", // Danish - dansk
        "Promoted", // Dutch - Nederlands
        "Promoted", // English
        "Promoted", // Filipino
        "Mainostettu", // Finnish - suomi
        "Sponsorisé", // French - français
        "Patrocinado", // Galician (beta) - galego
        "Gesponsert", // German - Deutsch
        "Προωθημένο", // Greek - Ελληνικά
        "પ્રચાર કરાયેલું", // Gujarati - ગુજરાતી
        "מקודם", // Hebrew - עברית
        "प्रचारित", // Hindi - हिन्दी
        "Ajánlott", // Hungarian - magyar
        "Dipromosikan", // Indonesian - Indonesia
        "Urraithe", // Irish (beta) - Gaeilge
        "Sponsorizzato", // Italian - italiano
        "プロモーション", // Japanese - 日本語
        "ಪ್ರಾಯೋಜಿತ", // Kannada - ಕನ್ನಡ
        "프로모션 중", // Korean - 한국어
        "Dipromosikan", // Malay - Melayu
        "प्रमोटेड", // Marathi - मराठी
        "Promotert", // Norwegian - norsk
        "تبلیغی", // Persian - فارسی
        "Promowane", // Polish - polski
        "Promovido", // Portuguese - português
        "Promovat", // Romanian - română
        "Реклама", // Russian - русский
        "Промовисано", // Serbian - српски
        "推廣", // Simplified Chinese - 简体中文
        "Sponzorovaný", // Slovak - slovenčina
        "Promocionado", // Spanish - español
        "Sponsrad", // Swedish - svenska
        "விளம்பரப்படுத்தப்பட்டது", // Tamil - தமிழ்
        "ประชาสัมพันธ์", // Thai - ไทย
        "推廣", // Traditional Chinese - 繁體中文
        "Sponsorlu", // Turkish - Türkçe
        "Реклама", // Ukrainian - українська
        "تشہیر شدہ", // Urdu (beta) - اردو
        "Được quảng bá", // Vietnamese - Tiếng Việt
        "Sponsored by", // English, attempt to start closing in on remaining 1% of ads
    ].filter(x => x.trim().length != 0);


    // Source: http://ryanmorr.com/using-mutation-observers-to-watch-for-element-availability/
    (function (win) {
        'use strict';

        var listeners = [],
            doc = win.document,
            MutationObserver = win.MutationObserver || win.WebKitMutationObserver,
            observer;

        function ready(selector, fn) {
            // Store the selector and callback to be monitored
            listeners.push({
                selector: selector,
                fn: fn
            });
            if (!observer) {
                // Watch for changes in the document
                observer = new MutationObserver(check);
                observer.observe(doc.documentElement, {
                    childList: true,
                    subtree: true
                });
            }
            // Check if the element is currently in the DOM
            check();
        }

        function check() {
            // Check the DOM for elements matching a stored selector
            for (var i = 0, len = listeners.length, listener, elements; i < len; i++) {
                listener = listeners[i];
                // Query for elements matching the specified selector
                elements = doc.querySelectorAll(listener.selector);
                for (var j = 0, jLen = elements.length, element; j < jLen; j++) {
                    element = elements[j];
                    // Make sure the callback isn't invoked with the 
                    // same element more than once
                    if (!element.ready) {
                        element.ready = true;
                        // Invoke the callback with the element
                        listener.fn.call(element, element);
                    }
                }
            }
        }

        // Expose `ready`
        win.ready = ready;

    })(this);

    // soundsLikeAd returns if the given text is likely an ad
    function soundsLikeAd(text) {
        return sponsoredTranslations.some(x => text.includes(x))
    }
    // removeIfAd returns a function that gets an HTML element and removes it if it's an ad
    function removeIfAd(adType) {
        var remove = removeWithReason(adType);
        return function (element) {
            if (soundsLikeAd(element.innerText)) {
                remove(element);
            }
        }
    }
    // removeWithReason returns a function that removes an element, logging why it was removed
    function removeWithReason(adType) {
        return function (element) {
            element.remove();
            log("Removed " + adType);
        }
    }

    // Whenever an ad tweet is added to the timeline, we remove it
    // Video players also have this id, so we only remove those that seem like sponsored stuff
    ready('[data-testid="placementTracking"]', removeIfAd("ad tweet"));

    // On profiles, there's stuff with "Promoted tweet" headers left over
    ready("div > h2", removeIfAd("promoted section header"));

    // Whenever a banner ad is added at the top of the "trending" section, we remove it
    ready('[data-testid="eventHero"]', removeWithReason("trends banner ad"));

    // We also want to remove sponsored trends
    ready('[data-testid="trend"]', removeIfAd("sponsored trend"));

    // The same for sponsored follow suggestions
    ready('[data-testid="UserCell"]', removeIfAd("sponsored follow suggestion"));

    log("All listeners attached");
};

if (document.readyState == 'complete') {
    scriptFun();
    log("Ran on completed document");
} else {
    window.addEventListener('load', scriptFun);
    log("Registered as 'load' event listener");
}
