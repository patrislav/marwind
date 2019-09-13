# Marwind WM

Marwind is a simple X11 tiling window manager written in Go. It was inspired by the [i3 window manager](https://i3wm.org/) and the [acme editor](http://acme.cat-v.org/) for Plan 9 and aims to combine the good parts of both for the best experience.

**Important! The project is under active development and is *not* stable. Use at your own risk.**

## Goals

- Sane defaults. It should be possible to install the WM and be immediately productive without spending hours on configuration
- Keyboard-driven without sacrificing the mouse. Marwind is focused on the keyboard not unlike most tiling managers, however mouse also has its place. Common actions - such as moving, resizing, or closing windows - should be possible using either of the input methods
- Dynamically reconfigurable. Provide standard HTTP / gRPC endpoints for on-the-fly configuration, without the need to reload the entire WM. These endpoints will also serve as points of communication with external applications.
- Clean code and documentation

## Licence

Copyright 2019 Patryk Kalinowski. All rights reserved.
