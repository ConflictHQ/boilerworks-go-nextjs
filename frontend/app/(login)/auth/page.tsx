"use client";

export default function LoginRedirect() {
  if (typeof window !== "undefined") {
    window.location.href = "/auth/login";
  }
  return null;
}
