package config

import (
	"strings"
	"testing"
)

// R1 checkpoint 3 (workspace-014 GOAL-002 D-001): storage.objects surface.
func TestObjectsConfig(t *testing.T) {
	t.Run("defaults to local with path-style on", func(t *testing.T) {
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.ObjectsDriver != "local" || !cfg.ObjectsS3UsePathStyle {
			t.Fatalf("objects defaults = driver %q pathStyle %v", cfg.ObjectsDriver, cfg.ObjectsS3UsePathStyle)
		}
	})

	t.Run("yaml s3 block parses", func(t *testing.T) {
		y := "app:\n  env: development\nstorage:\n  objects:\n    driver: s3\n    local:\n      root: /tmp/objects\n    s3:\n      endpoint: http://minio:9000\n      region: us-east-1\n      bucket: schema-ui\n      access_key_id: AKIAEXAMPLE\n      secret_access_key: shhh\n      use_path_style: false\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.ObjectsDriver != "s3" || cfg.ObjectsLocalRoot != "/tmp/objects" {
			t.Fatalf("driver/root = %q/%q", cfg.ObjectsDriver, cfg.ObjectsLocalRoot)
		}
		if cfg.ObjectsS3Endpoint != "http://minio:9000" || cfg.ObjectsS3Region != "us-east-1" || cfg.ObjectsS3Bucket != "schema-ui" ||
			cfg.ObjectsS3AccessKeyID != "AKIAEXAMPLE" || cfg.ObjectsS3SecretAccessKey != "shhh" || cfg.ObjectsS3UsePathStyle {
			t.Fatalf("s3 block = %+v pathStyle %v", cfg.ObjectsS3Endpoint, cfg.ObjectsS3UsePathStyle)
		}
	})

	t.Run("env overrides win over yaml", func(t *testing.T) {
		t.Setenv("STORAGE_OBJECTS_DRIVER", "s3")
		t.Setenv("STORAGE_OBJECTS_S3_ENDPOINT", "http://env-endpoint:9000")
		t.Setenv("STORAGE_OBJECTS_S3_BUCKET", "env-bucket")
		t.Setenv("STORAGE_OBJECTS_S3_ACCESS_KEY_ID", "env-key")
		t.Setenv("STORAGE_OBJECTS_S3_SECRET_ACCESS_KEY", "env-secret")
		y := "app:\n  env: development\nstorage:\n  objects:\n    driver: local\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError != nil {
			t.Fatalf("LoadError: %v", cfg.LoadError)
		}
		if cfg.ObjectsDriver != "s3" || cfg.ObjectsS3Endpoint != "http://env-endpoint:9000" || cfg.ObjectsS3Bucket != "env-bucket" {
			t.Fatalf("env override = %s/%s/%s", cfg.ObjectsDriver, cfg.ObjectsS3Endpoint, cfg.ObjectsS3Bucket)
		}
	})

	t.Run("unknown driver fails closed at load", func(t *testing.T) {
		y := "app:\n  env: development\nstorage:\n  objects:\n    driver: gcs\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "storage.objects.driver") {
			t.Fatalf("unknown driver must be a LoadError, got %v", cfg.LoadError)
		}
	})

	t.Run("s3 keys with driver=local fail closed (misconfig interception)", func(t *testing.T) {
		y := "app:\n  env: development\nstorage:\n  objects:\n    driver: local\n    s3:\n      endpoint: http://minio:9000\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "driver is local") {
			t.Fatalf("s3 keys under local driver must fail closed, got %v", cfg.LoadError)
		}
	})

	t.Run("driver=s3 missing secret fails closed at load", func(t *testing.T) {
		y := "app:\n  env: development\nstorage:\n  objects:\n    driver: s3\n    s3:\n      endpoint: http://minio:9000\n      bucket: b\n      access_key_id: k\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil || !strings.Contains(cfg.LoadError.Error(), "secret_access_key") {
			t.Fatalf("missing s3 secret must fail closed, got %v", cfg.LoadError)
		}
	})

	t.Run("ValidateProd re-checks zero-value configs", func(t *testing.T) {
		c := &Config{AppEnv: "development", ObjectsDriver: "s3"}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "storage.objects.s3.endpoint") {
			t.Fatalf("zero-value s3 config must fail ValidateProd, got %v", err)
		}
		c2 := &Config{AppEnv: "development"}
		if err := c2.ValidateProd(); err != nil {
			t.Fatalf("default (empty) objects driver must pass, got %v", err)
		}
	})
}

// A-002 F-001: the local-misconfig error names the KEY, never the value.
func TestObjectsMisconfigErrorDoesNotLeakSecret(t *testing.T) {
	t.Run("secret-only under local driver", func(t *testing.T) {
		const secret = "supersecret-value-42"
		y := "app:\n  env: development\nstorage:\n  objects:\n    s3:\n      secret_access_key: " + secret + "\n"
		writeConfig(t, y)
		cfg := Load()
		if cfg.LoadError == nil {
			t.Fatal("secret-only s3 key with local driver must fail closed")
		}
		if !strings.Contains(cfg.LoadError.Error(), "storage.objects.s3.secret_access_key") {
			t.Fatalf("error must name the offending key, got %v", cfg.LoadError)
		}
		if strings.Contains(cfg.LoadError.Error(), secret) {
			t.Fatal("error must never carry the secret value (A-002 F-001)")
		}
	})

	t.Run("ValidateProd re-checks local misconfig (A-002 R-001)", func(t *testing.T) {
		c := &Config{AppEnv: "development", ObjectsS3Endpoint: "http://minio:9000"}
		if err := c.ValidateProd(); err == nil || !strings.Contains(err.Error(), "storage.objects.s3.endpoint") {
			t.Fatalf("ValidateProd must re-check local+s3 misconfig, got %v", err)
		}
	})
}
