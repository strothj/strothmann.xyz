  * `latest` latest public (as described here) [![](https://badge.imagelayers.io/wernight/plex-media-server:latest.svg)](https://imagelayers.io/?images=wernight/plex-media-server:latest 'Get your own badge on imagelayers.io')
  * `autoupdate` installs latest on start (see below for differences) [![](https://badge.imagelayers.io/wernight/plex-media-server:autoupdate.svg)](https://imagelayers.io/?images=wernight/plex-media-server:autoupdate 'Get your own badge on imagelayers.io')
  * `0`, `0.9`, `0.9.14`, `0.9.14.6` (or similar) are like `latest` but for a specific version [![](https://badge.imagelayers.io/wernight/plex-media-server:0.svg)](https://imagelayers.io/?images=wernight/plex-media-server:0 'Get your own badge on imagelayers.io')

Dockerized [Plex Media Server](https://plex.tv/): Plex organizes your video, music, and photo collections and streams them to all of your screens (mobile, TV/Chromecast, laptop...).


### Usage

It is recommended to provide two mount points writable by user `797` (that `plex` random UID inside the container for safety, alternatively use `--user` flag):

  * `/config`: To somewhere to hold your Plex configuration (can be a data-only container). This will include all media listing, posters, collections and playlists you've setup...
  * `/media`: To one or more of your media files (videos, audio, images...).

Example:

    $ mkdir ~/plex-config
    $ chown 797:797 -R ~/plex-config
    $ docker run -d --restart=always -v ~/plex-config:/config -v ~/Movies:/media --net=host -p 32400:32400 wernight/plex-media-server

Once done, wait a few seconds and open `http://localhost:32400/web` in your browser.

The flag `--net=host` is only required for the first run, so that your can login locally without password (without SSH proxy) and see the "Server" tab in the web UI (see troubleshooting section below). If you want **Avahi broadcast** to work then keep `--net=host` even after being logged in, but this will be somewhat less secure.

The `--restart=always` is optional, it'll for example allow auto-start on boot.


### Features

  * **Small**: Built using official Docker [Debian](https://registry.hub.docker.com/_/debian/) and official [Plex download](https://plex.tv/downloads) (takes 85 MB instead of 180 MB for Ubuntu).
  * **Simple**: One command and you should be ready to go. All documented here.
  * **Secure**:
      * Runs Plex as `plex` user (not root as [Docker's Containers don't contain](http://www.projectatomic.io/blog/2014/09/yet-another-reason-containers-don-t-contain-kernel-keyrings/)).
      * Downloads and installs the official binaries.
      * Avoids [PID 1 / zombie reap problem](https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/) (if plex or one of its subprocesses dies) by running directly plex.

#### Comparison of main Plex Docker containers

Image                        | Size                 | [Runs As]  | [PID 1 Reap] | [Slim Container] | [Plex Pass]
---------------------------- | -------------------- | ---------- | ------------ | ---------------- | -----------
[wernight/plex-media-server] | ![][img-wernight]    | **user**   | **Safe**     | **Yes**          | **Supported**
[linuxserver/plex]           | ![][img-linuxserver] | **user**   | **Safe**     | No               | No
[timhaak/plex]               | ![][img-timhaak]     | root       | Unsafe       | No               | **Supported**
[needo/plex]                 | ![][img-needo]       | root       | **Safe**     | No               | No
[binhex/arch-plex]           | ![][img-binhex]      | root       | Unsafe       | No               | No


Based on current state as of December 2015 (if you find any mistake please open a ticket on GitHub).

[Runs As]: https://opensource.com/business/14/7/docker-security-selinux
[PID 1 Reap]: https://blog.phusion.nl/2015/01/20/docker-and-the-pid-1-zombie-reaping-problem/
[Slim Container]: https://blog.phusion.nl/2015/01/20/baseimage-docker-fat-containers-treating-containers-vms/
[Plex Pass]: https://support.plex.tv/hc/en-us/articles/201844613-Early-Access-Preview-Releases
[wernight/plex-media-server]: https://registry.hub.docker.com/u/wernight/plex-media-server/
[linuxserver/plex]:           https://registry.hub.docker.com/u/linuxserver/plex/
[timhaak/plex]:               https://registry.hub.docker.com/u/timhaak/plex/
[needo/plex]:                 https://registry.hub.docker.com/u/needo/plex/
[binhex/arch-plex]:           https://registry.hub.docker.com/u/binhex/arch-plex/

### Upgrades and Versions

*Plex Media Server* does *not* support auto-upgrade from the UI on Linux. If/once it does, we'd be more than happy to support it.

There are two ways to keep up to date:

  * Using `wernight/plex-media-server:latest` (default) – To upgrade to the latest public version do again a `docker pull wernight/plex-media-server` and restart your container; that should be it. You may use a *tagged version* to use a fixed or older version as well. It works as described here.
  * Using `wernight/plex-media-server:autoupdate` (for users who want the really latest) – Installs the latest public or **Plex Pass** release each time the container starts. It has a few differences compared to what is described here:
      * Runs as `root` initially so it can install Plex (required), after that it runs as `plex` user.
      * Supports PlexPass: Premium users get to download newer versions shortly before they get public. For that set two additional environment variables (only be used to retrieve the latest official download URL and cleared after that) like:

            $ docker run -d --restart=always -v ~/plex-config:/config -v ~/Movies:/media --net=host -p 32400:32400 -e PLEXPASS_LOGIN='<my_plex_login>' -e PLEXPASS_PASSWORD='<my_plex_password>' wernight/plex-media-server:autoupdate


### Environment Variables

You can change some settings by setting environement variables:

  * `PLEX_MEDIA_SERVER_MAX_STACK_SIZE` ulimit stack size (default: 3000).
  * `PLEX_MEDIA_SERVER_MAX_PLUGIN_PROCS` the number of plugins that can run at the same time (default: 6).
  * `PLEXPASS_LOGIN` your Plex Pass username or e-mail (used only on the `:autoupdate` tagged image).
  * `PLEXPASS_PASSWORD` your Plex Pass username or e-mail (used only on the `:autoupdate` tagged image).


### Troubleshooting

  * I have to accept EULA each time?!
      * Did you forget to mount `/config` directory? Check also that it's writable by user `797`.
  * Cannot see [**Server** tab](http://localhost:32400/web/index.html#!/settings/server) from settings!
      * Try running once with `--net=host`. You may allow more IPs without being logged in by then going to Plex Settings > Server > Network > List of networks that are allowed without auth; or edit `your_config_location/Plex Media Server/Preferences.xml` and add `allowedNetworks="192.168.1.0/255.255.255.0"` attribute the `<Preferences …>` node or what ever your local range is.
  * Why do I have a random server name each time?
      * Either set a friendly name undex Plex Settings > Server > General; or start with `-h some-name`.
  * Which port do I need to open on my firewall/router?
      * Even if you're using `--net=host` or `--port 0.0.0.0:32400:32400` flag, you'll still need to redirect port 32400 on your router to your machine running Plex, else you'll only be able to access it from within your LAN and you won't be able to Chromecast and other things. Remember to also check your firewall. Note that you can use another port if you so desire.

### Backup

Honestly I wish there was a more official documentation for this. What you really need to back-up (adapt `~/plex-config` to
your `/config` mounting point):

  * Your media, obviously!
  * `~/plex-config/Plex Media Server/Media/`
  * `~/plex-config/Plex Media Server/Metadata/`
  * `~/plex-config/Plex Media Server/Plug-in Support/Databases/`

In practice, you may want to be safer and back-up everything except may be `~/plex-config/Plex Media Server/Cache/`
which is pretty large and you can really just skip it. It'll be rebuild with the thumbnails, etc. as you had them.
But don't take my word for it, it's really easy for you to check.


### Feedbacks

Having more issues? [Report a bug on GitHub](https://github.com/wernight/docker-plex-media-server/issues).

[img-wernight]: https://badge.imagelayers.io/wernight/plex-media-server:latest.svg
[img-linuxserver]: https://badge.imagelayers.io/linuxserver/plex:latest.svg
[img-timhaak]: https://badge.imagelayers.io/timhaak/plex:latest.svg
[img-needo]: https://badge.imagelayers.io/needo/plex:latest.svg
[img-binhex]: https://badge.imagelayers.io/binhex/arch-plex:latest.svg

