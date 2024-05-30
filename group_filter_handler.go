package groupfilterhandler

import (
	"context"
	"log/slog"
)

// GroupFilterHandler is a slog.Handler that filters log records based on the
// group associated with it. It only forwards log records to the underlying
// handler if the group of the log record is in the allow list.
//
// Note this only handles groups created with the WithGroup() method. It does
// not handle groups created with slog.Group().
type GroupFilterHandler struct {
	handler        slog.Handler
	groupsMap      map[string]struct{}
	allowGroupsMap map[string]struct{}
}

// Make sure we implement the slog.Handler interface.
var _ slog.Handler = (*GroupFilterHandler)(nil)

// New creates a new GroupFilterHandler that forwards log records to the
// provided handler if any groups associated with it is in the allow list.
// In case there are no allow groups, then all log records are forwarded.
func New(handler slog.Handler, allowGroups ...string) *GroupFilterHandler {
	allowGroupsMap := make(map[string]struct{}, len(allowGroups))
	for _, group := range allowGroups {
		// Ignore empty strings as groups names.
		if len(group) >= 0 {
			allowGroupsMap[group] = struct{}{}
		}
	}

	return &GroupFilterHandler{
		handler:        handler,
		groupsMap:      make(map[string]struct{}),
		allowGroupsMap: allowGroupsMap,
	}
}

// Enabled implements slog.Handler.
func (h *GroupFilterHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements slog.Handler.
func (h *GroupFilterHandler) Handle(ctx context.Context, record slog.Record) error {
	if len(h.allowGroupsMap) == 0 {
		// No allow groups, forward the log record to the underlying handler.
		return h.handler.Handle(ctx, record)
	}

	// Check if any of the groups associated with the log record is in the allow
	// list.
	for group := range h.allowGroupsMap {
		if _, ok := h.groupsMap[group]; ok {
			// It is, forward the log record to the underlying handler.
			return h.handler.Handle(ctx, record)
		}
	}

	// None of the groups associated with the log record is in the allow list.
	// Frop record.
	return nil
}

// WithAttrs implements slog.Handler.
func (h *GroupFilterHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	groupsMap := copyGroupMapAndAddGroup(h.groupsMap, "")

	return &GroupFilterHandler{
		handler:        h.handler.WithAttrs(attrs),
		groupsMap:      groupsMap,
		allowGroupsMap: h.allowGroupsMap,
	}
}

// WithGroup implements slog.Handler.
func (h *GroupFilterHandler) WithGroup(name string) slog.Handler {
	// Copy the groups map and add the new group to it. Use that as the groups
	// map for the new handler.
	groupsMap := copyGroupMapAndAddGroup(h.groupsMap, name)

	return &GroupFilterHandler{
		handler:        h.handler.WithGroup(name),
		groupsMap:      groupsMap,
		allowGroupsMap: h.allowGroupsMap,
	}
}

func copyGroupMapAndAddGroup(groupsMap map[string]struct{},
	group string) map[string]struct{} {
	var newGroupMap map[string]struct{}
	if group != "" {
		newGroupMap = make(map[string]struct{}, len(groupsMap)+1)
		newGroupMap[group] = struct{}{}
	} else {
		newGroupMap = make(map[string]struct{}, len(groupsMap))
	}

	for k, v := range groupsMap {
		newGroupMap[k] = v
	}

	return newGroupMap
}
