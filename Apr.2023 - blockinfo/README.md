## Project: "blockinfo"

### Generated with
 - Types for network messaging: false
 - Enabled Cloud feature: false

### Supervision Tree

Applications
 - `BlockInfoApp{}` blockinfo/apps/blockinfoapp/blockinfoapp.go
   - `CommonSup{}` blockinfo/apps/blockinfoapp/commonsup.go
     - `Dispatcher{}` blockinfo/apps/blockinfoapp/dispatcher.go
     - `CrawlerSup{}` blockinfo/apps/blockinfoapp/crawlersup.go
	   - ... has more items. See source code

Process list that is starting by node directly
 - `Storage{}` blockinfo/cmd/storage.go
 - `Web{}` blockinfo/cmd/web.go


#### Used command
`ergo -init blockinfo -with-app BlockInfoApp -with-sup BlockInfoApp:CommonSup{type:rfo} -with-actor CommonSup:Dispatcher -with-sup CommonSup:CrawlerSup -with-actor CrawlerSup:CrawlerBlock -with-actor CrawlerSup:CrawlerAddress -with-actor Storage -with-web Web{ssl:yes,port:9090,handlers:3}`
