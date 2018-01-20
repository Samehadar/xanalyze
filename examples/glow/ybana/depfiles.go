package main

import (
	"path/filepath"
	"regexp"
)

var regexps map[string]*regexp.Regexp

func Find(path string) *DepFile {
	fname := filepath.Base(path)
	for _, f := range DepFiles {
		for _, fn := range f.Filenames {
			if fname == fn {
				return &f
			}
		}
		for _, p := range f.Patterns {
			if r, ok := regexps[p]; ok && r.MatchString(fname) {
				return &f
			}
		}
	}
	return nil
}

func init() {
	regexps = make(map[string]*regexp.Regexp)
	for _, f := range DepFiles {
		for _, p := range f.Patterns {
			regexps[p] = regexp.MustCompile(p)
		}
	}
}

const (
	pkgTypeGem       = "gem"
	pkgTypePackagist = "packagist"
	pkgTypePypi      = "pypi"
	pkgTypeNpm       = "npm"
	pkgTypeBower     = "bower"
	pkgTypeMaven     = "maven"
)

var (
	pmMaven    = PackageManager{Name: "Maven", URL: "http://maven.apache.org"}
	pmBundler  = PackageManager{Name: "Bundler", URL: "http://bundler.io"}
	pmGem      = PackageManager{Name: "gem", URL: "http://guides.rubygems.org/command-reference/"}
	pmNpm      = PackageManager{Name: "npm", URL: "https://www.npmjs.com"}
	pmYarn     = PackageManager{Name: "Yarn", URL: "https://yarnpkg.com"}
	pmBower    = PackageManager{Name: "Bower", URL: "https://bower.io/"}
	pmComposer = PackageManager{Name: "Composer", URL: "http://guides.rubygems.org/command-reference/"}
	pmPip      = PackageManager{Name: "pip", URL: "https://pip.pypa.io/en/stable/"}
)

type DepFile struct {
	Name            string           `json:"name"`                // used as a unique key
	Filenames       []string         `json:"filenames,omitempty"` // supported filenames
	Patterns        []string         `json:"patterns,omitempty"`  // patterns of supported filenames
	URL             string           `json:"url,omitempty"`       // documentation URL if any
	PackageType     string           `json:"package_type"`        // one of npm, pypi, etc.
	PackageManagers []PackageManager `json:"package_managers"`    // package managers capable of processing the file
}

type PackageManager struct {
	Name string `json:"name"`          // human readable name
	URL  string `json:"url,omitempty"` // documentation URL
}

var DepFiles = []DepFile{

	// java
	{
		Name: "pom.xml",
		Filenames: []string{
			"pom.xml",
		},
		URL:         "http://maven.apache.org/pom.html",
		PackageType: pkgTypeMaven,
		PackageManagers: []PackageManager{
			pmMaven,
		},
	},

	{
		Name: "maven-dependencies.json",
		Filenames: []string{
			"maven-dependencies.json",
			"gemnasium-maven-plugin.json",
		},
		PackageType: pkgTypeMaven,
		PackageManagers: []PackageManager{
			pmMaven,
		},
	},

	// ruby
	{
		Name: "Gemfile",
		Filenames: []string{
			"Gemfile",
			"gems.rb",
		},
		URL:         "http://bundler.io/man/gemfile.5.html",
		PackageType: pkgTypeGem,
		PackageManagers: []PackageManager{
			pmBundler,
		},
	},

	{
		Name: "Gemfile.lock",
		Filenames: []string{
			"Gemfile.lock",
			"gems.locked",
		},
		PackageType: pkgTypeGem,
		PackageManagers: []PackageManager{
			pmBundler,
		},
	},

	{
		Name: "gemspec",
		Patterns: []string{
			`\.gemspec$`,
		},
		URL:         "http://guides.rubygems.org/specification-reference/",
		PackageType: pkgTypeGem,
		PackageManagers: []PackageManager{
			pmGem,
		},
	},

	// npm
	{
		Name: "package.json",
		Filenames: []string{
			"package.json",
		},
		URL:         "https://docs.npmjs.com/files/package.json",
		PackageType: pkgTypeNpm,
		PackageManagers: []PackageManager{
			pmNpm,
			pmYarn,
		},
	},

	{
		Name: "npm-shrinkwrap.json",
		Filenames: []string{
			"npm-shrinkwrap.json",
		},
		URL:         "https://docs.npmjs.com/files/shrinkwrap.json",
		PackageType: pkgTypeNpm,
		PackageManagers: []PackageManager{
			pmNpm,
		},
	},

	{
		Name: "package-lock.json",
		Filenames: []string{
			"package-lock.json",
		},
		URL:         "https://docs.npmjs.com/files/package-lock.json",
		PackageType: pkgTypeNpm,
		PackageManagers: []PackageManager{
			pmNpm,
		},
	},

	{
		Name: "yarn.lock",
		Filenames: []string{
			"yarn.lock",
		},
		URL:         "https://yarnpkg.com/lang/en/docs/yarn-lock/",
		PackageType: pkgTypeNpm,
		PackageManagers: []PackageManager{
			pmYarn,
		},
	},

	// bower
	{
		Name: "bower.json",
		Filenames: []string{
			"bower.json",
		},
		URL:         "https://github.com/bower/spec/blob/master/json.md#name",
		PackageType: pkgTypeBower,
		PackageManagers: []PackageManager{
			pmBower,
		},
	},

	// php
	{
		Name: "composer.json",
		Filenames: []string{
			"composer.json",
		},
		URL:         "https://getcomposer.org/doc/04-schema.md",
		PackageType: pkgTypePackagist,
		PackageManagers: []PackageManager{
			pmComposer,
		},
	},

	{
		Name: "composer.lock",
		Filenames: []string{
			"composer.lock",
		},
		PackageType: pkgTypePackagist,
		PackageManagers: []PackageManager{
			pmComposer,
		},
	},

	// python
	{
		Name: "requirements.txt",
		Filenames: []string{
			"requires.txt",
			"requirements.txt",
			"requirements.pip",
		},
		Patterns: []string{
			`requirements.*\.txt$`,
		},
		URL:         "https://pip.pypa.io/en/stable/reference/pip_install/#requirements-file-format",
		PackageType: pkgTypePypi,
		PackageManagers: []PackageManager{
			pmPip,
		},
	},

	// TODO: https://github.com/naiquevin/pipdeptree

}
