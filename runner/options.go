package runner

import (
	"math"
	"os"
	"regexp"

	"github.com/projectdiscovery/fileutil"
	"github.com/projectdiscovery/goconfig"
	"github.com/projectdiscovery/goflags"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/gologger/formatter"
	"github.com/projectdiscovery/gologger/levels"
	"github.com/projectdiscovery/httpx/common/customheader"
	"github.com/projectdiscovery/httpx/common/customlist"
	customport "github.com/projectdiscovery/httpx/common/customports"
	fileutilz "github.com/projectdiscovery/httpx/common/fileutil"
	"github.com/projectdiscovery/httpx/common/stringz"
)

const (
	maxFileNameLength = 255
	two               = 2
	DefaultResumeFile = "resume.cfg"
)

type scanOptions struct {
	Methods                   []string
	StoreResponseDirectory    string
	RequestURI                string
	RequestBody               string
	VHost                     bool
	OutputTitle               bool
	OutputStatusCode          bool
	OutputLocation            bool
	OutputContentLength       bool
	StoreResponse             bool
	OutputServerHeader        bool
	OutputWebSocket           bool
	OutputWithNoColor         bool
	OutputMethod              bool
	ResponseInStdout          bool
	ChainInStdout             bool
	TLSProbe                  bool
	CSPProbe                  bool
	VHostInput                bool
	OutputContentType         bool
	Unsafe                    bool
	Pipeline                  bool
	HTTP2Probe                bool
	OutputIP                  bool
	OutputCName               bool
	OutputCDN                 bool
	OutputResponseTime        bool
	PreferHTTPS               bool
	NoFallback                bool
	NoFallbackScheme          bool
	TechDetect                bool
	StoreChain                bool
	MaxResponseBodySizeToSave int
	MaxResponseBodySizeToRead int
	OutputExtractRegex        string
	extractRegex              *regexp.Regexp
	ExcludeCDN                bool
	HostMaxErrors             int
}

func (s *scanOptions) Clone() *scanOptions {
	return &scanOptions{
		Methods:                   s.Methods,
		StoreResponseDirectory:    s.StoreResponseDirectory,
		RequestURI:                s.RequestURI,
		RequestBody:               s.RequestBody,
		VHost:                     s.VHost,
		OutputTitle:               s.OutputTitle,
		OutputStatusCode:          s.OutputStatusCode,
		OutputLocation:            s.OutputLocation,
		OutputContentLength:       s.OutputContentLength,
		StoreResponse:             s.StoreResponse,
		OutputServerHeader:        s.OutputServerHeader,
		OutputWebSocket:           s.OutputWebSocket,
		OutputWithNoColor:         s.OutputWithNoColor,
		OutputMethod:              s.OutputMethod,
		ResponseInStdout:          s.ResponseInStdout,
		ChainInStdout:             s.ChainInStdout,
		TLSProbe:                  s.TLSProbe,
		CSPProbe:                  s.CSPProbe,
		OutputContentType:         s.OutputContentType,
		Unsafe:                    s.Unsafe,
		Pipeline:                  s.Pipeline,
		HTTP2Probe:                s.HTTP2Probe,
		OutputIP:                  s.OutputIP,
		OutputCName:               s.OutputCName,
		OutputCDN:                 s.OutputCDN,
		OutputResponseTime:        s.OutputResponseTime,
		PreferHTTPS:               s.PreferHTTPS,
		NoFallback:                s.NoFallback,
		NoFallbackScheme:          s.NoFallbackScheme,
		TechDetect:                s.TechDetect,
		StoreChain:                s.StoreChain,
		OutputExtractRegex:        s.OutputExtractRegex,
		MaxResponseBodySizeToSave: s.MaxResponseBodySizeToSave,
		MaxResponseBodySizeToRead: s.MaxResponseBodySizeToRead,
		HostMaxErrors:             s.HostMaxErrors,
	}
}

