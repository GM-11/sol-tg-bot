package utils

import "os"

const DEX_SCREENER_API_BASE = "https://api.dexscreener.com/"

const PORT = ":8080"

var SOLANA_FM_API = os.Getenv("SOLANA_FM_API")

var HELIUS_API_KEY = os.Getenv("HELIUS_API")

const SOLANA_FM_API_URL = "https://api.solana.fm/v0/accounts/"

var BIRD_EYE_KEY = os.Getenv("BIRD_EYE_KEY")

const BIRD_EYE_API_URL = "https://public-api.birdeye.so/"

const SOLANA_TRACKER_API_URL = "https://data.solanatracker.io"
