package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tadasv/conditiond"
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
	whitelistMap := map[string]interface{}{}

	for _, wl := range cfg.EvaluatorConfig.FunctionWhitelist {
		whitelistMap[wl] = nil
	}

	handlerMap := map[string]condition.ExpressionFunc{}
	for newFuncName, registryFuncName := range cfg.EvaluatorConfig.FunctionMap {
		// if nothing is whitelisted, we're allowing all functions; otherwise, only the ones that were whitelisted.
		if _, ok := whitelistMap[registryFuncName]; ok || cfg.EvaluatorConfig.FunctionWhitelist == nil {
			// It's ok to do this without checking for keys in the registry.  We're
			// assuming that the configuration was validated on start up and should
			// contain valid keys.
			handlerMap[newFuncName] = condition.ExpressionRegistry[registryFuncName]
		}
	}

	evaluator := condition.NewEvaluator()
	for key, f := range handlerMap {
		evaluator.AddHandler(key, f)
	}

	return evaluator
}

func parseAndEvaluate(e *condition.Evaluator, msg *ConditionMessage) (interface{}, error) {
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
		log.Fatalf("unable to read configuration: %s", err.Error())
	}

	if err := config.Validate(); err != nil {
		log.Fatalf("invalid configuration: %s", err.Error())
	}

	if *configDump {
		fmt.Fprintf(os.Stdout, "%s", config.String())
		return
	}

	evaluator := evaluatorFromConfig(config)

	if *cli {
		var dec *json.Decoder
		var enc *json.Encoder

		if *cliIn == "-" {
			dec = json.NewDecoder(os.Stdin)
		} else {
			fd, err := os.Open(*cliIn)
			if err != nil {
				log.Fatalf("unable to open input file: %s\n", err.Error())
			}
			dec = json.NewDecoder(fd)
			defer fd.Close()
		}

		if *cliOut == "-" {
			enc = json.NewEncoder(os.Stdout)
		} else {
			fd, err := os.OpenFile(*cliOut, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
			if err != nil {
				log.Fatalf("unable to open output file %q: %s\n", *cliOut, err.Error())
			}
			enc = json.NewEncoder(fd)
			defer fd.Sync()
			defer fd.Close()
		}

		for dec.More() {
			msg := ConditionMessage{}
			if err := dec.Decode(&msg); err != nil {
				log.Fatalf("unable to decode message: %s", err.Error())
			}

			value, err := parseAndEvaluate(evaluator, &msg)
			resultMsg := EvaluationResult{
				Result: value,
			}
			if err != nil {
				errMsg := err.Error()
				resultMsg.Error = &errMsg
			}

			if err := enc.Encode(resultMsg); err != nil {
				log.Fatalf("unable to encode message: %s", err.Error())
			}
		}
	} else {
		// Healthcheck endpoint
		http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		})

		// Evaluation endpoint
		http.HandleFunc("/evaluate", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			defer r.Body.Close()

			dec := json.NewDecoder(r.Body)
			enc := json.NewEncoder(w)
			results := []EvaluationResult{}

			for dec.More() {
				msg := ConditionMessage{}
				if err := dec.Decode(&msg); err != nil {
					w.WriteHeader(http.StatusBadRequest)
					return
				}

				value, err := parseAndEvaluate(evaluator, &msg)
				resultMsg := EvaluationResult{
					Result: value,
				}
				if err != nil {
					errMsg := err.Error()
					resultMsg.Error = &errMsg
				}

				results = append(results, resultMsg)
			}

			for _, result := range results {
				if err := enc.Encode(result); err != nil {
					// TODO
				}
			}
		})

		log.Printf("starting conditiond server on %s", *listenAddress)
		http.ListenAndServe(*listenAddress, nil)
	}

}
