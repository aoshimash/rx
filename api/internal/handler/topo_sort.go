package handler

import (
	"github.com/aoshimash/rx/api/internal/domain"
	"github.com/google/uuid"
)

// topoSortGroups returns groups ordered so every parent appears before its children.
// Uses Kahn's algorithm. Returns a ValidationError if a cycle is detected.
func topoSortGroups(groups []domain.ProgramGroup) ([]domain.ProgramGroup, error) {
	byID := make(map[uuid.UUID]int, len(groups))
	for i, g := range groups {
		byID[g.ID] = i
	}

	inDegree := make([]int, len(groups))
	children := make([][]int, len(groups))
	for i, g := range groups {
		if g.ParentGroupID == nil {
			continue
		}
		parentIdx, ok := byID[*g.ParentGroupID]
		if !ok {
			continue
		}
		inDegree[i]++
		children[parentIdx] = append(children[parentIdx], i)
	}

	queue := make([]int, 0, len(groups))
	for i, deg := range inDegree {
		if deg == 0 {
			queue = append(queue, i)
		}
	}

	sorted := make([]domain.ProgramGroup, 0, len(groups))
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		sorted = append(sorted, groups[cur])
		for _, child := range children[cur] {
			inDegree[child]--
			if inDegree[child] == 0 {
				queue = append(queue, child)
			}
		}
	}

	if len(sorted) != len(groups) {
		return nil, &domain.ValidationError{
			Field:   "groups",
			Message: "circular parent_group_id reference detected",
		}
	}

	return sorted, nil
}
