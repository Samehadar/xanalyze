package main

import (
	"regexp"
)

// KeywordInfo contains information used by the front end to build the label
type KeywordInfo struct {
	Icon   string
	Color  string
	Regexp []string
}

// Keywords takes in any meta information strings and returns the keywordMapping key
// by testing the string against any regex in the KeywordMapping map
func Keywords(s string) (keywords []string) {
	// use the keyword map to deduplicate keywords for same line
	keywordMap := make(map[string]struct{})

	for key, info := range KeywordMapping {
		// compare the passed string against every regexp possibility
		for _, regex := range info.Regexp {
			// don't check again if the keyword has already been added
			if _, ok := keywordMap[key]; !ok {
				// add i flag to make case insenstive
				r := regexp.MustCompile(`(?i)` + regex)
				if r.Match([]byte(s)) {
					keywordMap[key] = struct{}{}
					// log.Printf("SubexpNames: %s\n", r.SubexpNames())
				}
			}
		}
	}

	// return the slice of keys for the keyword map
	for keyword := range keywordMap {
		keywords = append(keywords, keyword)
	}
	return keywords
}

var (
	Numeric              = `^(\d+)$`
	AlphaNumeric         = `^([0-9A-Za-z]+)$`
	Alpha                = `^([A-Za-z]+)$`
	AlphaCapsOnly        = `^([A-Z]+)$`
	AlphaNumericCapsOnly = `^([0-9A-Z]+)$`
	Url                  = `^((http?|https?|ftps?):\/\/)?([\da-z\.-]+)\.([a-z\.]{2,6})([\/\w \.-]*)*\/?$`
	Email                = `^(.+@([\da-z\.-]+)\.([a-z\.]{2,6}))$`
	HashtagHex           = `^#([a-f0-9]{6}|[a-f0-9]{3})$`
	ZeroXHex             = `^0x([a-f0-9]+|[A-F0-9]+)$`
	IPv4                 = `^(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)$`
	IPv6                 = `^([0-9A-Fa-f]{0,4}:){2,7}([0-9A-Fa-f]{1,4}$|((25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)(\.|$)){4})$`
)

var BrandsMapping = map[string]KeywordInfo{
	`google`:        {`google-icon`, `medium-orange`, []string{`google`, `angular`, `googlecloudplatform`, `googlechrome`, `golang`, `gwtproject`, `zxing`, `v8`}},
	`facebook`:      {`facebook-icon`, `medium-orange`, []string{`facebook`, `facebookarchive`, `boltsframework`}},
	`postgres`:      {`postgres-icon`, `medium-orange`, []string{`postgres`, `postgresql`}},
	`elasticsearch`: {`elastic-icon`, `medium-orange`, []string{`elastic`, `elasticsearch`}},
	`mongodb`:       {`mongo-icon`, `medium-orange`, []string{`mongodb`, `mongo`}},
	`zeromq`:        {`zmq-icon`, `medium-orange`, []string{`zeromq`, `zmq`, `0mq`}},
	`kubernetes`:    {`k8s-icon`, `medium-orange`, []string{`kubernetes`, `k8s`}},
	`boilerplate`:   {`boilerplate-icon`, `medium-orange`, []string{`boilerplate`, `seed`}},
	`phantom`:       {`phantom-icon`, `medium-orange`, []string{`phantom`, `phantomjs`}},
	`twitter`:       {`twitter-icon`, `medium-orange`, []string{`twbs`, `twitter`, `bower`, `flightjs`}},
	`microsoft`:     {`microsoft-icon`, `medium-orange`, []string{`microsoft`, `dotnet`, `aspnet`, `exceptionless`, `mono`, `winjs`}},
}

