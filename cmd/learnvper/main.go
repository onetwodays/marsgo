package main

import (
    "bytes"
    "fmt"
    "github.com/spf13/viper"
    yaml "gopkg.in/yaml.v2"
    "os"
)

func main()  {
    //Establishing Defaults
    viper.SetDefault("ContentDir", "content")
    viper.SetDefault("LayoutDir", "layouts")
    viper.SetDefault("Taxonomies", map[string]string{"tag": "tags", "category": "categories"})

    //Reading Config Files
    viper.SetConfigName("core") // name of config file (without extension)
    viper.SetConfigType("yaml") // REQUIRED if the config file does not have the extension in the name
    viper.AddConfigPath("/etc/appname/")   // path to look for the config file in
    viper.AddConfigPath("$HOME/.appname")  // call multiple times to add many search paths
    viper.AddConfigPath(".")               // optionally look for config in the working directory
    err := viper.ReadInConfig() // Find and read the config file
    if err != nil { // Handle errors reading the config file
        if _,ok:= err.(viper.ConfigFileNotFoundError);ok {
            fmt.Println(err.Error(),"Config file not found") //Config File "core1" Not Found in "[/etc/appname /home/iliu/.appname /home/iliu/wish/marsgo]" Config file not found
        }else{
            fmt.Println(err.Error(),"Config file was found but another error was produced")
        }
        return
        //panic(fmt.Errorf("Fatal error config file: %s \n", err))
    }

    //As a rule of the thumb, everything marked with safe won't overwrite any file, but just create if not existent, whilst the default behavior is to create or truncate
    //viper.WriteConfig() // writes current config to predefined path set by 'viper.AddConfigPath()' and 'viper.SetConfigName'
    //viper.SafeWriteConfig()
    viper.WriteConfigAs("./test.yaml") // writes the current viper configuration to the given filepath. Will overwrite the given file, if it exists.
    //viper.SafeWriteConfigAs("./test.ini") // will error since it has already been written
    //viper.SafeWriteConfigAs("/path/to/my/.other_config")//writes the current viper configuration to the given filepath. Will not overwrite the given file, if it exists.

    //Watching and re-reading config files
    //Make sure you add all of the configPaths prior to calling WatchConfig()
    /*
    viper.WatchConfig()
    onChange:= func(e fsnotify.Event) {
        fmt.Println("Config file changed:", e.Name)
    }
    viper.OnConfigChange(onChange)
    */


    c := viper.AllSettings()
    bs, err := yaml.Marshal(c)
    if err != nil {
        //log.Fatalf("unable to marshal config to YAML: %v", err)
    }
    fmt.Println(string(bs))

    //
    viper.SetConfigType("yaml") // or viper.SetConfigType("YAML")

    // any approach to require this configuration into your program.
    var yamlExample = []byte(`
            Hacker: true
            name: steve
            hobbies:
            - skateboarding
            - snowboarding
            - go
            clothing:
              jacket: leather
              trousers: denim
            age: 35
            eyes : brown
            beard: true
    `)

    viper.ReadConfig(bytes.NewBuffer(yamlExample))

     // this would be "steve"
    fmt.Println(viper.Get("name"))


    //Setting Overrides
    //These could be from a command line flag, or from your own application logic.
    viper.Set("Verbose", true)



    viper.SetEnvPrefix("spf") // will be uppercased automatically
    viper.BindEnv("id")

    os.Setenv("SPF_ID", "13") // typically done outside of the app

    id := viper.Get("id") // 13
    fmt.Println("id=",id)
    /*
    Getting Values From Viper
    In Viper, there are a few ways to get a value depending on the value’s type. The following functions and methods exist:

    Get(key string) : interface{}
    GetBool(key string) : bool
    GetFloat64(key string) : float64
    GetInt(key string) : int
    GetIntSlice(key string) : []int
    GetString(key string) : string
    GetStringMap(key string) : map[string]interface{}
    GetStringMapString(key string) : map[string]string
    GetStringSlice(key string) : []string
    GetTime(key string) : time.Time
    GetDuration(key string) : time.Duration
    IsSet(key string) : bool
    AllSettings() : map[string]interface{}
    One important thing to recognize is that each Get function will return a zero value if it’s not found. To check if a given key exists, the IsSet() method has been provided.

    Example:

    viper.GetString("logfile") // case-insensitive Setting & Getting
    if viper.GetBool("verbose") {
        fmt.Println("verbose enabled")
    }

    Accessing nested keys
    The accessor methods also accept formatted paths to deeply nested keys. For example, if the following JSON file is loaded:

    {
        "host": {
            "address": "localhost",
            "port": 5799
        },
        "datastore": {
            "metric": {
                "host": "127.0.0.1",
                "port": 3099
            },
            "warehouse": {
                "host": "198.0.0.1",
                "port": 2112
            }
        }
    }
    Viper can access a nested field by passing a . delimited path of keys:

    GetString("datastore.metric.host") // (returns "127.0.0.1")
     */







}
