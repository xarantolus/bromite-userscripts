# bromite-userscripts
This is a repository that hosts [User Scripts for Bromite](https://github.com/bromite/bromite/wiki/UserScripts), an Android browser.

### Scripts
Here scripts are sorted by category. See below on how to install them into your browser.

#### Ad Blockers
Bromite has a built-in ad blocker (also see my [Bromite ad blocking project](https://github.com/xarantolus/filtrite)), but some sites are very good at bypassing it. We can improve blocking on these sites using scripts.

* [**Twitter**](block/twitter.user.js?raw=true): block ads on Twitter (sponsored tweets and sponsored trends)

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
You can also create your own scripts, see the [official documentation](https://github.com/bromite/bromite/wiki/UserScripts) on how to do that.

I recommend trying out [remote debugging](https://developer.chrome.com/docs/devtools/remote-debugging/) via a desktop browser if your script doesn't behave as expected.

If you want to remove certain elements on dynamic pages (like Twitter), I recommend [this snippet of code](http://ryanmorr.com/using-mutation-observers-to-watch-for-element-availability/), it's very helpful.


### Issues & Contributing
If you have any issues with these scripts (e.g. some ads aren't blocked for an ad blocking script), please feel free to open an issue. Also if you want to add something, feel free to do that :)


### [License](LICENSE)
This is free as in freedom software. Do whatever you like with it.
