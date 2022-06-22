package mvnlib

type Options struct {
  GroupId string
  ArtifactId string
  Version string
  Verbose bool
  Parallel int
  Pom bool
  Recursive bool
  Repo string
  OutputDir string
}
