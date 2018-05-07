// Copyright 2018 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package runtimevar_test

import (
	"context"
	"log"

	"github.com/google/go-cloud/runtimevar"
	"github.com/google/go-cloud/runtimevar/filevar"
)

type DBConfig struct {
	Host     string
	Port     string
	Username string
	Name     string
}

func initConfig() (*runtimevar.Variable, func()) {
	// Construct a runtimevar.Variable object.
	cfg, err := filevar.NewVariable("/etc/myapp/db.json",
		runtimevar.NewDecoder(&DBConfig{}, runtimevar.JSONDecode))
	if err != nil {
		log.Fatal(err)
	}

	return cfg, func() {
		cfg.Close()
	}
}

func Example() {
	cfg, cleanup := initConfig()
	defer cleanup()

	ctx := context.Background()
	// Call Watch to retrieve initial value before proceeding.
	snap, err := cfg.Watch(ctx)
	if err != nil {
		log.Fatalf("Error in retrieving initial variable: %v", err)
	}
	log.Printf("Value: %+v", snap.Value.(*DBConfig))

	// Get a Context with cancel func to stop the Watch call.
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Have a separate goroutine that waits for changes.
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				snap, err := cfg.Watch(ctx)
				if err != nil {
					// Handle errors.
					log.Printf("Error: %v", err)
					continue
				}
				// Use updated configuration accordingly.
				log.Printf("Value: %+v", snap.Value.(*DBConfig))
			}
		}
	}()
}