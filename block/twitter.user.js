// ==UserScript==
// @name         Ad Block: Twitter
// @namespace    google
// @version      0.0.1
// @description  Removes sponsored tweets on Twitter
// @author       xarantolus
// @match        *://*.twitter.com/*
// @match        *://twitter.com/*
// @grant        none
// @run-at       document-start
// ==/UserScript==

// This is basically a copy of the chrome extension https://github.com/jodylecompte/twitter-adblock-chrome (MIT License), but in userscript form

// ------------------------------------------------------------------------------------------------------------------------

// Some languages share the english text of "Promoted" such as Fillipino
// Leaving duplicates in place in case they change later
const languages = [
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
];

const hidePromotedTweets = () => {
    Array.from(document.querySelectorAll("span"))
        .filter((el) => languages.includes(el.textContent))
        .forEach((el) => {
            const parent = el.closest("div[data-testid=placementTracking]");

            if (parent) {
                el.closest("div[data-testid=placementTracking]").remove();
            }
        });
};

const targetNode = document.querySelector("body");
const config = { attributes: false, childList: true, subtree: true };

const observer = new MutationObserver((mutationsList, observer) => {
    hidePromotedTweets();
});

observer.observe(targetNode, config);

document.addEventListener("load", () => hidePromotedTweets());
