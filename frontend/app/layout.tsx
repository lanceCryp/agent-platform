import type { Metadata, Viewport } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AgentHub - AI Agent 服务平台",
  description: "191个专业AI Agent，随需调用。企业级AI服务，降低90%成本。",
  keywords: ["AI Agent", "人工智能", "代码审查", "品牌设计", "营销助手"],
  authors: [{ name: "AgentHub Team" }],
  icons: {
    icon: "/favicon.ico",
  },
};

export const viewport: Viewport = {
  width: "device-width",
  initialScale: 1,
  themeColor: "#0ea5e9",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="zh-CN">
      <body className="antialiased">
        {children}
      </body>
    </html>
  );
}
