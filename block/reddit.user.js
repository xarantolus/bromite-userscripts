// ==UserScript==
// @name         Ad Block: Mobile Reddit
// @namespace    xarantolus
// @version      0.0.3
// @description  Removes sponsored posts, ads and recommended posts on Reddit
// @author       xarantolus
// @match        *://reddit.com/*
// @match        *://*.reddit.com/*
// @grant        none
// @run-at       document-end
// @homepage     https://userscripts.010.one
// @url_source   https://github.com/xarantolus/bromite-userscripts/releases/latest/download/reddit.user.js
// ==/UserScript==

var log = function (...data) {
  console.log("[Ad Block: Reddit]:", ...data);
};

var scriptFun = function () {
  // Source: http://ryanmorr.com/using-mutation-observers-to-watch-for-element-availability/
  (function (win) {
    "use strict";

    var listeners = [],
      doc = win.document,
      MutationObserver = win.MutationObserver || win.WebKitMutationObserver,
      observer;

    function ready(selector, fn) {
      // Store the selector and callback to be monitored
      listeners.push({
        selector: selector,
        fn: fn,
      });
      if (!observer) {
        // Watch for changes in the document
        observer = new MutationObserver(check);
        observer.observe(doc.documentElement, {
          childList: true,
          subtree: true,
        });
      }
      // Check if the element is currently in the DOM
      check();
    }

    function check() {
      // Check the DOM for elements matching a stored selector
      for (
        var i = 0, len = listeners.length, listener, elements;
        i < len;
        i++
      ) {
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

  function removeWithReason(adType) {
    return function (element) {
      element.remove();
      log("Removed " + adType);
    };
  }

  // removeIfAd returns a function that gets an HTML element and removes it if it's an ad
  function removeAd(adType) {
    var remove = removeWithReason(adType);
    return function (element) {
		remove(element);
    };
  }

  ready("shreddit-ad-post", removeAd("feed ad"));
  ready("shreddit-comments-page-ad", removeAd("comments page ad"));

  // Posts that Reddit injects from communities you are not part of ("recommended" posts). They
  // carry a "Join" button in the post's community bar and a `recommendation-source` attribute
  // (seen: "geo_popular_hot_sort", "global_good_visits", "two_hop_time_on_subreddits_6hr"). A
  // normal subreddit feed has neither (0 of 27 posts in /r/de/), so posts there are left alone.
  var recommendedPosts =
    "shreddit-post[recommendation-source]:not([recommendation-source=''])";

  // A subreddit's own page keeps its posts even when you are not joined to it, so only the
  // aggregate feeds (home, /r/popular, /r/all) are filtered by the Join button.
  var pathSubreddit = (location.pathname.match(/^\/r\/([^/]+)/i) || [])[1];
  var aggregateFeed = !pathSubreddit || /^(popular|all)$/i.test(pathSubreddit);

  // The style rules hide a post the moment the marker appears - also for feed items that are loaded
  // lazily or whose attribute is set after insertion, which the MutationObserver above (childList
  // only) would not catch. Two separate rules, so a browser without :has() only loses the second.
  var css = recommendedPosts + " { display: none !important; }";
  if (aggregateFeed) {
    css += "\nshreddit-post:has(shreddit-join-button) { display: none !important; }";
  }
  var style = document.createElement("style");
  style.textContent = css;
  (document.head || document.documentElement).appendChild(style);

  ready(recommendedPosts, removeWithReason("recommended post"));

  // Fallback for browsers without :has(): remove the post a Join button sits in. The observer sees
  // the button when it is inserted, so lazily rendered cards are caught too.
  if (aggregateFeed) {
    ready("shreddit-join-button", function (button) {
      var post = button.closest("shreddit-post");
      if (post) {
        post.remove();
        log("Removed recommended post (join button)");
      }
    });
  }
};


if (document.readyState == "complete") {
  scriptFun();
  log("Ran on completed document");
} else {
  window.addEventListener("load", scriptFun);
  log("Registered as 'load' event listener");
}
