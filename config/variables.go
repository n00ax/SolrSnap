package config

import "os"

const ProductName = "solr-snapshot-service"
const CommitterEmail = "noah.whiteis@bluestembrands.com"

var GitCommit string = "unavailable"

// ENV Pulled Here
var GitRemote = os.Getenv("GIT_REMOTE")
var GitUsername = os.Getenv("GIT_USERNAME")
var GitPassword = os.Getenv("GIT_PASSWORD")
var FusionBaseUrl = os.Getenv("FUSION_BASE_URL")
