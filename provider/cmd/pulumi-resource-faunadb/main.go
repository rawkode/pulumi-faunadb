package main

import (
	"github.com/rawkode/pulumi-faunadb/provider/pkg/provider"
	"github.com/rawkode/pulumi-faunadb/provider/pkg/version"
)

var providerName = "faunadb"

func main() {
	provider.Serve(providerName, version.Version)
}
