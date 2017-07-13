# NoBots
A web bots protection plugin for caddy

For manual installation you shoud apply the following patch

```
diff --git a/caddy/caddymain/run.go b/caddy/caddymain/run.go
index b889971..90907ac 100644
--- a/caddy/caddymain/run.go
+++ b/caddy/caddymain/run.go
@@ -21,6 +21,7 @@ import (

        "github.com/mholt/caddy/caddytls"
        // This is where other plugins get plugged in (imported)
+        _ "github.com/Xumeiquer/nobots"
 )

 func init() {
diff --git a/caddyhttp/httpserver/plugin.go b/caddyhttp/httpserver/plugin.go
index a12ff0e..3239bb0 100644
--- a/caddyhttp/httpserver/plugin.go
+++ b/caddyhttp/httpserver/plugin.go
@@ -497,6 +497,7 @@ var directives = []string{
        "grpc",      // github.com/pieterlouw/caddy-grpc
        "gopkg",     // github.com/zikes/gopkg
        "restic",    // github.com/restic/caddy
+       "nobots",
 }

 const (
~
~
```

Once the patch is applied you must compile Caddy by running `github.com/mholt/caddy/caddy/build.bash/` from this directory `github.com/mholt/caddy/caddy/`

## Configuration

The configuration is quite easy. First of all place the bomb path then set a list of strings that match with the UA you want to ban from quering your site. You can set a regexp as well.

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  "Bing Bot"
  "Yahoo Bot"
  regexp "^bot"
}
```

## How to create a bomb
In Linux is really easy, for example you can use the following command.

```
dd if=/dev/zero bs=1M count=10240 | gzip > 10G.gzip
```

