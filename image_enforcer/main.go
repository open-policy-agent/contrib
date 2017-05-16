package main

import (
	"flag"
	"log"
	"os"
	"strings"
	"time"

	clairpkg "github.com/open-policy-agent/contrib/image_enforcer/clair"
	"github.com/open-policy-agent/contrib/image_enforcer/docker"
	opapkg "github.com/open-policy-agent/contrib/image_enforcer/opa"
)

var repoName = flag.String("repository", "", "Docker repository names (comma-separated)")
var pollDelay = flag.Duration("poll-delay", time.Second*60, "Polling loop delay")
var registryURL = flag.String("registry-url", "https://registry-1.docker.io", "Docker Registry URL")
var clairURL = flag.String("clair-url", "http://localhost:6060/v1", "Clair URL")
var opaURL = flag.String("opa-url", "http://localhost:8181/v1", "OPA URL")
var dockerDataRoot = flag.String("opa-layer-root", "/data/docker/layers", "Root path of Docker layer data")
var clairDataRoot = flag.String("opa-clair-root", "/data/clair/layers", "Root path of Clair layer index data")

var exit = make(chan struct{})

func main() {

	flag.Parse()

	user := os.Getenv("DOCKER_USER")
	password := os.Getenv("DOCKER_PASSWORD")

	opa := opapkg.New(*opaURL)
	clair := clairpkg.New(*clairURL)

	for _, repoName := range strings.Split(*repoName, ",") {
		registry := docker.New(*registryURL, user, password)
		poller := newPoller(*pollDelay, repoName, opa, registry, clair, *dockerDataRoot, *clairDataRoot)
		go poller.Run()
	}

	<-exit
}

type poller struct {
	delay          time.Duration
	repoName       string
	opa            *opapkg.OPA
	dockerDataRoot string
	clairDataRoot  string
	registry       *docker.Registry
	clair          *clairpkg.Clair
}

func newPoller(delay time.Duration, repoName string, opa *opapkg.OPA, registry *docker.Registry, clair *clairpkg.Clair, dockerDataRoot, clairDataRoot string) *poller {
	return &poller{
		delay:          delay,
		repoName:       repoName,
		opa:            opa,
		registry:       registry,
		clair:          clair,
		dockerDataRoot: dockerDataRoot,
		clairDataRoot:  clairDataRoot,
	}
}

func (p *poller) Run() {
	for {
		p.poll()
		log.Printf("Waiting %v before next poll cycle", p.delay)
		time.Sleep(p.delay)
	}

}

func (p *poller) poll() {

	log.Printf("Querying tags for %v", p.repoName)
	tags, err := p.registry.Tags(p.repoName)
	if err != nil {
		log.Printf("Unexpected error listing tags for %v: %v", p.repoName, err)
		return
	}

	for _, tag := range tags {

		log.Printf("Querying for manifest for %v:%v", p.repoName, tag)
		manifest, err := p.registry.Manifest(p.repoName, tag)
		if err != nil {
			log.Printf("Unexpected error getting manifest for %v:%v:", p.repoName, tag)
			continue
		}

		err = p.opa.Push(p.dockerDataRoot+"/"+p.repoName+"/"+tag, manifest)
		if err != nil {
			log.Printf("Failed to push Docker data to OPA for %v:%v: %v", p.repoName, tag, err)
		}

		for _, layer := range manifest.Layers {

			log.Printf("Building index for %v:%v digest: %v (this may take some time)", p.repoName, tag, layer.Digest)

			props := clairpkg.IndexProps{
				Name:   layer.Digest,
				Path:   p.registry.Path(p.repoName, layer),
				Format: "Docker",
				Headers: map[string]string{
					"Authorization": p.registry.Authorization(),
				},
			}

			err := p.clair.Index(props)
			if err != nil {
				log.Printf("Unexpected error while indexing %v:%v (digest: %v): %v", p.repoName, tag, layer.Digest, err)
				continue
			}

			log.Printf("Querying index for %v:%v digest: %v", p.repoName, tag, layer.Digest)
			info, err := p.clair.Layer(layer.Digest)
			if err != nil {
				log.Printf("Unexpected error while retrieving vulnerabilities for %v:%v (digest: %v): %v", p.repoName, tag, layer.Digest, err)
				continue
			}

			err = p.opa.Push(p.clairDataRoot+"/"+layer.Digest, info)
			if err != nil {
				log.Printf("Failed to push Clair data to OPA for %v: %v", layer.Digest, err)
			}
		}

	}

}
