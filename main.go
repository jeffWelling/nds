package main

import (
  "strings"
  "fmt"
  "os"
  "path/filepath"
  "io/ioutil"
  "path"

  "github.com/urfave/cli"
)

func main() {
  app := cli.NewApp()
  app.Name = "dnstrigger"
  app.Usage = "Monitor DNS records for changes"
  app.Flags = []cli.Flag {
    cli.BoolFlag{ Name: "verbose" },
  }
  app.Action = func( c *cli.Context ) error {
    readAllConfigs()
    return nil
  }
  app.Run(os.Args)
}

func log(msg string) {
  fmt.Println(msg)
}

func readConfigErrCheck(e error, filename string) {
  if e != nil {
    log(strings.Join([]string{"Failed to load", filename}, ", "))
    panic(e)
  }
}

func readSetFile(config_file string) []nftSet {
  data, err := ioutil.ReadFile(config_file)
  readConfigErrCheck(err, config_file)
  data_string := strings.Fields( string(data) )

  in_set_decl   := false
  breakpoint    := 0
  nftAddrFamily := "<notset>"
  nftTableName  := "<notset>"
  nftSetName       := "<notset>"
  nftType       := "<notset>"
  nftElements   := make([]string, 1)
  nftSets       := make([]nftSet, 0)
  set           := nftSet{}

  for i, _ := range data_string {
    fmt.Printf("word: %s\n", data_string[i])
    fmt.Printf("breakpoint: %d\n", breakpoint)
    fmt.Printf("i: %d\n\n", i)
    if i < breakpoint {
      continue
    }
    if data_string[i] == "table" {
      nftAddrFamily = data_string[i + 1]
      nftTableName  = data_string[i + 2]
      breakpoint = i + 4
      continue
    }
    if data_string[i] == "set" {
      in_set_decl = true
      nftSetName = data_string[i + 1]
      breakpoint = i + 3
      continue
    }
    if data_string[i] == "type" {
      nftType = data_string[i + 1]
      breakpoint = i + 2
      continue
    }
    if data_string[i] == "elements" {
      i = i + 3
      for true {
        if data_string[i] == "}" {
          break
        }
        nftElements = append(nftElements, data_string[i])
        i = i + 1
      }
      breakpoint = i + 1
      continue
    }
    if in_set_decl == true {
			if data_string[i] == "}" {
        fmt.Printf("Adding a thingy...")
        set = nftSet{ nftAddrFamily: nftAddrFamily,
          nftTableName: nftTableName,
          nftSetName: nftSetName,
          nftType: nftType,
          nftFilename: path.Base(config_file),
          nftElements: nftElements }
        nftSets = append( nftSets, set )
				in_set_decl = false
        continue
      }
    }
    if data_string[i] == "}" {
      // End of 'table' declaration in file.
    }
  }
  return nftSets
}

func readConfig(config_file string) []nftSet {
  log(strings.Join([]string{"Loading ", config_file}, ": "))
  return readSetFile(config_file)
}

type nftSet struct {
  nftFilename string
  nftAddrFamily string
  nftTableName string
  nftSetName string
  nftType string
  nftElements []string
}

// Read `/etc/nft.conf.d/sets.s/*.conf`
// return config structure
func readAllConfigs() {
  log("Reading config")
  configs := make(map[string]nftSet)
  loadable_files, _ := filepath.Glob("/etc/nft.conf.d/sets.d/domains.d/*.conf")
  for  _, config_file := range loadable_files {
    file_sets := readConfig(config_file)
    for _,set := range file_sets {
      configs[set.nftFilename] = set
    }
  }
  log("Done reading config:")
  fmt.Println(configs)
}
