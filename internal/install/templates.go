package install

import "fmt"

type Target string

const (
	TargetCodex  Target = "codex"
	TargetClaude Target = "claude"
	TargetAll    Target = "all"
)

type SkillTemplate struct {
	Name        string
	Description string
	Content     string
}

func Templates() []SkillTemplate {
	return []SkillTemplate{
		{
			Name:        "diary-get",
			Description: "Retrieve prior Diary context before starting work.",
			Content:     diaryGetContent(),
		},
		{
			Name:        "diary-record",
			Description: "Compact the latest implementation work and record it with Diary.",
			Content:     diaryRecordContent(),
		},
		{
			Name:        "diary-list",
			Description: "List available Diary projects, records, ids, and hash prefixes.",
			Content:     diaryListContent(),
		},
	}
}

func ValidateTarget(target Target) error {
	switch target {
	case TargetCodex, TargetClaude:
		return nil
	default:
		return fmt.Errorf("unsupported target: %s", target)
	}
}

func ExpandTargets(target Target) ([]Target, error) {
	switch target {
	case TargetCodex, TargetClaude:
		return []Target{target}, nil
	case TargetAll:
		return []Target{TargetCodex, TargetClaude}, nil
	default:
		return nil, fmt.Errorf("unsupported target: %s", target)
	}
}

func diaryGetContent() string {
	return "---\n" +
		"name: diary-get\n" +
		"description: Retrieve prior Diary context before starting implementation work.\n" +
		"---\n\n" +
		"# Diary Get\n\n" +
		"Use this when starting or resuming work in a repository.\n\n" +
		"Run the default context retriever:\n\n" +
		"```bash\n" +
		"diary get\n" +
		"```\n\n" +
		"Use the returned context to understand the latest handoff, blockers, decisions, and next steps.\n\n" +
		"When the user asks for a specific stored record, retrieve it exactly:\n\n" +
		"```bash\n" +
		"diary get --id <id>\n" +
		"diary get --hash <prefix>\n" +
		"```\n\n" +
		"Do not read `.env` files while gathering context.\n"
}

func diaryRecordContent() string {
	return "---\n" +
		"name: diary-record\n" +
		"description: Compact the latest implementation work and record it with Diary.\n" +
		"---\n\n" +
		"# Diary Record\n\n" +
		"Use this when ending a meaningful implementation segment or before context may be lost.\n\n" +
		"## Gather Worktree Context\n\n" +
		"If Git is available, inspect changed files without modifying the worktree:\n\n" +
		"```bash\n" +
		"git status --short\n" +
		"git diff --stat\n" +
		"```\n\n" +
		"Do not read `.env` files or record secrets.\n\n" +
		"## Write A Compact Summary\n\n" +
		"Create a concise handoff with this structure:\n\n" +
		"```md\n" +
		"## Summary\n" +
		"Briefly describe what changed and why.\n\n" +
		"## Implementation\n" +
		"- Key implementation details and decisions.\n\n" +
		"## Files In Scope\n" +
		"- `path/to/file` - what changed and why.\n\n" +
		"## Files Out Of Scope\n" +
		"- `path/to/file` - dirty but unrelated or pre-existing.\n" +
		"- None, if no unrelated changes were observed.\n\n" +
		"## Verification\n" +
		"- Commands run and results.\n" +
		"- Commands not run and why.\n\n" +
		"## Blockers\n" +
		"- Known issues, missing information, or risks.\n\n" +
		"## Next Steps\n" +
		"- What the next run should do first.\n" +
		"```\n\n" +
		"Then record the summary:\n\n" +
		"```bash\n" +
		"diary record \"<compaction summary>\"\n" +
		"```\n\n" +
		"Include file references when useful:\n\n" +
		"```bash\n" +
		"diary record --file internal/storage/record.go \"<compaction summary>\"\n" +
		"```\n"
}

func diaryListContent() string {
	return "---\n" +
		"name: diary-list\n" +
		"description: List available Diary projects, records, ids, and hash prefixes.\n" +
		"---\n\n" +
		"# Diary List\n\n" +
		"Use this to discover stored Diary context without retrieving full prompt-ready context.\n\n" +
		"Run:\n\n" +
		"```bash\n" +
		"diary list\n" +
		"diary list --projects\n" +
		"```\n\n" +
		"Use listed ids or hash prefixes with `diary get --id <id>` or `diary get --hash <prefix>`.\n"
}
