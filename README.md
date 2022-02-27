# Bromite User Scripts
This is a repository that hosts [User Scripts for Bromite](https://github.com/bromite/bromite/wiki/UserScripts), an Android browser. They might also work on other mobile or desktop browsers, but I only use them on [Bromite](https://www.bromite.org/).

### Scripts
Here scripts are sorted by category. See below on how to install them into your browser.

#### Blockers
Bromite has a built-in ad blocker (also see my [Bromite ad blocking project](https://github.com/xarantolus/filtrite)), but some sites are very good at bypassing it. We can improve blocking on these sites using scripts.

* [**Twitter**](https://github.com/xarantolus/bromite-userscripts/releases/latest/download/twitter.user.js): block ads on Twitter (sponsored tweets and sponsored trends)
* [**I don't care about cookies**](https://github.com/xarantolus/bromite-userscripts/releases/latest/download/idcac.user.js): block all kinds of cookie prompts. 
  * This script is based on the ["I don't care about cookies" browser extension](https://addons.mozilla.org/de/firefox/addon/i-dont-care-about-cookies/) (GPL). 
  * This script is automatically regenerated from time to time, keeping up to date with the latest rules from the browser extension
  * **Security consideration**: if the author of the browser extension inserts malicious code, this script would likely also contain that code
* [**Cosmetic AdBlock**](https://github.com/xarantolus/bromite-userscripts/releases/latest/download/cosmetic.user.js): block annoying elements 
  * The Bromite AdBlock engine does not support cosmetic filtering, so this script implements that capability (to a very basic extent)
  * Rules are regenerated once a week from the filter lists defined in [this file](generate/cosmetic/filter-lists.txt)


---

### Installing a script 
Make sure you have a recent version of [Bromite](https://www.bromite.org/) installed. Then you can follow these steps:
1. At first you need to download the script file. You can do this for the scripts in this repository by holding on the link until the menu appears, then selecting "Download link"
2. Now you can go to Bromite settings (three dots at the top right, then Settings)
3. Scroll down to open the "User Scripts" section under the "Advanced" menu
4. Make sure "Activate User Scripts" is enabled/on
5. Select the "Add script" button
6. Now select the file you just downloaded (likely in your downloads directory)
7. Confirm the installation
8. (Optionally) Click the "View source" button to verify the content of the script
9. Enable the newly installed script using the switch at the left side



### Contributing / Creating your own scripts
You can also create your own scripts, see the [official documentation](https://github.com/bromite/bromite/wiki/UserScripts) on how to do that. You should read the [Chromium User Script Design Document](https://www.chromium.org/developers/design-documents/user-scripts) to learn about existing pitfalls.

I recommend trying out [remote debugging](https://developer.chrome.com/docs/devtools/remote-debugging/) via a desktop browser if your script doesn't behave as expected.

If you want to remove certain elements on dynamic pages (like Twitter), I recommend [this snippet of code](http://ryanmorr.com/using-mutation-observers-to-watch-for-element-availability/), it's very helpful.


### Issues & Contributing
If you have any issues with these scripts (e.g. some ads aren't blocked for an ad blocking script), please feel free to open an issue. Also if you want to add something, feel free to do that :)


### [License](LICENSE)
All scripts unless otherwise noted are published under the MIT License (see the LICENSE file). Some scripts might be licensed differently (e.g. because they are derived from GPL-licensed works), which is indicated by the license header at the top of the file