// KeywordMapping keywords, icons, colors, and extensions taken heavily from file-icons atom
// see https://github.com/file-icons/atom for their amazing work
var KeywordMapping = map[string]KeywordInfo{
	`Markdown`:       {`markdown-icon`, `medium-orange`, []string{`[a-zA-Z0-9-_]+\.md|[a-zA-Z0-9-_]+\.markdown`}},
	`ChangeLog`:      {`markdown-icon`, `medium-orange`, []string{`^CHANGELOG.md|^CHANGELOG|CHANGELOG`}},
	`Readme`:         {`markdown-icon`, `medium-orange`, []string{`^README|^README\.md|^readme\.md|^readme\.txt|^README\.txt|^README[a-zA-Z0-9-_]+\.md|^README[a-zA-Z0-9-_]+\.markdown`}},
	`Bower`:          {`bower-icon`, `medium-orange`, []string{`bower[-_]components`, `\.(bowerrc|bower\.json|Bowerfile)`}},
	`TravisCI`:       {`circleci-icon`, `dark-purple`, []string{`\.travis\.yml`, `\.travis\.yaml`}},
	`CircleCI`:       {`circleci-icon`, `dark-purple`, []string{`\.circleci`, `circle\.yml`}},
	`Ignore`:         {`ignore-icon`, `dark-blue`, []string{`\.dockerignore|\.gitignore|\.hgignore`}},
	`Docker`:         {`docker-icon`, `dark-blue`, []string{`dockerfile|dockerfile\s+|Dockerfile|Dockerfile\s+`}},
	`Docker-Compose`: {`docker-compose-icon`, `dark-blue`, []string{`docker-compose\.yml|docker-compose\.yaml|docker-compose[a-zA-Z0-9-_]+\.yml|docker-compose[a-zA-Z0-9-_]+\.yaml\s+`}},
	`Docker-Crane`:   {`k8s-icon`, `dark-blue`, []string{`crane\.yml|crane\.yaml|crane[a-zA-Z0-9-_]+\.yml|crane[a-zA-Z0-9-_]+\.yaml`}},
	`Kubernetes`:     {`k8s-icon`, `dark-blue`, []string{`k8s[a-zA-Z0-9-_]+\.yml`}},
	`Dropbox`:        {`dropbox-icon`, `medium-blue`, []string{`(Dropbox|\.dropbox\.cache)`}},
	`Emacs`:          {`emacs-icon`, `medium-purple`, []string{`\.emacs`}},
	`Framework`:      {`dylib-icon`, `medium-yellow`, []string{`\.framework`}},
	`Git`:            {`git-icon`, `medium-red`, []string{`\.git`}},
	`Github`:         {`fa-github`, `medium-orange`, []string{`\.github`}},
	`Gitlab`:         {`gitlab-icon`, `medium-orange`, []string{`\.gitlab`, `\.gitlab-ci\.yml`}},
	`Meteor`:         {`meteor-icon`, `dark-orange`, []string{`\.meteor|versions\.json`}},
	`Mercurial`:      {`hg-icon`, `dark-orange`, []string{`\.hg`, `mercurial`, `\.hgtags`, `\.hgignore`}},
	`Packet`:         {`package-icon`, `medium-green`, []string{`(bundle|paket)`}},
	`SVN`:            {`svn-icon`, `medium-yellow`, []string{`\.svn`, `subversion`}},
	`Textmate`:       {`textmate-icon`, `medium-green`, []string{`\.tmbundle`}},
	`Vagrant`:        {`vagrant-icon`, `medium-cyan`, []string{`\.vagrant`}},
	`Visual Studio`:  {`vs-icon`, `medium-blue`, []string{`\.vscode`}},
	`Xcode`:          {`appstore-icon`, `medium-cyan`, []string{`\.xcodeproj`}},
	`Debian`:         {`debian-icon`, `medium-red`, []string{`[a-zA-Z0-9]+\.deb\s+`}},
	`CMake`:          {`cmake-icon`, `medium-blue`, []string{`[a-zA-Z0-9]+\.cmake\s+`, `\.cmake `, `CMakeLists\.txt`, `cmakelists\.txt`, `cmakelists`}},
	`GoPkgs`:         {`go-icon`, `medium-blue`, []string{`^glide\.yaml`, `^glide\.lock`, `^Godeps`, `^Godeps/Godeps\.json`, `^vendor/manifest`, `^vendor/vendor\.json`, `^Gomfile`, `^Gopkg\.toml`, `^Gopkg\.lock`}},
	`NPM`:            {`npm-icon`, `medium-red`, []string{`package\.json|yarn\.lock|\.npmignore|\.?npmrc|npm-debug\.log|npm-shrinkwrap\.json|package-lock\.json`}},
	`NodeJS`:         {`node-icon`, `medium-green`, []string{`\.(njs|nvmrc|node|node-version)`, `node_modules`}},
	`Maven`:          {`maven-icon`, `medium-blue`, []string{`pom\.xml`, `ivy\.xml`, `build\.gradle`}},
	`RubyGems`:       {`rubygems-icon`, `medium-blue`, []string{`Gemfile`, `Gemfile\.lock`, `gems\.rb`, `gems\.locked`, `\.gemspec`}},
	`Packagist`:      {`packagist-icon`, `medium-blue`, []string{`composer\.json`, `composer\.lock`}},
	`PyPi`:           {`pypi-icon`, `medium-blue`, []string{`setup\.py`, `req[a-zA-Z0-9-_.]+txt`, `req[a-zA-Z0-9-_.]+pip`, `requirements/[a-zA-Z0-9-_.]+txt`, `requirements/[a-zA-Z0-9-_.]+pip`, `Pipfile`, `Pipfile\.lock`}},
	`Nuget`:          {`nuget-icon`, `medium-blue`, []string{`packages\.config`, `Project\.json`, `Project\.lock\.json`, `\.nuspec`, `paket\.lock`, `\.csproj`}},
	`CPAN`:           {`cpan-icon`, `medium-blue`, []string{`META\.json`, `META\.yml`}},
	`CocoaPods`:      {`cocoapods-icon`, `medium-blue`, []string{`Podfile`, `Podfile\.lock`, `\.podspec`}},
	`Clojars`:        {`clojars-icon`, `medium-blue`, []string{`project\.clj`}},
	`CRAN`:           {`cran-icon`, `medium-blue`, []string{`DESCRIPTION`}},
	`Cargo`:          {`cargo-icon`, `medium-blue`, []string{`Cargo\.toml`, `Cargo\.lock`}},
	`Hex`:            {`hex-icon`, `medium-blue`, []string{`mix\.exs`, `mix\.lock`}},
	`Swift`:          {`swift-icon`, `medium-blue`, []string{`Package\.swift`}},
	`Pub`:            {`pub-icon`, `medium-blue`, []string{`pubspec\.yaml`, `pubspec\.lock`}},
	`Carthage`:       {`carthage-icon`, `medium-blue`, []string{`Cartfile`, `Cartfile\.private`, `Cartfile\.resolved`}},
	`Dub`:            {`dub-icon`, `medium-blue`, []string{`dub\.json`, `dub\.sdl`}},
	`Julia`:          {`julia-icon`, `medium-blue`, []string{`REQUIRE`}},
	`Shards`:         {`shards-icon`, `medium-blue`, []string{`shard\.yml`, `shard\.lock`}},
	`Elm`:            {`elm-icon`, `medium-blue`, []string{`elm-package\.json`, `elm_dependencies\.json", "elm-stuff/exact-dependencies\.json`}},
	`PHPUnit`:        {`phpunit-icon`, `medium-purple`, []string{`\.phpunit\.xml `}},
	`NGINX`:          {`nginx-icon`, `dark-green`, []string{`nginx\.conf`}},
	`Jenkins`:        {`jenkins-icon`, `medium-red`, []string{`Jenkinsfile`}},
	`Gulp`:           {`gulp-icon`, `medium-red`, []string{`gulpfile\.js`, `gulpfile\.coffe`, `gulpfile\.babel\.js `}},
	`Ansible`:        {`ansible-icon`, `dark-blue`, []string{`\.(ansible|ansible\.yaml|ansible\.yml) `}},

	// `Meteor`:         {`meteor-icon`, `medium-blue`, []string{`versions\.json`}},
	// `Bower`:          {`bower-icon`, `medium-blue`, []string{`bower\.json`}},
	// `ArtText`:        {`arttext-icon`, `dark-purple`, []string{`\.artx`}},
	// `Atom`:           {`atom-icon`, `dark-green`, []string{`\.atom`}},
	// `Chef`:           {`chef-icon`, `dark-purple`, []string{`\.chef`}},
	// `Ruby`:           {`ruby-icon`, `medium-red`, []string{`\.(rb|ru|ruby|erb|gemspec|god|mspec|pluginspec|podspec|rabl|rake|opal|rails) `}},
	// `NPM`:           {`npm-icon`, `medium-blue`, []string{`package.json`, `package-lock.json`, `npm-shrinkwrap.json`, `yarn.lock`}},
	// `Zimpl`:         {`zimpl-icon`, `medium-orange`, []string{`\.(zimpl|zmpl|zpl) `}},
	// `Vue`:           {`vue-icon`, `light-green`, []string{`\.vue `}},
	// `Rust`:          {`rust-icon`, `medium-maroon`, []string{`\.(rs|\.rlib)`}},
	// `R`:             {`r-icon`, `medium-blue`, []string{`\.(r|Rprofile|rsx|rd) `}},
	// `Python`:        {`python-icon`, `dark-blue`, []string{`\.(py|\.ipy|pep|py3|\.pyi) `}},
	// `Perl`:          {`perl-icon`, `medium-blue`, []string{`\.(pl|perl|pm) `}},
	// `Perl6`:         {`perl6-icon`, `medium-purple`, []string{`\.(pl6|p6l|p6m) `}},
	// `PHP`:           {`php-icon`, `dark-blue`, []string{`\.php `}},
	// `Objective-C`:   {`objc-icon`, `dark-red`, []string{`\.objc`}},
	// `MATLAB`:        {`matlab-icon`, `medium-yellow`, []string{`\.matlab`}},
	// `Kotlin`:        {`kotlin-icon`, `dark-blue`, []string{`\.(kt|ktm|kts) `}},
	// `Java`:          {`java-icon`, `medium-purple`, []string{`\.java `}},
	// `Javascript`:    {`js-icon`, `medium-yellow`, []string{`\.(js|_js|jsb|jsm|jss|es6|es|mjs|sjs|ssjs|xsjs|dust) `}},
	// `HTML`: {`html5-icon`, `medium-orange`, []string{`\.html `}},
	// `Erlang`:        {`erlang-icon`, `medium-red`, []string{`\.erl `}},
	// `Dart`:          {`dart-icon`, `medium-cyan`, []string{`\.dart `}},
	// `CoffeeScript`:  {`coffee-icon`, `medium-maroon`, []string{`\.coffee `}},
	// `Clojure`:       {`clojure-icon`, `medium-cyan`, []string{`\.(clj|cl2|cljc|cljx|hic) `}},
	// `ClojureScript`: {`cljs-icon`, `dark-cyan`, []string{`\.cljs `}},
	// `C`:             {`c-icon`, `medium-blue`, []string{`\.(c|h) `}},
	// `C++`:           {`cpp-icon`, `light-blue`, []string{`\.(cpp|hpp) `}},
	// `C#`:            {`csharp-icon`, `darker-blue`, []string{`\.(csharp|cs) `}},
	// `Alpine Linux`:  {`alpine-icon`, `dark-blue`, []string{`\.APKBUILD `, ` apk `}},
}
