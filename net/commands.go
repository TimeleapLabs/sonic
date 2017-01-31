/*
Copyright 2017 Mosaic Networks Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package net

import "github.com/arrivets/go-swirlds/hashgraph"

type KnownRequest struct {
	From string
}

type KnownResponse struct {
	Known map[string]int
}

type SyncRequest struct {
	Head   []byte
	Events []hashgraph.Event
}

type SyncResponse struct {
	// We may not succeed if we have a conflicting entry
	Success bool
}