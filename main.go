package main

import (
	"context"
	_ "embed"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
	"path/filepath"
	"sync"
	"time"
)

//go:embed pre.js
var preScript string

//go:embed post.js
var postScript string

var defaultUrl = "https://mp.weixin.qq.com/s/wa98LIcSmoRz-1LxAfWZ9w"

func main() {

	extPath, _ := filepath.Abs("SingleFile")

	u := launcher.New().
		// Must use abs path for an extension
		Set("load-extension", extPath).
		MustLaunch()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	b := rod.New().ControlURL(u).MustConnect()
	defer b.Close()

	save := b.MustPage().Context(ctx)
	save.MustNavigate("chrome-extension://cpneebmdjnifhgajfhmcjmdkoknohimd/src/extension/ui/pages/options.html").MustWaitLoad()
	// set auto close
	time.Sleep(time.Millisecond * 100)
	save.MustEval(`document.getElementById('autoCloseInput').click()`)

	save.MustNavigate("chrome-extension://cpneebmdjnifhgajfhmcjmdkoknohimd/src/extension/core/bg/background.html").MustWaitLoad()

	var wg sync.WaitGroup
	var pages sync.Map
	go b.EachEvent(func(e *proto.TargetTargetCreated) {
		switch e.TargetInfo.Type {
		case proto.TargetTargetInfoTypePage:
			{
				_, loaded := pages.LoadOrStore(e.TargetInfo.TargetID, e.TargetInfo)
				if !loaded {
					wg.Add(1)
				}
			}
		}
	})()
	go b.EachEvent(func(e *proto.TargetTargetDestroyed) {
		_, loaded := pages.Load(e.TargetID)
		if loaded {
			wg.Done()
		}
	})()

	// create a new page
	p := b.MustPage().Context(ctx)

	// inject js for some special usage
	p.MustEvalOnNewDocument(preScript)

	p.MustNavigate(defaultUrl).MustWaitLoad()

	//p.EnableDomain(proto.NetworkEnable{})
	//go p.EachEvent(func(e *proto.NetworkRequestWillBeSent) {
	//	fmt.Printf("NetworkRequestWillBeSent %+v\n", e)
	//})()
	//go p.EachEvent(func(e *proto.NetworkResponseReceived) {
	//	fmt.Printf("NetworkResponseReceived  %+v\n", e)
	//})()

	// in this example, it slowly scrolls to the bottom for those lazy loading images
	duration := p.MustEval(postScript)

	// the post script can define how long we wait
	time.Sleep(time.Millisecond * time.Duration(duration.Int()))

	// save page with singlefile
	// https://github.com/puppeteer/puppeteer/issues/2486#issuecomment-602116047
	save.MustEval(`chrome.tabs.query({ active: true }, tabs => {
	chrome.browserAction.onClicked.dispatch(tabs[0]);
})`)

	wg.Wait()

	// check the file path
}
