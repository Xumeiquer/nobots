# NoBots
Caddy Server plugin to protect your website against web crawlers and bots

**This is not an official plugin yet so you have to install it manually, please have a look at the end of this readme.**

## Usage

The directive for the Caddyfile is really simple. First, you have to place the bomb path next to the `nobots` keyword, for example `bomb.gz` in the example below.

Then you can specify user agent either as strings or regular expresions. When using regular expresions you must add `regexp` keyword in fornt of the regex.

Caddyfile example:

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  "DuckDuckBot"
  regexp "^[Bb]ot"
  regexp "bingbot"
}
```

There is another keyword that is useful in case you want to allow crawlers and bots navigate through specific part of your website. The keyword is `public` and its values are regular expresions so you can use it as following:

```
nobots "bomb.gz" {
  "Googlebot/2.1 (+http://www.googlebot.com/bot.html)"
  public "^/public"
  public "^/[a-z]{,5}/public"
}
```

The above example will send the to all URI except those that match with `/public` and `[a-z]{,5}/public`.

NOTE: By default all URI


## How to create a bomb
The bomb is not provided within the plugin so you have to create one.

In Linux is really easy, you can use the following commands.

```
dd if=/dev/zero bs=1M count=1024 | gzip > 1G.gzip
dd if=/dev/zero bs=1M count=10240 | gzip > 10G.gzip
dd if=/dev/zero bs=1M count=1048576 | gzip > 1T.gzip
```

In terms to optimize the final bomb you may compress them several times

```
cat 10G.gzip | gzip > 10G.gzipx2
cat 1T.gzip | gzip | gzip | gzip > 1T.gzipx4
 ```
*NOTE*: The extension `.gzipx2` or `.gzipx4` is only to highlight how many times the file was compress.



## Manual installation
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
