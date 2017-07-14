# NoBots
Caddy Server plugin to protect your website against web crawlers and bots

This is not an official plugin yet so you have to install it manually.

## Installation
You need to compile CaddyServer manually by running `github.com/mholt/caddy/caddy/build.bash/` from this directory `github.com/mholt/caddy/caddy/` on you repository clone. But first of all you must cahnge two files in terms to plug this pluing into CaddyServer.

The following code is the patch you sould apply on the CaddyServer source code.

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


## Configuration

The directive for the Caddyfile is really simple. First place the bomb path. In the example below `bomb.gz`. Then set a list of strings that will match with the UA you want to ban. Moreover you can set a regexp as well by preceding the keyword `regexp`.

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  "Bing Bot"
  "Yahoo Bot"
  regexp "^bot"
}
```

## How to create a bomb
In Linux is really easy, you can use the following command.

```
dd if=/dev/zero bs=1M count=1024 | gzip > 1G.gzip
dd if=/dev/zero bs=1M count=10240 | gzip > 10G.gzip
dd if=/dev/zero bs=1M count=1048576 | gzip > 1T.gzip
```

In terms to optimize the final bomb you can compress them several times

```
cat 10G.gzip | gzip > 10G.gzipx2
cat 1T.gzip | gzip | gzip | gzip > 1T.gzipx4
 ```
*NOTE*: The extension `.gzipx2` or `.gzipx4` is only to highlight how many times the file was compress.
