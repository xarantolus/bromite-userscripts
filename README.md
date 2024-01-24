# [Bromite User Scripts](https://userscripts.010.one)
This repository provides [User Scripts](https://github.com/bromite/bromite/wiki/UserScripts) for [Bromite](https://www.bromite.org/) and [Cromite](https://www.cromite.org/), which are browsers for Android. The scripts might also work on other mobile or desktop browsers.

Bromite has a built-in ad blocker (also see my [Bromite ad blocking project](https://github.com/xarantolus/filtrite)), but some sites are very good at bypassing it. We can improve blocking on these sites using scripts. Cromite already has a good built-in Ad Blocker.

### Downloads
You can get an overview and download links on [**the website**](https://userscripts.010.one).

**Acknowledgement**: The filter lists for the cosmetic filter script are sourced from [collinbarrett's FilterLists](https://github.com/collinbarrett/FilterLists) project (see [filterlists.com](https://filterlists.com)). Thanks to everyone that contributes! :)


**Deprecated**
* [**Twitter**](https://github.com/xarantolus/bromite-userscripts/releases/latest/download/twitter.user.js): block ads on Twitter (sponsored tweets, trends etc.)
  * This script is no longer maintained as I no longer use Twitter - please feel free to contribute the necessary changes if it no longer works


---

### Installing a script
Make sure you have a recent version of [Bromite](https://www.bromite.org/) installed. Then you can follow these steps:
1. At first you need to download a script file from the project website
2. Click the "Open" button that appears once the script has been downloaded
3. Confirm the installation
4. Make sure "Activate User Scripts" is enabled/on
5. Enable the newly installed script using the switch at the left side

You might need to go to settings (via the three dots at the top right), then "User Scripts" to enable the "Activate User Scripts" option first.

If the browser doesn't prompt you to install the script, you can also just go to the User Script settings and add the file manually.

---

### Auto-generated scripts
Some scripts are auto-generated (because they need to be regenerated from time to time to include up to date sources). The source code for the generators is in subdirectories of the [`generate`](generate/) directory.

You can see statistics (e.g. number of included rules) in the [latest release](https://github.com/xarantolus/bromite-userscripts/releases/latest).

### Creating your own scripts
You can also create your own scripts, see the [official documentation](https://github.com/bromite/bromite/wiki/UserScripts) on how to do that. You should read the [Chromium User Script Design Document](https://www.chromium.org/developers/design-documents/user-scripts) to learn about existing pitfalls.

I recommend trying out [remote debugging](https://developer.chrome.com/docs/devtools/remote-debugging/) via a desktop browser if your script doesn't behave as expected.

If you want to remove certain elements on dynamic pages (like Twitter), I recommend [this snippet of code](http://ryanmorr.com/using-mutation-observers-to-watch-for-element-availability/), it's very helpful.


### Issues & Contributing
If you have any issues with these scripts (e.g. some ads aren't blocked for an ad blocking script), please feel free to open an issue. Also if you want to add something, feel free to do that :)


### [License](LICENSE)
All scripts unless otherwise noted are published under the MIT License (see the LICENSE file). Some scripts might be licensed differently (e.g. because they are derived from GPL-licensed works), which is indicated by the license header at the top of the file
