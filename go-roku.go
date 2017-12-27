package main

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	ssdp "github.com/bcurren/go-ssdp"
	flag "github.com/ogier/pflag"
)

const usage = `Usage: go-roku [--channel=<channel>|--pause|--play|--on|--off]`

type rokuApp struct {
	Name string `xml:",chardata"`
	ID   string `xml:"id,attr"`
}

type rokuApps struct {
	XMLName xml.Name  `xml:"apps"`
	Apps    []rokuApp `xml:"app"`
}

func main() {
	os.Exit(mainWithReturn())
}

func mainWithReturn() int {

	res := 0

	var help *bool = flag.Bool("help", false, "Print usage.")
	var pause *bool = flag.Bool("pause", false, "Pause button.")
	var play *bool = flag.Bool("play", false, "Play button.")
	var powerOff *bool = flag.Bool("off", false, "Turn off TV.")
	var powerOn *bool = flag.Bool("on", false, "Turn on TV.")
	var channel *string = flag.String("channel", "", "Select channel.")

	flag.Parse()

	if *help {
		fmt.Println(usage)
	} else {

		done := false
		for _, tm := range []time.Duration{1000, 3000, 5000} {
			devs, _ := ssdp.Search("roku:ecp", tm*time.Millisecond)
			if len(devs) > 0 {

				done = true
				dev := devs[0]

				if *pause || *play {
					http.PostForm(fmt.Sprintf("%vkeypress/%v", dev.Location, "Play"), nil)
				} else if *powerOff {
					http.PostForm(fmt.Sprintf("%vkeypress/%v", dev.Location, "PowerOff"), nil)
				} else if *powerOn {
					http.PostForm(fmt.Sprintf("%vkeypress/%v", dev.Location, "PowerOn"), nil)
				} else {
					resp, err := http.Get(fmt.Sprintf("%squery/apps", dev.Location))
					if err != nil {
						fmt.Println(err)
						res = 1
					} else {

						defer resp.Body.Close()

						var dict rokuApps
						if err := xml.NewDecoder(resp.Body).Decode(&dict); err != nil {
							fmt.Println(err)
							res = 1
						} else {
							for _, app := range dict.Apps {
								if strings.HasPrefix(app.Name, *channel) {
									fmt.Println("Launching channel %s...\n")
									http.PostForm(fmt.Sprintf("%vlaunch/%v", dev.Location, app.ID), nil)
									break
								}
							}
						}
					}
				}
			}
			if done {
				break
			}
		}
	}

	return res
}

//======================================================================
// Local Variables:
// indent-tabs-mode: nil
// tab-width: 4
// fill-column: 70
// End:
