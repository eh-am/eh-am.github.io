+++ 
date = 2024-03-02T14:59:42Z
title = "Fixing Remarkable's 'Screen Share' Screen in Yabai"
categories = ["technical"]
lang = "en"
+++

Pretty niche post. If you are running yabai and using the Screen Share feature from remarkable, yabai won't manage the window. [It looks like it doesn't manage Dialogs by default](https://github.com/koekeishiya/yabai/issues/2046), which that screen is (see the `"subrole": "AXDialog"` part):


```bash
❮ sleep 3 && yabai -m query --windows --window
```
```json
{
        "id":35147,
        "pid":50138,
        "app":"reMarkable",
        "title":"Screen Share",
        "frame":{
                "x":1020.0000,
                "y":85.0000,
                "w":909.0000,
                "h":778.0000
        },
        "role":"AXWindow",
        "subrole":"AXDialog",
        "root-window":true,
        "display":1,
        "space":2,
        "level":0,
        "sub-level":0,
        "layer":"normal",
        "sub-layer":"normal",
        "opacity":1.0000,
        "split-type":"none",
        "split-child":"none",
        "stack-index":0,
        "can-move":true,
        "can-resize":true,
        "has-focus":true,
        "has-shadow":true,
        "has-parent-zoom":false,
        "has-fullscreen-zoom":false,
        "is-native-fullscreen":false,
        "is-visible":true,
        "is-minimized":false,
        "is-hidden":false,
        "is-floating":true,
        "is-sticky":false,
        "is-grabbed":false
}
```

As per the docs above:
> Even if a window is eligible for management, only standard windows (subset of a real window -- e.g not dialogs) can be tiled unless a rule specifies the AXRole and AXSubrole filter with the manage=on property.

So we just need to add a rule to manage it automatically.

Add to your `.yabairc`:

```bash
yabai -m rule --add app="reMarkable" role="AXWindow" subrole="AXDialog" manage=on
```

And refresh by executing the config file (weird when you say it out loud huh?).
