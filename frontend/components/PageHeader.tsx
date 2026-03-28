"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import { ChevronRight } from "lucide-react";
import { routeLabels } from "@/lib/routes";

export function PageHeader() {
  const pathname = usePathname();
  const segments = pathname.split("/").filter(Boolean);

  return (
    <header className="border-b border-border bg-background px-6 py-4">
      <nav className="flex items-center gap-1 text-sm text-muted-foreground">
        {segments.map((segment, index) => {
          const href = "/" + segments.slice(0, index + 1).join("/");
          const label =
            routeLabels[segment] ||
            segment.charAt(0).toUpperCase() + segment.slice(1).replace(/-/g, " ");
          const isLast = index === segments.length - 1;

          return (
            <span key={href} className="flex items-center gap-1">
              {index > 0 && <ChevronRight className="h-3 w-3" />}
              {isLast ? (
                <span className="text-foreground font-medium">{label}</span>
              ) : (
                <Link href={href} className="hover:text-foreground transition-colors">
                  {label}
                </Link>
              )}
            </span>
          );
        })}
      </nav>
    </header>
  );
}
