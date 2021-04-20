package tracingviz

import (
	_ "embed"
	"encoding/json"
	"html/template"
	"io"
	"os"

	"github.com/golang/protobuf/proto"
	"github.com/jamesrr39/go-tracing"
	"github.com/jamesrr39/goutil/errorsx"
	"github.com/jamesrr39/goutil/streamtostorage"
)

//go:embed main.js
var mainJS string

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
			return nil, errorsx.Wrap(err)
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
		TracerDataJSON: string(tracesJSONData),
		MainJS:         mainJS,
	}

	err = gotpl.Execute(outFile, data)
	if err != nil {
		return errorsx.Wrap(err)
	}

	return nil
}

type tplData struct {
	TracerDataJSON, MainJS string
}

var gotpl = template.Must(template.New("profileviz").Parse(`
<html>
    <head>
    <meta charset="UTF-8">
    <title>React + htm Demo</title>

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
    </style>
    
    <script src="https://unpkg.com/htm@2.2.1" crossorigin></script>
    <!--<script src="https://unpkg.com/react@16/umd/react.production.min.js" crossorigin></script>-->
    <script src="https://unpkg.com/react@16/umd/react.development.js" crossorigin></script>
    <!--<script src="https://unpkg.com/react-dom@16/umd/react-dom.production.min.js" crossorigin></script>-->
    <script src="https://unpkg.com/react-dom@16/umd/react-dom.development.js" crossorigin></script>
    
    <script type="module">
		window.tracerData = {{.TracerDataJSON}}
		{{.MainJS}}
    </script>
    </head>


    <body>
        <h1> React + htm Demo</h1>
        <div id="App"></div>
    </body>
</html>
`))
