package prompt

import (
	"fmt"

	"github.com/Originate/git-town/src/cfmt"
	"github.com/Originate/git-town/src/git"
	"go.uber.org/dig"
)

type parentBranchService struct {
	branchService    IBranchService
	gitBranchService git.IBranchService
	gitConfigService git.IConfigService
	headerShown      bool
}

// ParentBranchServiceOpts are the opts to NewParentBranchService
type ParentBranchServiceOpts struct {
	dig.In

	BranchService    IBranchService
	GitBranchService git.IBranchService
	GitConfigService git.IConfigService
}

// NewParentBranchService returns a new IParentBranchService
func NewParentBranchService(opts ParentBranchServiceOpts) IParentBranchService {
	return &parentBranchService{
		branchService:    opts.BranchService,
		gitBranchService: opts.GitBranchService,
		gitConfigService: opts.GitConfigService,
	}
}

// EnsureKnowsParentBranches asserts that the entire ancestry for all given branches
// is known to Git Town.
// Missing ancestry information is queried from the user.
func (p *parentBranchService) EnsureKnowsParentBranches(branchNames []string) {
	for _, branchName := range branchNames {
		if p.gitConfigService.IsMainBranch(branchName) || p.gitConfigService.IsPerennialBranch(branchName) || p.gitConfigService.HasParentBranch(branchName) {
			continue
		}
		p.AskForBranchAncestry(branchName, p.gitConfigService.GetMainBranch())
		if parentBranchHeaderShown {
			fmt.Println()
		}
	}
}

// AskForBranchAncestry prompts the user for all unknown ancestors of the given branch
func (p *parentBranchService) AskForBranchAncestry(branchName, defaultBranchName string) {
	current := branchName
	for {
		parent := p.gitConfigService.GetParentBranch(current)
		if parent == "" {
			p.printParentBranchHeader()
			parent = p.AskForBranchParent(current, defaultBranchName)
			if parent == perennialBranchOption {
				p.gitConfigService.AddToPerennialBranches(current)
				break
			}
			p.gitConfigService.SetParentBranch(current, parent)
		}
		if parent == p.gitConfigService.GetMainBranch() || p.gitConfigService.IsPerennialBranch(parent) {
			break
		}
		current = parent
	}
}

// AskForBranchParent prompts the user for the parent of the given branch
func (p *parentBranchService) AskForBranchParent(branchName, defaultBranchName string) string {
	choices := p.gitBranchService.GetLocalBranchesWithMainBranchFirst()
	filteredChoices := p.filterOutSelfAndDescendants(branchName, choices)
	return p.branchService.askForBranch(askForBranchOptions{
		branchNames:       append([]string{perennialBranchOption}, filteredChoices...),
		prompt:            fmt.Sprintf(parentBranchPromptTemplate, branchName),
		defaultBranchName: defaultBranchName,
	})
}

// Helpers

func (p *parentBranchService) filterOutSelfAndDescendants(branchName string, choices []string) []string {
	result := []string{}
	for _, choice := range choices {
		if choice == branchName || p.gitConfigService.IsAncestorBranch(choice, branchName) {
			continue
		}
		result = append(result, choice)
	}
	return result
}

func (p *parentBranchService) printParentBranchHeader() {
	if !p.headerShown {
		p.headerShown = true
		cfmt.Printf(parentBranchHeaderTemplate, p.gitConfigService.GetMainBranch())
	}
}
