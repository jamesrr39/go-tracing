package tracingvizpkg

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/jamesrr39/go-tracing"
	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/goutil/streamtostorage"
)

//go:embed main.js
var mainJS string

//go:embed libs/htm/2.2.1/htm@2.2.1.js
var htmLib string

//go:embed libs/react/16/react.production.min.js
var reactLib string

//go:embed libs/react-dom/16/react-dom.production.min.js
var reactDomLib string

func streamToStorageReaderToRuns(reader *streamtostorage.Reader) ([]*tracing.Trace, errorsx.Error) {
	traces := []*tracing.Trace{}
	for {
		b, err := reader.ReadNextMessage()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, errorsx.Wrap(err)
		}

		trace := new(tracing.Trace)
		err = proto.Unmarshal(b, trace)
		if err != nil {
			// return nil, errorsx.Wrap(err)
			log.Printf("failed to unmarshal trace. Error: %q\n", err)
			continue
		}

		traces = append(traces, trace)
	}

	return traces, nil
}

func Generate(dataFilePath, outFilePath string) errorsx.Error {
	file, err := os.Open(dataFilePath)
	if err != nil {
		return errorsx.Wrap(err)
	}
	defer file.Close()

	reader := streamtostorage.NewReader(file, streamtostorage.MessageSizeBufferLenDefault)

	traces, err := streamToStorageReaderToRuns(reader)
	if err != nil {
		return errorsx.Wrap(err)
	}

	tracesJSONData, err := json.Marshal(traces)
	if err != nil {
		return errorsx.Wrap(err)
	}

	outFile, err := os.Create(outFilePath)
	if err != nil {
		return errorsx.Wrap(err)
	}
	defer outFile.Close()

	data := tplData{
		TracerDataJSON: template.JS(tracesJSONData),
		MainJS:         template.JS(mainJS),
		HtmLib:         template.JS(htmLib),
		ReactLib:       template.JS(reactLib),
		ReactDomLib:    template.JS(reactDomLib),
	}

	err = gotpl.Execute(outFile, data)
	if err != nil {
		return errorsx.Wrap(err)
	}

	return nil
}

type tplData struct {
	TracerDataJSON, MainJS, HtmLib, ReactLib, ReactDomLib template.JS
}

var gotpl = template.Must(template.New("profileviz").Parse(`
<html>
    <head>
    <meta charset="UTF-8">
    <title>Tracing</title>

    <style type="text/css">
        .events-table {
            width: 100%;
            background: lightblue;
        }
        .event-percentage-through-cell {
            width: 100%;
            border-left: 1px solid grey;
            border-right: 1px solid grey;
        }
        .event-name {
            min-width: 100px;
        }
        .event-since-start-of-run {
            min-width: 100px;
        }
		.traces-table {
			width: 100%;
		}
    </style>
    
    <script type="text/javascript">{{.HtmLib}}</script>
	<script type="text/javascript">{{.ReactLib}}</script>
	<script type="text/javascript">{{.ReactDomLib}}</script>
    
    <script type="module">
		window.traces = {{.TracerDataJSON}};
		{{.MainJS}}
    </script>
    </head>


    <body>
        <h1>Tracing</h1>
        <div id="App"></div>
    </body>
</html>
`))