// Options contains configuration options for httpx.
type Options struct {
	CustomHeaders             customheader.CustomHeaders
	CustomPorts               customport.CustomPorts
	matchStatusCode           []int
	matchContentLength        []int
	filterStatusCode          []int
	filterContentLength       []int
	Output                    string
	StoreResponseDir          string
	HTTPProxy                 string
	SocksProxy                string
	InputFile                 string
	Methods                   string
	RequestURI                string
	RequestURIs               string
	requestURIs               []string
	OutputMatchStatusCode     string
	OutputMatchContentLength  string
	OutputFilterStatusCode    string
	OutputFilterContentLength string
	InputRawRequest           string
	rawRequest                string
	RequestBody               string
	OutputFilterString        string
	OutputMatchString         string
	OutputFilterRegex         string
	OutputMatchRegex          string
	Retries                   int
	Threads                   int
	Timeout                   int
	filterRegex               *regexp.Regexp
	matchRegex                *regexp.Regexp
	VHost                     bool
	VHostInput                bool
	Smuggling                 bool
	ExtractTitle              bool
	StatusCode                bool
	Location                  bool
	ContentLength             bool
	FollowRedirects           bool
	StoreResponse             bool
	JSONOutput                bool
	CSVOutput                 bool
	Silent                    bool
	Version                   bool
	Verbose                   bool
	NoColor                   bool
	OutputServerHeader        bool
	OutputWebSocket           bool
	responseInStdout          bool
	chainInStdout             bool
	FollowHostRedirects       bool
	MaxRedirects              int
	OutputMethod              bool
	TLSProbe                  bool
	CSPProbe                  bool
	OutputContentType         bool
	OutputIP                  bool
	OutputCName               bool
	Unsafe                    bool
	Debug                     bool
	Pipeline                  bool
	HTTP2Probe                bool
	OutputCDN                 bool
	OutputResponseTime        bool
	NoFallback                bool
	NoFallbackScheme          bool
	TechDetect                bool
	TLSGrab                   bool
	protocol                  string
	ShowStatistics            bool
	RandomAgent               bool
	StoreChain                bool
	Deny                      customlist.CustomList
	Allow                     customlist.CustomList
	MaxResponseBodySizeToSave int
	MaxResponseBodySizeToRead int
	OutputExtractRegex        string
	RateLimit                 int
	Probe                     bool
	Resume                    bool
	resumeCfg                 *ResumeCfg
	ExcludeCDN                bool
	HostMaxErrors             int
}

