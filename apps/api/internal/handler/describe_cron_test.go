package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	tasksstore "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/store"
)

func TestDescribeCronPatterns(t *testing.T) {
	cases := []struct {
		expr string
		zh   string
		en   string
	}{
		{"* * * * *", "每分钟", "Every minute"},
		{"*/5 * * * *", "每 5 分钟", "Every 5 minutes"},
		{"0 * * * *", "每小时的第 0 分钟", "Every hour at minute 0"},
		{"15 * * * *", "每小时的第 15 分钟", "Every hour at minute 15"},
		{"0 2 * * *", "每天 02:00", "Every day at 02:00"},
		{"30 8 * * 1", "每周周一 08:30", "Every Monday at 08:30"},
		{"0 0 1 * *", "每月 1 日 00:00", "On day 1 of every month at 00:00"},
		{"0 2 1,15 * *", "5 段 Cron 计划", "5-field cron schedule"},
	}
	for _, tc := range cases {
		fields, err := tasksstore.ParseCron(tc.expr)
		if err != nil {
			t.Fatalf("ParseCron(%q): %v", tc.expr, err)
		}
		if got := describeCron(fields, "zh-CN"); got != tc.zh {
			t.Errorf("zh %q: got %q want %q", tc.expr, got, tc.zh)
		}
		if got := describeCron(fields, "en-US"); got != tc.en {
			t.Errorf("en %q: got %q want %q", tc.expr, got, tc.en)
		}
	}
}

func TestCronPreviewNegotiatesChineseDescription(t *testing.T) {
	env := newAuthTestEnv(t)
	req := bearer(t, adminToken(t, env), http.MethodPost, "/api/scheduled-tasks/cron/preview", `{"cron":"0 2 * * *"}`)
	req.Header.Set("Accept-Language", "zh-CN")
	rr := httptest.NewRecorder()
	env.mux.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d: %s", rr.Code, rr.Body.String())
	}
	var body struct {
		Description string `json:"description"`
	}
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body.Description != "每天 02:00" {
		t.Fatalf("description = %q, want 每天 02:00", body.Description)
	}
}
