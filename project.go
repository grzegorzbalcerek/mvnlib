package mvnlib

import (
  "fmt"
  "encoding/xml"
)

type Property struct {
  XMLName xml.Name
  Value string `xml:",innerxml"`
}
type Properties struct {
  XMLName xml.Name
  Values []Property `xml:",any"`
}
type Project struct {
  GroupId string `xml:"groupId"`
  ArtifactId string `xml:"artifactId"`
  Version string `xml:"version"`
  ParentArtifact Artifact `xml:"parent"`
  Properties Properties `xml:"properties"`
  Dependencies []Artifact `xml:"dependencies>dependency"`
}

func parsePom(options Options, b []byte) Project {
  var project Project
  xml.Unmarshal(b, &project)
  if options.Pom { fmt.Println(string(b)) }
  if project.GroupId == "" && project.ParentArtifact.GroupId != "" { project.GroupId = project.ParentArtifact.GroupId }
  if project.ArtifactId == "" && project.ParentArtifact.ArtifactId != "" { project.ArtifactId = project.ParentArtifact.ArtifactId }
  if project.Version == "" && project.ParentArtifact.Version != "" { project.Version = project.ParentArtifact.Version }
//   var parentProject Project
//   if project.ParentArtifact.GroupId != "" && project.ParentArtifact.ArtifactId != "" && project.ParentArtifact.Version != ""{
//     parentProject = project.ParentArtifact.resolvePom(options)
//   }
  for j, a := range project.Dependencies {
    if a.GroupId == "${project.groupId}" { project.Dependencies[j].GroupId = project.GroupId }
    if a.ArtifactId == "${project.artifactId}" { project.Dependencies[j].ArtifactId = project.ArtifactId }
    if a.Version == "${project.version}" { project.Dependencies[j].Version = project.Version }
    for _, p := range project.Properties.Values {
      if a.GroupId == "${"+p.XMLName.Local+"}" { project.Dependencies[j].GroupId = p.Value }
      if a.ArtifactId == "${"+p.XMLName.Local+"}" { project.Dependencies[j].ArtifactId = p.Value }
      if a.Version == "${"+p.XMLName.Local+"}" { project.Dependencies[j].Version = p.Value }
    }
  }
//   if options.Verbose { fmt.Println(project) }
  return project
}
