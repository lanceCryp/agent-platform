"use client";

import * as React from "react";
import Link from "next/link";
import { useParams } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";

// Mock data - in real app, fetch from API
const mockAgent = {
  id: "1",
  agent_id: "engineering-code-reviewer",
  name: "代码审查专家",
  name_en: "Code Review Expert",
  description: "专业代码审查，发现潜在bug、安全漏洞和性能问题。提供详细的改进建议。支持多种编程语言，包括Python、JavaScript、TypeScript、Go、Rust等。",
  category: "engineering",
  category_name: "工程开发",
  category_icon: "💻",
  tags: ["代码审查", "质量", "安全", "性能"],
  tier: 1,
  tier_name: "基础",
  runtime_type: "claude",
  price_per_request: 9.90,
  avg_duration_seconds: 45,
  success_rate: 98.5,
  rating: 4.8,
  total_tasks: 15234,
  total_success: 15005,
  total_failed: 229,
  input_example: "请帮我审查这段Python代码的潜在问题:\n\ndef calculate(a, b):\n    return a / b",
  output_example: "✅ 代码审查完成\n\n发现 1 个潜在问题:\n\n**除零风险** (Line 2)\n- 问题: 当 b 为 0 时会抛出 ZeroDivisionError\n- 建议: 添加除零检查",
};

export default function AgentDetailPage() {
  const params = useParams();
  const agentId = params.agentId as string;
  const [prompt, setPrompt] = React.useState("");
  const [isLoading, setIsLoading] = React.useState(false);
  const [result, setResult] = React.useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!prompt.trim()) return;

    setIsLoading(true);
    setResult(null);

    // Simulate API call
    await new Promise((resolve) => setTimeout(resolve, 2000));
    
    setResult(mockAgent.output_example);
    setIsLoading(false);
  };

  const tierColors: Record<number, string> = {
    1: "bg-green-100 text-green-700",
    2: "bg-blue-100 text-blue-700",
    3: "bg-purple-100 text-purple-700",
    4: "bg-amber-100 text-amber-700",
  };

  return (
    <div className="max-w-6xl mx-auto">
      {/* Breadcrumb */}
      <div className="flex items-center gap-2 text-sm text-gray-500 mb-6">
        <Link href="/agents" className="hover:text-sky-600">Agent 市场</Link>
        <span>/</span>
        <span>{mockAgent.category_name}</span>
        <span>/</span>
        <span className="text-gray-900">{mockAgent.name}</span>
      </div>

      <div className="grid lg:grid-cols-3 gap-6">
        {/* Left: Agent Info */}
        <div className="lg:col-span-2 space-y-6">
          {/* Header */}
          <div className="bg-white rounded-xl border border-gray-200 p-6">
            <div className="flex items-start justify-between mb-4">
              <div>
                <div className="flex items-center gap-3 mb-2">
                  <span className="text-4xl">{mockAgent.category_icon}</span>
                  <div>
                    <h1 className="text-2xl font-bold text-gray-900">{mockAgent.name}</h1>
                    <p className="text-gray-500">{mockAgent.name_en}</p>
                  </div>
                </div>
              </div>
              <span className={`px-3 py-1 rounded-full text-sm font-medium ${tierColors[mockAgent.tier]}`}>
                {mockAgent.tier_name}
              </span>
            </div>

            <p className="text-gray-600 mb-4">{mockAgent.description}</p>

            {/* Tags */}
            <div className="flex flex-wrap gap-2 mb-6">
              {mockAgent.tags.map((tag) => (
                <span key={tag} className="px-3 py-1 bg-gray-100 text-gray-600 rounded-full text-sm">
                  {tag}
                </span>
              ))}
            </div>

            {/* Stats */}
            <div className="grid grid-cols-4 gap-4 p-4 bg-gray-50 rounded-lg">
              <div className="text-center">
                <p className="text-2xl font-bold text-sky-600">{mockAgent.rating}</p>
                <p className="text-sm text-gray-500">评分</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-green-600">{mockAgent.success_rate}%</p>
                <p className="text-sm text-gray-500">成功率</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-purple-600">{mockAgent.total_tasks.toLocaleString()}</p>
                <p className="text-sm text-gray-500">使用次数</p>
              </div>
              <div className="text-center">
                <p className="text-2xl font-bold text-gray-600">{mockAgent.avg_duration_seconds}s</p>
                <p className="text-sm text-gray-500">平均耗时</p>
              </div>
            </div>
          </div>

          {/* Examples */}
          <Card>
            <CardHeader>
              <CardTitle>使用示例</CardTitle>
            </CardHeader>
            <CardContent className="space-y-4">
              <div>
                <p className="text-sm font-medium text-gray-500 mb-2">输入示例</p>
                <div className="p-4 bg-gray-900 text-gray-100 rounded-lg font-mono text-sm whitespace-pre-wrap">
                  {mockAgent.input_example}
                </div>
              </div>
              <div>
                <p className="text-sm font-medium text-gray-500 mb-2">输出示例</p>
                <div className="p-4 bg-gray-50 border border-gray-200 rounded-lg text-sm whitespace-pre-wrap">
                  {mockAgent.output_example}
                </div>
              </div>
            </CardContent>
          </Card>
        </div>

        {/* Right: Task Form */}
        <div className="lg:col-span-1">
          <div className="bg-white rounded-xl border border-gray-200 p-6 sticky top-24">
            <div className="text-center mb-6">
              <span className="text-4xl font-bold text-sky-600">¥{mockAgent.price_per_request}</span>
              <span className="text-gray-500">/次</span>
            </div>

            <form onSubmit={handleSubmit} className="space-y-4">
              <div>
                <label className="block text-sm font-medium text-gray-700 mb-2">
                  输入提示词
                </label>
                <textarea
                  value={prompt}
                  onChange={(e) => setPrompt(e.target.value)}
                  placeholder="描述您需要完成的任务..."
                  rows={6}
                  className="w-full px-3 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none"
                />
              </div>

              <Button type="submit" className="w-full" size="lg" isLoading={isLoading}>
                {isLoading ? "处理中..." : "立即使用"}
              </Button>
            </form>

            {/* Result */}
            {result && (
              <div className="mt-6">
                <p className="text-sm font-medium text-gray-700 mb-2">执行结果</p>
                <div className="p-4 bg-gray-50 border border-gray-200 rounded-lg text-sm whitespace-pre-wrap max-h-96 overflow-y-auto">
                  {result}
                </div>
              </div>
            )}

            {/* Tips */}
            <div className="mt-6 pt-6 border-t border-gray-200">
              <p className="text-sm text-gray-500">
                💡 提示：详细描述任务需求可获得更准确的结果
              </p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
