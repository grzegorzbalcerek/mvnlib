package mvnlib

import (
  "fmt"
  "io"
  "os"
  "net/http"
  "strings"
)

type Artifact struct {
  GroupId string `xml:"groupId"`
  ArtifactId string `xml:"artifactId"`
  Version string `xml:"version"`
  Optional string `xml:"optional"`
  Scope string `xml:"scope"`
}

func Download(options Options) {
  a := Artifact{GroupId: options.GroupId, ArtifactId: options.ArtifactId, Version: options.Version}
  err := os.MkdirAll(options.OutputDir, 0755)
  if err != nil { fmt.Println(err); os.Exit(1) }
  dependencies := a.resolve(options)
  var artifactsCh = make(chan Artifact)
  var ackCh = make(chan int)
  for j := 1; j <= options.Parallel; j++ {
    go downloadArtifactJar(options, j, artifactsCh, ackCh)
  }
  for d, _ := range dependencies { artifactsCh <- d }
  close(artifactsCh)
  for j := 1; j <= options.Parallel; j++ {
    select { case <-ackCh: }
  }
}

func List(options Options) {
  a := Artifact{GroupId: options.GroupId, ArtifactId: options.ArtifactId, Version: options.Version}
  dependencies := a.resolve(options)
  for a, _ := range dependencies {
    fmt.Println(a.jarUrl(options))
  }
}

func downloadArtifactJar(options Options, n int, artifactsCh chan Artifact, ackCh chan int) {
  for a := range artifactsCh {
    a.downloadJar(options, n)
  }
  ackCh <- n
}

func (a Artifact) resolve(options Options) map[Artifact]bool {
  return a.resolve1(options, make(map[Artifact]bool))
}

func (a Artifact) resolve1(options Options, resolved map[Artifact]bool) map[Artifact]bool {
  project := a.resolvePom(options)
  if options.Recursive {
      for _, a := range project.Dependencies {
        if !resolved[a] && a.Version != "" && a.Scope != "provided" && a.Scope != "test" {
          resolved = a.resolve1(options, resolved)
        }
      }
  }
  resolved[a] = true
  return resolved
}

func (a Artifact) resolvePom(options Options) Project {
  body := a.getPom(options)
  project := parsePom(options, body)
  return project
}

func (a Artifact) toUrl(options Options) string {
  return options.Repo + "/" + strings.ReplaceAll(a.GroupId, ".", "/") + "/" + a.ArtifactId + "/" + a.Version + "/"
}

func (a Artifact) pomUrl(options Options) string {
  return a.toUrl(options) + a.ArtifactId + "-" + a.Version + ".pom"
}

func (a Artifact) jarUrl(options Options) string {
  return a.toUrl(options) + a.jarFileName()
}

func (a Artifact) jarFileName() string {
  return a.ArtifactId + "-" + a.Version + ".jar"
}

func (a Artifact) getPom(options Options) []byte {
  url := a.pomUrl(options)
  if options.Verbose { fmt.Println("get", url) }
  resp, err := http.Get(url)
  if err != nil { fmt.Println(err); os.Exit(1) }
  if resp.StatusCode != 200 { fmt.Println(resp); os.Exit(1) }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body)
  if err != nil { fmt.Println(err); os.Exit(1) }
  return body
}

func (a Artifact) downloadJar(options Options, n int) {
  url := a.jarUrl(options)
  dest := options.OutputDir + "/" + a.jarFileName()
  fmt.Println(n, "download:", url, "->", dest)
  resp, err := http.Get(url)
  if err != nil { fmt.Println(err); os.Exit(1) }
  if resp.StatusCode != 200 { fmt.Println(resp); os.Exit(1) }
  defer resp.Body.Close()
  body, err := io.ReadAll(resp.Body)
  if err != nil { fmt.Println(err); os.Exit(1) }
  err = os.WriteFile(options.OutputDir + "/" + a.jarFileName(), body, 0644)
  if err != nil { fmt.Println(err); os.Exit(1) }
  if options.Verbose { fmt.Println(n, "saved", dest) }
}
