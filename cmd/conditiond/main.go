package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tadasv/conditiond"
	"os"
)

var (
	configFilePath = flag.String("config", "config.json", "path to configuration file")
	configDump     = flag.Bool("config-dump", false, "print current configuration and exit")
	cli            = flag.Bool("cli", false, "use CLI instead of http server")
	cliIn          = flag.String("cliIn", "-", "path to input file or - (dash) for stdin")
	cliOut         = flag.String("cliOut", "-", "path to result file or - (dash) for stdout")
	listenAddress  = flag.String("listenAddress", ":9000", "address to listen to for incoming http connections")
)

type ConditionMessage struct {
	Condition json.RawMessage `json:"condition"`
	Context   json.RawMessage `json:"context"`
}

type EvaluationResult struct {
	Error  *string     `json:"error"`
	Result interface{} `json:"result"`
}

func evaluatorFromConfig(cfg *Config) *condition.Evaluator {
	// registry -> name

	handlerMap := map[string]condition.ExpressionFunc{}
	for key, value := range cfg.EvaluatorConfig.FunctionMap {
		// It's ok to do this without checking for keys in the registry.  We're
		// assuming that the configuration was validated on start up and should
		// contain valid keys.
		handlerMap[key] = condition.ExpressionRegistry[value]
	}

	evaluator := condition.NewEvaluator()

	if len(cfg.EvaluatorConfig.FunctionWhitelist) > 0 {
		for _, whitelistedValue := range cfg.EvaluatorConfig.FunctionWhitelist {
			if handler, ok := handlerMap[whitelistedValue]; ok {
				evaluator.AddHandler(whitelistedValue, handler)
			}
		}
	} else {
		for key, f := range handlerMap {
			evaluator.AddHandler(key, f)
		}
	}

	return evaluator
}

func parseAndEvaluate(cfg *Config, msg *ConditionMessage) (interface{}, error) {
	e := condition.NewDefaultEvaluator()

	root, err := condition.Parse(string(msg.Condition))
	if err != nil {
		return nil, err
	}

	result, err := e.Evaluate(msg.Context, root)
	return result, err
}

func main() {
	flag.Parse()

	config, err := GetConfig(*configFilePath)
	if err != nil {
		panic(err)
	}

	if err := config.Validate(); err != nil {
		panic(err)
	}

	if *configDump {
		fmt.Fprintf(os.Stderr, "%s", config.String())
		return
	}

	if *cli {
		var dec *json.Decoder
		var enc *json.Encoder

		if *cliIn == "-" {
			dec = json.NewDecoder(os.Stdin)
		} else {
			fd, err := os.Open(*cliIn)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to open input file: %s\n", err.Error())
				os.Exit(-1)
			}
			dec = json.NewDecoder(fd)
			defer fd.Close()
		}

		if *cliOut == "-" {
			enc = json.NewEncoder(os.Stdout)
		} else {
			fd, err := os.OpenFile(*cliOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				fmt.Fprintf(os.Stderr, "unable to open output file %q: %s\n", *cliOut, err.Error())
				os.Exit(-1)
			}
			enc = json.NewEncoder(fd)
			defer fd.Sync()
			defer fd.Close()
		}

		for dec.More() {
			msg := ConditionMessage{}
			if err := dec.Decode(&msg); err != nil {
				panic(err)
			}

			value, err := parseAndEvaluate(config, &msg)
			resultMsg := EvaluationResult{
				Result: value,
			}
			if err != nil {
				errMsg := err.Error()
				resultMsg.Error = &errMsg
			}

			if err := enc.Encode(resultMsg); err != nil {
				panic(err)
			}
		}
	} else {
		panic("not implemented")
	}

}
