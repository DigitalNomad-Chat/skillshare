package config

import (
	"os"
	"path/filepath"
	"strings"
)

// DiscoveredAgentsFromPath scans OpenClaw's workspace/agents/ directory for
// sub-agents that contain a skills/ subdirectory.
//
// This function is a thin wrapper around discovery.DiscoverOpenClawAgents that
// delegates skill counting to avoid an import cycle (config → sync → config).
func DiscoveredAgentsFromPath(openclawSkillsPath string) ([]DiscoveredAgent, error) {
	baseDir := filepath.Dir(openclawSkillsPath)
	workspaceAgentsDir := filepath.Join(baseDir, "workspace", "agents")

	var agents []DiscoveredAgent

	agentDirs, err := os.ReadDir(workspaceAgentsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	for _, dir := range agentDirs {
		if !dir.IsDir() {
			continue
		}
		agentName := dir.Name()
		agentBaseDir := filepath.Join(workspaceAgentsDir, agentName)
		skillsPath := filepath.Join(agentBaseDir, "skills")

		skillsInfo, err := os.Stat(skillsPath)
		if err != nil || !skillsInfo.IsDir() {
			continue
		}

		// Count skill entries using dir listing (lightweight — no full discovery)
		skillCount := countSkillEntries(skillsPath)

		linkedCount := 0
		localCount := 0
		if entries, err := os.ReadDir(skillsPath); err == nil {
			for _, entry := range entries {
				if !entry.IsDir() {
					continue
				}
				skillPath := filepath.Join(skillsPath, entry.Name())
				if isSymlinkPointingTo(skillPath, baseDir) {
					linkedCount++
				} else {
					localCount++
				}
			}
		}

		agents = append(agents, DiscoveredAgent{
			Name:          agentName,
			AgentDir:      agentBaseDir,
			SkillsPath:    skillsPath,
			LinkedCount:   linkedCount,
			LocalCount:    localCount,
			ExpectedCount: skillCount,
		})
	}

	return agents, nil
}

// isSymlinkPointingTo reports whether path is a symlink (or junction) pointing to
// a directory under the given source path. This mirrors the logic in
// sync.CheckStatusMerge without importing the sync package.
func isSymlinkPointingTo(path, sourceDir string) bool {
	linkTarget, err := os.Readlink(path)
	if err != nil {
		return false
	}
	if !filepath.IsAbs(linkTarget) {
		linkTarget = filepath.Join(filepath.Dir(path), linkTarget)
	}
	linkTarget = filepath.Clean(linkTarget)
	sourceDir = filepath.Clean(sourceDir)
	return strings.HasPrefix(linkTarget, sourceDir+string(filepath.Separator)) || linkTarget == sourceDir
}

// countSkillEntries counts the number of skill directories (folders with SKILL.md)
// in a skills directory. This is a lightweight approximation used by the OpenClaw
// agent scanner to avoid importing the sync package (which would create an
// import cycle: config → sync → config).
func countSkillEntries(dir string) int {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return 0
	}
	count := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(dir, entry.Name(), "SKILL.md")); err == nil {
			count++
		}
	}
	return count
}

// OpenClawWorkspaceSkillsPath returns the path to OpenClaw's workspace-level skills directory.
// This is the "also_scans" path for openclaw target.
// Returns empty string if the openclaw skills path is not provided.
func OpenClawWorkspaceSkillsPath(openclawSkillsPath string) string {
	if openclawSkillsPath == "" {
		return ""
	}
	baseDir := filepath.Dir(openclawSkillsPath)
	path := filepath.Join(baseDir, "workspace", "skills")
	if info, err := os.Stat(path); err == nil && info.IsDir() {
		return path
	}
	return ""
}

// GetOpenClawBasePath extracts the OpenClaw root path from its skills directory path.
// ~/.openclaw/skills → ~/.openclaw
func GetOpenClawBasePath(skillsPath string) string {
	return filepath.Dir(skillsPath)
}

// BuildOpenClawBaseFromWorkspace extracts the OpenClaw base path from a workspace-level path.
// ~/.openclaw/workspace/skills → ~/.openclaw
func BuildOpenClawBaseFromWorkspace(workspacePath string) string {
	// workspace/skills → workspace → openclaw
	return filepath.Dir(filepath.Dir(workspacePath))
}

// IsOpenClawTarget reports whether the target name matches OpenClaw or its aliases.
func IsOpenClawTarget(name string) bool {
	lower := strings.ToLower(name)
	return lower == "openclaw" || lower == "open-claw" || lower == "open_claw"
}
