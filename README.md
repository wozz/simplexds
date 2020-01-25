# simplexds

this repo is meant to be used as an XDS server for envoy in integration testing. the first use case is for testing EDS. the current implementation when used with pre v1.13.0 releases of envoy causes a crash on EDS updates after running for a while with hosts being added and removed over time.
