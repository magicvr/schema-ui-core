// Account avatar upload (W13 T-05 · GOAL-014).
//
// Self-service avatar surface on the personal-center Profile tab: any
// authenticated user may upload their avatar image. Files are processed by
// the shared RasterAssetStore (server re-encode, dimension-limited to 256px,
// never raw user bytes) and served publicly through the same pinned headers
// as brand assets (nosniff + sandbox + immutable). The avatar URL is then
// committed to the profile via PATCH /api/account/profile (avatarUrl).
//
// Cleanup model: uploading a new avatar deletes the user's previous avatar
// file (best-effort); clearing avatarUrl in the profile PATCH deletes it too.
// There is no startup GC for avatars (the referenced set lives in the users
// table); replace/clear covers the common paths and a crashed upload leaves
// at most one orphan file, which is acceptable for this auxiliary surface.
package handler

import (
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// AccountAvatarRoutes returns the account avatar surface: authenticated
// upload (self-service, no permission key) + public GET.
func AccountAvatarRoutes(a authMiddleware, store *RasterAssetStore, repository AccountRepository, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	h := &accountAvatarHandler{store: store, repository: repository, operations: operations}
	identity := func(method, pattern string) kernel.ContributionIdentity {
		return kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)}
	}
	return []kernel.RouteContribution{
		{ContributionIdentity: identity("POST", "/api/account/avatar"), Method: "POST", Pattern: "/api/account/avatar", Handler: a.Middleware(h.upload())},
		{ContributionIdentity: identity("GET", "/api/account/avatars/{id}"), Method: "GET", Pattern: "/api/account/avatars/{id}", Handler: store.file(), Public: true},
	}
}

type accountAvatarHandler struct {
	store      *RasterAssetStore
	repository AccountRepository
	operations operationlog.Recorder
}

// upload handles POST /api/account/avatar (multipart "file" part). On
// success the previous avatar file (if any and stored in this store) is
// deleted best-effort, then the new asset URL is returned.
func (h *accountAvatarHandler) upload() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.UserIdentityFrom(r.Context())
		if !ok {
			writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
			return
		}
		// W7 F-004: enforce a per-user avatar file cap before storing. A user
		// may upload the same total number across replace cycles; referenced
		// avatars still count until the profile switches away (startup GC then
		// reclaims unreferenced files).
		owned, err := h.store.CountOwner(user.ID)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "STORAGE_UNAVAILABLE", "could not check avatar quota")
			return
		}
		if owned >= maxAvatarPerUser {
			writeLocalizedError(w, r, http.StatusRequestEntityTooLarge, "AVATAR_QUOTA_EXCEEDED", "avatar upload quota exceeded; assign or clear your current avatar first")
			return
		}
		payload, ok := h.store.storeUploadForOwner(w, r, user.ID)
		if !ok {
			return
		}
		// Replace semantics: remove the user's previous avatar asset (if any)
		// after the new file is safely persisted. Best-effort — a leftover file
		// is an orphan, never a security issue (server-produced raster).
		if err := h.dropPreviousAvatar(user.ID); err != nil {
			slog.Error("avatar replace cleanup failed", "user", user.ID, "err", err)
		}
		h.record(user, payload)
		writeJSON(w, http.StatusOK, payload)
	})
}

// maxAvatarPerUser bounds the number of stored avatar files an account may
// upload before assigning one to the profile (W7 F-004): even if a client never
// PATCHes the profile, a single account cannot create unbounded orphan files.
// The startup GC (composition) removes unreferenced assets on restart; the cap
// bounds growth between restarts.
const maxAvatarPerUser = 10

// dropPreviousAvatar deletes the user's stored avatar file when its URL
// points into the avatar store AND the asset is owned by the user (W7 F-003:
// a URL leak must never let one account delete another user's avatar file).
func (h *accountAvatarHandler) dropPreviousAvatar(userID string) error {
	current, err := h.repository.GetUser(userID)
	if err != nil {
		return err
	}
	if id, ok := h.store.AssetIDFromURL(current.AvatarURL); ok {
		_, meta, err := h.store.load(id)
		if err == nil && strings.TrimSpace(meta["owner"]) == userID {
			if err := h.store.Delete(id); err != nil {
				return err
			}
		}
	}
	return nil
}

func (h *accountAvatarHandler) record(user account.User, payload map[string]any) {
	recordAudit(h.operations, user, operationlog.EventAccountAvatarChange, "", auditDetail("avatar-change", payload), time.Now().UTC(), nil)
}
