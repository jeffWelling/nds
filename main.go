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

func updateConfigs(configs map[string]nftConfigFile) map[string]nftConfigFile {
  return configs
}

func writeConfigs(configs map[string]nftConfigFile) {
  f_content := ""
  _ = f_content
  for filename, config := range configs {
		f_content = "table "
    _ = filename
    for i, set := range config.nftTables {
      _ = i
      _ = set
    }
  }
  log("Done writing config files")
}

func main() {
  app := cli.NewApp()
  app.Name = "dnstrigger"
  app.Usage = "Monitor DNS records for changes"
  app.Flags = []cli.Flag {
    cli.BoolFlag{ Name: "verbose" },
  }
  app.Action = func( c *cli.Context ) error {
    configs := readAllConfigs()
    configs =  updateConfigs(configs)
    writeConfigs(configs)
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

func readSetFile(config_file string) nftConfigFile {
  log(strings.Join([]string{"Loading ", config_file}, ": "))
  data, err := ioutil.ReadFile(config_file)
  readConfigErrCheck(err, config_file)
  data_string := strings.Fields( string(data) )

  in_set_decl   := false
  breakpoint    := 0
  nftAddrFamily := "<notset>"
  nftTableName  := "<notset>"
  nftSetName       := "<notset>"
  nftType       := "<notset>"
  nftElements   := make([]string, 0)
  nftSets       := make([]nftSet, 0)
  set           := nftSet{}
  table         := nftTable{}
  config        := nftConfigFile{nftConfigFileName: path.Base(config_file)}

  for i, _ := range data_string {
    if i < breakpoint {
      continue
    }
    if data_string[i] == "table" {
      nftAddrFamily = data_string[i + 1]
      nftTableName  = data_string[i + 2]
      breakpoint = i + 4
      table = nftTable{nftAddrFamily: nftAddrFamily, nftTableName: nftTableName}
      table.nftSets = make([]nftSet, 0)
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
      nftElements = make([]string, 0)
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
        set = nftSet{ nftAddrFamily: nftAddrFamily,
          nftTableName: nftTableName,
          nftSetName: nftSetName,
          nftType: nftType,
          nftFilename: path.Base(config_file),
          nftElements: nftElements }
        nftSets = append( nftSets, set )
        table.nftSets = append( table.nftSets, set )
				in_set_decl = false
        continue
      }
    }
    if data_string[i] == "}" {
      // End of 'table' declaration in file.
      config.nftTables = append( config.nftTables, table )
    }
  }
  return config
}

type nftConfigFile struct {
  nftConfigFileName string
  nftTables []nftTable
}

type nftTable struct {
  nftAddrFamily string
  nftTableName string
  nftSets []nftSet
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
func readAllConfigs() map[string]nftConfigFile {
  configs := make(map[string]nftConfigFile)
  loadable_files, _ := filepath.Glob("/etc/nft.conf.d/sets.d/domains.d/*.conf")
  for  _, config_file := range loadable_files {
    config := readSetFile(config_file)
    configs[config.nftConfigFileName] = config
  }
  fmt.Println(configs)
  return configs
}
