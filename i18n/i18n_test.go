package i18n

import "testing"

func TestResolveLocalePriority(t *testing.T) {
	got := ResolveLocale("en-US", "zh-CN,zh;q=0.9", "zh-CN")
	if got != LocaleENUS {
		t.Fatalf("expected preferred locale en-US, got %s", got)
	}

	got = ResolveLocale("", "en-US,en;q=0.9", "zh-CN")
	if got != LocaleENUS {
		t.Fatalf("expected browser locale en-US, got %s", got)
	}

	got = ResolveLocale("", "", "zh-CN")
	if got != LocaleZHCN {
		t.Fatalf("expected fallback locale zh-CN, got %s", got)
	}
}

func TestNormalizeLocale(t *testing.T) {
	if got := NormalizeLocale("zh"); got != LocaleZHCN {
		t.Fatalf("expected zh-CN, got %s", got)
	}
	if got := NormalizeLocale("en_GB"); got != LocaleENUS {
		t.Fatalf("expected en-US, got %s", got)
	}
}

func TestTranslate(t *testing.T) {
	if got := T("en-US", "common.home", ""); got != "Home" {
		t.Fatalf("expected Home, got %s", got)
	}
	if got := T("zh-CN", "common.home", ""); got != "首页" {
		t.Fatalf("expected 首页, got %s", got)
	}
}
