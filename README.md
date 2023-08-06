# Electrolux OCP

This repository contains the source code for communicating with the Electrolux OCP API (`api.ocp.electrolux.one`), see [`ocpapi`](./ocpapi).

The Electrolux OCP API (OnE Connected Platform?) is a REST API for controlling Electrolux connected appliances. It is used by e.g. the Electrolux app (previously Wellbeing, which used the Electrolux Delta API).

**Note:** This is not an official API and is not supported by Electrolux. Caveat emptor.

Only tested with a few Air Purifier appliances.

## Requirements

To use this library, you'll need an API key, client ID and client secret.

## Usage

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mafredri/electrolux-ocp/ocpapi"
)

func main() {
	client, err := ocpapi.New(ocpapi.Config{
		APIKey:       "[APIKey]",
		Brand:        "electrolux",
		ClientID:     "[ClientID]",
		ClientSecret: "[ClientSecret]",
		CountryCode:  "FI",
	})
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	err = client.Login(ctx, "[email]", "[password]")
	if err != nil {
		panic(err)
	}

	appliances, err := client.Appliances(ctx, true)
	if err != nil {
		panic(err)
	}

	var applianceIDs []string
	for _, a := range appliances {
		applianceIDs = append(applianceIDs, a.ApplianceID)
		fmt.Printf("%#v\n", a)
	}

	applianceInfo, err := client.AppliancesInfo(ctx, applianceIDs...)
	if err != nil {
		panic(err)
	}

	for _, a := range applianceInfo {
		fmt.Printf("%#v\n", a)
	}
}
```