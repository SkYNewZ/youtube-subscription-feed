package main

import (
	"context"
	"embed"
	"flag"
	"fmt"
	"os"
	"text/template"

	log "github.com/sirupsen/logrus"
)

var (
	//go:embed opml.tmpl
	content        embed.FS
	secretFilePath = "client_secret.json"
	raw, debug     bool
	output         = "youtube-subscription-feed.opml"
)

const youtubePartSnippet string = "snippet"

func main() {
	flag.StringVar(&secretFilePath, "secret-file", secretFilePath, "Client secret file to use")
	flag.StringVar(&output, "output", output, "Filename to output generated opml")
	flag.BoolVar(&raw, "raw", false, "Print output instead of writing to file")
	flag.BoolVar(&debug, "debug", false, "More verbose logs")
	flag.Parse()

	if debug {
		log.SetLevel(log.TraceLevel)
	}

	ctx := context.Background()
	service, err := GetYouTubeClient(ctx)
	handleError(err)

	channels, err := listSubscriptionsChannelIDs(ctx, service)
	handleError(err)

	if raw {
		for _, channel := range channels {
			fmt.Printf("https://www.youtube.com/feeds/videos.xml?channel_id=%s\n", channel.Snippet.ResourceId.ChannelId)
		}
		return
	}

	f, err := os.Create(output)
	handleError(err)
	defer f.Close()

	t := template.Must(template.New("opml.tmpl").ParseFS(content, "opml.tmpl"))
	handleError(t.Execute(f, channels))
}

func handleError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
