"use client";

import { useEffect, useState } from "react";

import Link from "next/link";
import { usePathname } from "next/navigation";

import { useTranslations } from "next-intl";

import { Button } from "@/components/ui/button";

export function OnboardingEntry() {
  const t = useTranslations("onboarding");
  return (
    <Button variant="outline" asChild>
      <Link href="/onboarding">{t("title")}</Link>
    </Button>
  );
}

export function OnboardingReturn() {
  const pathname = usePathname();
  const [visible, setVisible] = useState(false);
  const t = useTranslations("onboarding");
  // biome-ignore lint/correctness/useExhaustiveDependencies: Read the URL again after client-side navigation.
  useEffect(() => {
    setVisible(new URLSearchParams(window.location.search).get("from") === "onboarding");
  }, [pathname]);
  if (!visible) return null;
  return (
    <div className="mb-4">
      <Button variant="outline" asChild>
        <Link href="/onboarding">{t("back")}</Link>
      </Button>
    </div>
  );
}
