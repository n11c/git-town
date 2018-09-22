package git

import (
	"sort"
	"strconv"
	"strings"

	"github.com/Originate/git-town/src/util"
)

// PrintableService provides methods for printing Git Town configuration
type PrintableService struct {
	configService IConfigService
}

// NewPrintableService returns a new PrintableService
func NewPrintableService(configService IConfigService) *PrintableService {
	return &PrintableService{
		configService: configService,
	}
}

// GetPrintableMainBranch returns a user printable main branch
func (p *PrintableService) GetPrintableMainBranch() string {
	output := p.configService.GetMainBranch()
	if output == "" {
		return noneString
	}
	return output
}

// GetPrintablePerennialBranches returns a user printable list of perennial branches
func (p *PrintableService) GetPrintablePerennialBranches() string {
	output := strings.Join(p.configService.GetPerennialBranches(), "\n")
	if output == "" {
		return noneString
	}
	return output
}

// GetPrintablePerennialBranchTrees returns a user printable list of perennial branches trees
func (p *PrintableService) GetPrintablePerennialBranchTrees() string {
	trees := []string{}
	for _, perennialBranch := range p.configService.GetPerennialBranches() {
		trees = append(trees, p.GetPrintableBranchTree(perennialBranch))
	}
	if len(trees) == 0 {
		return noneString
	}
	return strings.Join(trees, "\n")
}

// GetPrintableNewBranchPushFlag returns a user printable new branch push flag
func (p *PrintableService) GetPrintableNewBranchPushFlag() string {
	return strconv.FormatBool(p.configService.ShouldNewBranchPush())
}

// GetPrintableBranchTree returns a user printable branch tree
func (p *PrintableService) GetPrintableBranchTree(branchName string) (result string) {
	result += branchName
	childBranches := p.configService.GetChildBranches(branchName)
	sort.Strings(childBranches)
	for _, childBranch := range childBranches {
		result += "\n" + util.Indent(p.GetPrintableBranchTree(childBranch), 1)
	}
	return
}

// GetPrintableOfflineFlag returns a user printable offline flag
func (p *PrintableService) GetPrintableOfflineFlag() string {
	return strconv.FormatBool(p.configService.IsOffline())
}
