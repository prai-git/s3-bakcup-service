package main

import (
    "io/ioutil"
    "encoding/json"
    "os"
    "log"
    "os/exec"
    "strings"
    "time"
    "net/http"
    "bytes"

)

type Json struct{
    Path string

}

type Paths struct{
    Sources []Json
}

var source Paths
var s3bucket string = os.Getenv("S3_BUCKET")
var date string = time.Now().Local().Format("2006-01-02")
var s3prefix string = "s3://"+s3bucket+"/backup-"+date+"/"
var mountpoint string = "/data"

func main(){
   loadpaths("config.json")
   verifybucket()
   syncfiles()
   if os.Getenv("SLACK_TOKEN") != ""{
     slackpush()
   }

}



func loadpaths(file string){

  raw, err := ioutil.ReadFile(file)
  if err!=nil{
      log.Fatal("Error:",err)
  }


  err=json.Unmarshal(raw, &source)
  if err!=nil{
      log.Fatal("Error:",err)
  }

  for i := 0; i < len(source.Sources); i++ {
    _,err := os.Stat(mountpoint + source.Sources[i].Path)
    if err != nil{
      log.Fatal("Something went, wrong your path," + source.Sources[i].Path +  " was not found")
    }
    log.Println(source.Sources[i].Path+ " Path loaded!")
  }
}

func verifybucket(){
  if s3bucket != "" {
    err := exec.Command("aws", "s3", "ls",s3bucket).Run()
	   if err != nil {
		     log.Fatal("Sorry it seems your bucket does't exist or is wrong")
	   }else{
       log.Println("The bucket " + s3bucket + " is valid!")
     }
  }else{
    log.Println("Please specify a bucket")

  }

}

func syncfiles(){
  var destdir string
  var sourcedir string
  for i := 0; i < len(source.Sources); i++ {
    words:= strings.Split(source.Sources[i].Path, "/")
    destdir = s3prefix + words[len(words) - 1]
    sourcedir = mountpoint + source.Sources[i].Path

    if s3bucket != "" {
      log.Println("Syncing " + sourcedir + " to " + destdir)
      err := exec.Command("aws", "s3", "sync", sourcedir, destdir).Run()
       if err != nil {
           log.Fatal("Sorry, something went wrong", err)
       }else{
         log.Println("Done!")
       }
     }
    }
}

 
func slackpush(){
  url := os.Getenv("SLACK_TOKEN")
    var jsonStr = []byte(`{"text":"Backup to S3 finished!"}`)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil{
      log.Fatal("Meh, something went wrong :( ")
    }else{ log.Println("response Body:", string(body),"So... everithing is fine!")}
}
