package templates

import (
	"fmt"
	"strings"

	"github.com/logrusorgru/aurora"
	"github.com/projectdiscovery/gologger"
	"github.com/projectdiscovery/nuclei/v3/pkg/model/types/severity"
	mapsutil "github.com/projectdiscovery/utils/maps"
)

var (
	Colorizer                       aurora.Aurora
	SeverityColorizer               func(severity.Severity) string
	deprecatedProtocolNameTemplates = mapsutil.SyncLockMap[string, bool]{Map: mapsutil.Map[string, bool]{}} //templates that still use deprecated protocol names
)

// TemplateLogMessage returns a beautified log string for a template
func TemplateLogMessage(id, name string, authors []string, templateSeverity severity.Severity) string {
	if Colorizer == nil || SeverityColorizer == nil {
		return ""
	}
	// Display the message for the template
	return fmt.Sprintf("[%s] %s (%s) [%s]",
		Colorizer.BrightBlue(id).String(),
		Colorizer.Bold(name).String(),
		Colorizer.BrightYellow(appendAtSignToAuthors(authors)).String(),
		SeverityColorizer(templateSeverity))
}

// appendAtSignToAuthors appends @ before each author and returns the final string
func appendAtSignToAuthors(authors []string) string {
	if len(authors) == 0 {
		return "@none"
	}

	values := make([]string, 0, len(authors))
	for _, k := range authors {
		if !strings.HasPrefix(k, "@") {
			values = append(values, fmt.Sprintf("@%s", k))
		} else {
			values = append(values, k)
		}
	}
	return strings.Join(values, ",")
}

// PrintDeprecatedProtocolNameMsgIfApplicable prints a message if deprecated protocol names are used
// Unless mode is silent we print a message for deprecated protocol name
func PrintDeprecatedProtocolNameMsgIfApplicable(isSilent bool, verbose bool) {
	count := 0
	_ = deprecatedProtocolNameTemplates.Iterate(func(k string, v bool) error {
		count++
		return nil
	})
	if count > 0 && !isSilent {
		gologger.Print().Msgf("[%v] 发现 %v 个模板加载了不推荐使用的协议语法\n", aurora.Yellow("WRN").String(), count)

		_ = deprecatedProtocolNameTemplates.Iterate(func(k string, v bool) error {
			gologger.Print().Msgf("  - %s\n", k)
			return nil
		})
	}

	deprecatedProtocolNameTemplates.Lock()
	deprecatedProtocolNameTemplates.Map = make(map[string]bool)
	deprecatedProtocolNameTemplates.Unlock()
}
