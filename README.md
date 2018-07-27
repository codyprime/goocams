======================================================
Go Open Cameras - CLI tools for configuring IP cameras
======================================================

Important note
===============
This is partially a playground for learning Go, with the side benefit of
producing something that is useful.

The usefulness part of it is to have some command-line tools to configure
various IP cameras.  Some cameras, such as the Reolink, require that you
login and receive a token before calling any API functions, and so can be
difficult to control with some open security camera software.


Cameras
=======
Camera support is planned for at least the following cameras:

    * Reolink
    * Amcrest
    * SV3C
    * Foscam
    * Wanview

As of now, the following cameras have at least some minimal support:

    * Reolink



Reolink
-------
This is really more of a proof-of-concept to talk to the camera, at this point.

This camera requires a login token to call other APIs.  To acquire a token:


    ./reolink -ip 192.168.1.50 -username admin -password mypassword

This will return a token that can be used for a period of time.


To perform a command:

    ./reolink -ip <ip> -username <username> -token <token> -cmd <cmd> -data <cmddata>

e.g.:

    ./reolink -ip 192.168.1.50 -username admin -token 7eaeac2a9af55fb -cmd set-daynight -data night
