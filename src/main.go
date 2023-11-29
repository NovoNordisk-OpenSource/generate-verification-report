package main

import (
    "flag"
    
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func main() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
    debug := flag.Bool("debug", false, "sets log level to debug")

    flag.Parse()

    // Default level 
    zerolog.SetGlobalLevel(zerolog.InfoLevel)
    if *debug {
        zerolog.SetGlobalLevel(zerolog.DebugLevel)
    }

    log.Info().Msg("Running..")
}
