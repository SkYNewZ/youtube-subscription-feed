# youtube-subscription-feed

Connect to your Google account, list all your subscriptions and generate
a [.opml](https://fr.wikipedia.org/wiki/Outline_Processor_Markup_Language) file to import in your favorite RSS reader.

## Run

1. Get a `client_secret.json` file by following https://developers.google.com/youtube/v3/getting-started#before-you-start
2. Build
3. Run !

## Usage

```
Usage of youtube-subscription-feed:
  -debug
    	More verbose logs
  -output string
    	Filename to output generated opml (default "youtube-subscription-feed.opml")
  -raw
    	Print output instead of writing to file
  -secret-file string
    	Client secret file to use (default "client_secret.json")

```