// ParseOptions parses the command line options for application
func ParseOptions() *Options {
	options := &Options{}

	flagSet := goflags.NewFlagSet()
	flagSet.SetDescription(`httpx is a fast and multi-purpose HTTP toolkit allow to run multiple probers using [retryablehttp](https://github.com/projectdiscovery/retryablehttp-go) library, it is designed to maintain the result reliability with increased threads.`)

	createGroup(flagSet, "input", "Target",
		flagSet.BoolVar(&options.VHostInput, "vhost-input", false, "Get a list of vhosts as input"),
		flagSet.Var(&options.CustomHeaders, "H", "Custom Header to send with request"),
		flagSet.Var(&options.CustomPorts, "ports", "Port ranges to scan (nmap syntax: eg 1,2-10,11)"),
		flagSet.StringVar(&options.HTTPProxy, "http-proxy", "", "HTTP Proxy, eg http://127.0.0.1:8080"),
		flagSet.StringVar(&options.InputFile, "l", "", "Input file containing list of hosts to process"),
		flagSet.StringVar(&options.Methods, "x", "", "Request Methods to use, use 'all' to probe all HTTP methods"),
		flagSet.StringVar(&options.RequestURI, "path", "", "Request path/file (example '/api')"),
		flagSet.StringVar(&options.RequestURIs, "paths", "", "Command separated paths or file containing one path per line (example '/api/v1,/apiv2')"),
		flagSet.StringVar(&options.RequestBody, "body", "", "Content to send in body with HTTP request"),
	)

	createGroup(flagSet, "template", "Template",
		flagSet.BoolVar(&options.TLSGrab, "tls-grab", false, "Perform TLS(SSL) data grabbing"),
		flagSet.BoolVar(&options.TechDetect, "tech-detect", false, "Perform wappalyzer based technology detection"),
		flagSet.IntVar(&options.Threads, "threads", 50, "Number of threads"),
		flagSet.IntVar(&options.Retries, "retries", 0, "Number of retries"),
		flagSet.IntVar(&options.Timeout, "timeout", 5, "Timeout in seconds"),
		flagSet.BoolVar(&options.VHost, "vhost", false, "Check for VHOSTs"),
		flagSet.BoolVar(&options.FollowRedirects, "follow-redirects", false, "Follow HTTP Redirects"),
		flagSet.BoolVar(&options.FollowHostRedirects, "follow-host-redirects", false, "Only Follow redirects on the same host"),
		flagSet.IntVar(&options.MaxRedirects, "max-redirects", 10, "Max number of redirects to follow per host"),
		flagSet.BoolVar(&options.TLSProbe, "tls-probe", false, "Send HTTP probes on the extracted TLS domains"),
		flagSet.BoolVar(&options.CSPProbe, "csp-probe", false, "Send HTTP probes on the extracted CSP domains"),
		flagSet.BoolVar(&options.Unsafe, "unsafe", false, "Send raw requests skipping golang normalization"),
		flagSet.BoolVar(&options.Pipeline, "pipeline", false, "HTTP1.1 Pipeline probe"),
		flagSet.BoolVar(&options.HTTP2Probe, "http2", false, "HTTP2 probe"),
		flagSet.BoolVar(&options.NoFallback, "no-fallback", false, "Probe both protocol (HTTPS and HTTP)"),
		flagSet.BoolVar(&options.NoFallbackScheme, "no-fallback-scheme", false, "Probe with input protocol scheme"),
		flagSet.BoolVar(&options.RandomAgent, "random-agent", true, "Use randomly selected HTTP User-Agent header value"),
		flagSet.Var(&options.Allow, "allow", "Allow list of IP/CIDR's to process (file or comma separated)"),
		flagSet.Var(&options.Deny, "deny", "Deny list of IP/CIDR's to process (file or comma separated)"),
		flagSet.IntVar(&options.MaxResponseBodySizeToSave, "response-size-to-save", math.MaxInt32, "Max response size to save in bytes (default - unlimited)"),
		flagSet.IntVar(&options.MaxResponseBodySizeToRead, "response-size-to-read", math.MaxInt32, "Max response size to read in bytes (default - unlimited)"),
		flagSet.BoolVar(&options.Resume, "resume", false, "Resume scan using resume.cfg"),
		flagSet.BoolVar(&options.ExcludeCDN, "exclude-cdn", false, "Skip full port scans for CDNs (only checks for 80,443)"),
		flagSet.IntVar(&options.HostMaxErrors, "max-host-error", 30, "Max error count per host before skipping remaining path/s"),
	)

	createGroup(flagSet, "filters", "Filtering",
		flagSet.StringVar(&options.OutputMatchStatusCode, "mc", "", "Match response with specific status code (-mc 200,302)"),
		flagSet.StringVar(&options.OutputMatchContentLength, "ml", "", "Match response with specific content length (-ml 102)"),
		flagSet.StringVar(&options.OutputFilterStatusCode, "fc", "", "Filter response with specific status code (-fc 403,401)"),
		flagSet.StringVar(&options.OutputFilterContentLength, "fl", "", "Filter response with specific content length (-fl 23)"),
		flagSet.StringVar(&options.OutputFilterString, "filter-string", "", "Filter response with specific string"),
		flagSet.StringVar(&options.OutputMatchString, "match-string", "", "Match response with specific string"),
		flagSet.StringVar(&options.OutputFilterRegex, "filter-regex", "", "Filter response with specific regex"),
		flagSet.StringVar(&options.OutputMatchRegex, "match-regex", "", "Match response with specific regex"),
		flagSet.StringVar(&options.OutputExtractRegex, "extract-regex", "", "Display response content with matched regex"),
	)

	createGroup(flagSet, "rate-limit", "Rate-limit",
		flagSet.IntVar(&options.RateLimit, "rate-limit", 150, "Maximum requests to send per second"),
	)

	createGroup(flagSet, "output", "Output",
		flagSet.StringVar(&options.Output, "o", "", "File to write output to (optional)"),
		flagSet.BoolVar(&options.StatusCode, "status-code", false, "Display HTTP response status code"),
		flagSet.BoolVar(&options.ExtractTitle, "title", false, "Display page title"),
		flagSet.BoolVar(&options.Location, "location", false, "Display location header"),
		flagSet.BoolVar(&options.ContentLength, "content-length", false, "Display HTTP response content length"),
		flagSet.BoolVar(&options.StoreResponse, "sr", false, "Store HTTP response to directory (default 'output')"),
		flagSet.StringVar(&options.StoreResponseDir, "srd", "output", "Custom directory to store HTTP responses"),
		flagSet.BoolVar(&options.JSONOutput, "json", false, "Display output in JSON format"),
		flagSet.BoolVar(&options.CSVOutput, "csv", false, "Display output in CSV format"),
		flagSet.BoolVar(&options.OutputMethod, "method", false, "Display request method"),
		flagSet.BoolVar(&options.Silent, "silent", false, "Silent mode"),
		flagSet.BoolVar(&options.Version, "version", false, "Show version of httpx"),
		flagSet.BoolVar(&options.Verbose, "verbose", false, "Verbose Mode"),
		flagSet.BoolVar(&options.NoColor, "no-color", false, "Disable colored output"),
		flagSet.BoolVar(&options.OutputServerHeader, "web-server", false, "Display server header"),
		flagSet.BoolVar(&options.OutputWebSocket, "websocket", false, "Display server using websocket"),
		flagSet.BoolVar(&options.responseInStdout, "response-in-json", false, "Show Raw HTTP response In Output (-json only) (deprecated)"),
		flagSet.BoolVar(&options.responseInStdout, "include-response", false, "Show Raw HTTP response In Output (-json only)"),
		flagSet.BoolVar(&options.chainInStdout, "include-chain", false, "Show Raw HTTP Chain In Output (-json only)"),
		flagSet.BoolVar(&options.OutputContentType, "content-type", false, "Display content-type header"),
		flagSet.BoolVar(&options.OutputIP, "ip", false, "Display Host IP"),
		flagSet.StringVar(&options.InputRawRequest, "request", "", "File containing raw request"),
		flagSet.BoolVar(&options.Debug, "debug", false, "Debug mode"),
		flagSet.BoolVar(&options.OutputCName, "cname", false, "Display Host cname"),
		flagSet.BoolVar(&options.OutputCDN, "cdn", false, "Display CDN"),
		flagSet.BoolVar(&options.OutputResponseTime, "response-time", false, "Display the response time"),
		flagSet.BoolVar(&options.ShowStatistics, "stats", false, "Enable statistic on keypress (terminal may become unresponsive till the end)"),
		flagSet.BoolVar(&options.StoreChain, "store-chain", false, "Save chain to file (default 'output')"),
		flagSet.BoolVar(&options.Probe, "probe", false, "Display probe status"),
	)
	flagSet.Parse()
	// Read the inputs and configure the logging
	options.configureOutput()

	err := options.configureResume()
	if err != nil {
		gologger.Fatal().Msgf("%s\n", err)
	}

	showBanner()

	if options.Version {
		gologger.Info().Msgf("Current Version: %s\n", Version)
		os.Exit(0)
	}

	options.validateOptions()

	return options
}

