package plugins

import (
	"fmt"
	"path/filepath"
	"plugin"

	"github.com/hatchify/scribe"
)

func newPlugin(dir, key string, update bool) (pp *Plugin, err error) {
	var p Plugin
	p.importKey = key
	key, p.alias = parseKey(key)
	p.update = update

	switch {
	case filepath.Ext(key) != "":
		if len(p.alias) == 0 {
			p.alias = getPluginKey(key)
		}

		p.filename = key

	case isGitReference(key):
		var repoName string
		if _, repoName, p.branch, err = getGitURLParts(key); err != nil {
			return
		}

		if len(p.alias) == 0 {
			p.alias = repoName
		}

		// Set gitURL
		p.gitURL = removeBranchHash(key)

		// Set filename
		p.filename = filepath.Join(dir, p.alias+".so")

	default:
		err = fmt.Errorf("plugin type not supported: %s", key)
		return
	}

	p.out = scribe.New(p.alias)
	pp = &p
	return
}

// Plugin represents a plugin entry
type Plugin struct {
	out *scribe.Scribe
	p   *plugin.Plugin

	// Original import key
	importKey string
	// Alias given to plugin (e.g. github.com/user/myplugin would be myplugin)
	alias string
	// The git URL for the plugin
	gitURL string
	// The filename of the plugin's .so file
	filename string
	// The target branch of the plugin
	branch string

	// Signals if the plugin was loaded with an active update state
	update bool
}

func (p *Plugin) retrieve() (err error) {
	if len(p.gitURL) == 0 {
		return
	}

	switch {
	case !doesPluginSourceExist(p.gitURL):
		p.out.Notification("Plugin source does not yet exist, downloading repository")
		// We don't have the source yet. Download it
		if err = goGet(p.gitURL); err != nil {
			err = fmt.Errorf("error performing go get: %v", err)
			return
		}

	default:
		// We have the source already. Perform git pull to make sure the branch is synced
		p.out.Notification("Pulling most recent version")
		if _, err = gitPull(p.gitURL); err != nil {
			err = fmt.Errorf("error performing git pull: %v", err)
			return
		}

		// Update all the plugin dependencies
		p.out.Notification("Updating plugin dependencies")
		if err = updatePluginDependencies(p.gitURL); err != nil {
			err = fmt.Errorf("error updating plugin dependencies: %v", err)
			return
		}
	}

	p.out.Success("Download complete")
	return
}

func (p *Plugin) checkout() (err error) {
	if len(p.gitURL) == 0 || len(p.branch) == 0 {
		return
	}

	p.out.Notification("Checking out " + p.branch)
	var status string
	if status, err = gitCheckout(p.gitURL, p.branch); err != nil {
		err = fmt.Errorf("error encountered while switching to \"%s\": %v", p.branch, err)
		return
	} else if len(status) != 0 {
		p.out.Successf("Switched to \"%s\" branch", p.branch)
	}

	// Ensure we're up to date with the given branch
	if status, err = gitPull(p.gitURL); len(status) == 0 || err != nil {
		return
	}

	p.out.Successf("%s", status)

	// Ensure we have all the current dependencies
	if err = updatePluginDependencies(p.gitURL); err != nil {
		return
	}

	p.out.Success("Dependencies downloaded")
	return
}

func (p *Plugin) build() (err error) {
	if err = goBuild(p.gitURL, p.filename); err != nil {
		return
	}

	p.out.Success("Build complete")
	return
}

func (p *Plugin) test() (err error) {
	if doesPluginExist(p.filename) && !p.update {
		return
	}

	var pass bool
	if pass, err = goTest(p.gitURL); err != nil {
		p.out.Error("Test failed :(")
		return fmt.Errorf("%s failed test", p.alias)
	}

	if pass {
		p.out.Success("Test passed!")
	} else {
		p.out.Warning("No test files")
	}

	return
}

func (p *Plugin) init() (err error) {
	p.p, err = plugin.Open(p.filename)
	return
}
