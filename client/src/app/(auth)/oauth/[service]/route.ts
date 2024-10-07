import { apiBaseUrl } from "@/shared/api-interceptor";
import { type NextRequest, NextResponse } from "next/server";

export function GET(
  _: NextRequest,
  { params }: { params: { service: string } },
) {
  return NextResponse.redirect(`${apiBaseUrl}/auth/${params.service}`);
}
