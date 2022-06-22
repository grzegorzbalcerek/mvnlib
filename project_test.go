package mvnlib

import (
    "os"
    "path"
    "testing"
)

func isEq(t *testing.T, label string, got any, wanted any) {
    if got != wanted {
        t.Errorf(`%s == %v; want %v`, label, got, wanted)
    }
}

func TestParsePom1(t *testing.T) {
    t.Log("in TestParsePom")
    bytes, err := os.ReadFile(path.Join("testdata","pom1.xml"))
    if err != nil { t.Fatal(err) }
    options := Options{}
    project := parsePom(options, bytes)
    isEq(t, "project.GroupId", project.GroupId, "group1")
    isEq(t, "project.ArtifactId", project.ArtifactId, "artifact1")
}