func (options *Options) validateOptions() {
	if options.InputFile != "" && !fileutilz.FileNameIsGlob(options.InputFile) && !fileutil.FileExists(options.InputFile) {
		gologger.Fatal().Msgf("File %s does not exist.\n", options.InputFile)
	}

	if options.InputRawRequest != "" && !fileutil.FileExists(options.InputRawRequest) {
		gologger.Fatal().Msgf("File %s does not exist.\n", options.InputRawRequest)
	}

	multiOutput := options.CSVOutput && options.JSONOutput
	if multiOutput {
		gologger.Fatal().Msg("Results can only be displayed in one format: 'JSON' or 'CSV'\n")
	}

	var err error
	if options.matchStatusCode, err = stringz.StringToSliceInt(options.OutputMatchStatusCode); err != nil {
		gologger.Fatal().Msgf("Invalid value for match status code option: %s\n", err)
	}
	if options.matchContentLength, err = stringz.StringToSliceInt(options.OutputMatchContentLength); err != nil {
		gologger.Fatal().Msgf("Invalid value for match content length option: %s\n", err)
	}
	if options.filterStatusCode, err = stringz.StringToSliceInt(options.OutputFilterStatusCode); err != nil {
		gologger.Fatal().Msgf("Invalid value for filter status code option: %s\n", err)
	}
	if options.filterContentLength, err = stringz.StringToSliceInt(options.OutputFilterContentLength); err != nil {
		gologger.Fatal().Msgf("Invalid value for filter content length option: %s\n", err)
	}
	if options.OutputFilterRegex != "" {
		if options.filterRegex, err = regexp.Compile(options.OutputFilterRegex); err != nil {
			gologger.Fatal().Msgf("Invalid value for regex filter option: %s\n", err)
		}
	}
	if options.OutputMatchRegex != "" {
		if options.matchRegex, err = regexp.Compile(options.OutputMatchRegex); err != nil {
			gologger.Fatal().Msgf("Invalid value for match regex option: %s\n", err)
		}
	}
}

// configureOutput configures the output on the screen
func (options *Options) configureOutput() {
	// If the user desires verbose output, show verbose output
	if options.Verbose {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelVerbose)
	}
	if options.Debug {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelDebug)
	}
	if options.NoColor {
		gologger.DefaultLogger.SetFormatter(formatter.NewCLI(true))
	}
	if options.Silent {
		gologger.DefaultLogger.SetMaxLevel(levels.LevelSilent)
	}
}

func (options *Options) configureResume() error {
	options.resumeCfg = &ResumeCfg{}
	if options.Resume && fileutil.FileExists(DefaultResumeFile) {
		return goconfig.Load(&options.resumeCfg, DefaultResumeFile)

	}
	return nil
}

// ShouldLoadResume resume file
func (options *Options) ShouldLoadResume() bool {
	return options.Resume && fileutil.FileExists(DefaultResumeFile)
}

// ShouldSaveResume file
func (options *Options) ShouldSaveResume() bool {
	return true
}

func createGroup(flagSet *goflags.FlagSet, groupName, description string, flags ...*goflags.FlagData) {
	flagSet.SetGroup(groupName, description)
	for _, currentFlag := range flags {
		currentFlag.Group(groupName)
	}
}
