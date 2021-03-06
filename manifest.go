package main

import (
	"code.cloudfoundry.org/cli/plugin/models"
	"fmt"
	"github.com/bluemixgaragelondon/cf-blue-green-deploy/from-cf-codebase/manifest"
)

type ManifestReader func(manifest.Repository, string) *plugin_models.GetAppModel

type ManifestAppFinder struct {
	Repo         manifest.Repository
	ManifestPath string
	AppName      string
}

// TODO This function was interesting, and now is boring and should be eliminated?
func (f *ManifestAppFinder) RoutesFromManifest(defaultDomain string) []plugin_models.GetApp_RouteSummary {
	if appParams := f.AppParams(defaultDomain); appParams != nil {
		return appParams.Routes
	}
	return nil
}

func IsHostEmpty(app plugin_models.GetAppModel) bool {
	for _, route := range app.Routes {
		if route.Host != "" {
			return false
		}
	}
	return true
}

func (f *ManifestAppFinder) AppParams(defaultDomain string) *plugin_models.GetAppModel {
	var manifest *manifest.Manifest
	var err error
	if f.ManifestPath == "" {
		manifest, err = f.Repo.ReadManifest("./")
	} else {
		manifest, err = f.Repo.ReadManifest(f.ManifestPath)
	}

	if err != nil {
		fmt.Println(err)
		return nil
	}

	apps, err := manifest.Applications(defaultDomain)

	if err != nil {
		fmt.Println(err)
		return nil
	}

	for index, app := range apps {
		if IsHostEmpty(app) {
			continue
		}

		if app.Name != "" && app.Name != f.AppName {
			continue
		}

		return &apps[index]
	}
	return nil
}
