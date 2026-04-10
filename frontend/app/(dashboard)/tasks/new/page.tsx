"use client";

import * as React from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

// Mock agents for selection
const mockAgents = [
  { agent_id: "engineering-code-reviewer", name: "代码审查专家", category: "工程开发", price: 9.90, tier: 1 },
  { agent_id: "design-brand-guardian", name: "品牌视觉设计师", category: "设计创意", price: 29.90, tier: 2 },
  { agent_id: "marketing-content-writer", name: "营销文案专家", category: "市场营销", price: 14.90, tier: 1 },
  { agent_id: "data-analysis-expert", name: "数据分析专家", category: "数据分析", price: 39.90, tier: 2 },
  { agent_id: "game-narrative-designer", name: "游戏叙事设计师", category: "游戏开发", price: 49.90, tier: 3 },
  { agent_id: "finance-risk-analyst", name: "财务风险分析师", category: "财务管理", price: 99.90, tier: 4 },
];

export default function NewTaskPage() {
  const router = useRouter();
  const [selectedAgent, setSelectedAgent] = React.useState<typeof mockAgents[0] | null>(null);
  const [prompt, setPrompt] = React.useState("");
  const [priority, setPriority] = React.useState(5);
  const [isLoading, setIsLoading] = React.useState(false);
  const [showAgentSelector, setShowAgentSelector] = React.useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!selectedAgent || !prompt.trim()) return;

    setIsLoading(true);

    try {
      // Simulate API call
      await new Promise((resolve) => setTimeout(resolve, 1500));
      
      // In real implementation:
      // const response = await api.tasks.create({
      //   agent_id: selectedAgent.agent_id,
      //   prompt,
      //   priority,
      // });
      // router.push(`/tasks/${response.data.data.task_id}`);
      
      router.push("/tasks");
    } catch (error) {
      console.error("Failed to create task:", error);
    } finally {
      setIsLoading(false);
    }
  };

  const estimatedCost = selectedAgent ? selectedAgent.price : 0;

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">创建新任务</h1>
        <p className="text-gray-500">选择 Agent 并描述您的任务需求</p>
      </div>

      <form onSubmit={handleSubmit} className="space-y-6">
        {/* Agent Selection */}
        <Card>
          <CardHeader>
            <CardTitle>选择 Agent</CardTitle>
          </CardHeader>
          <CardContent>
            {selectedAgent ? (
              <div className="flex items-center justify-between p-4 bg-sky-50 rounded-lg border border-sky-200">
                <div>
                  <p className="font-medium text-gray-900">{selectedAgent.name}</p>
                  <p className="text-sm text-gray-500">{selectedAgent.category}</p>
                </div>
                <div className="flex items-center gap-4">
                  <span className="text-lg font-bold text-sky-600">¥{selectedAgent.price}/次</span>
                  <Button
                    type="button"
                    variant="outline"
                    size="sm"
                    onClick={() => setShowAgentSelector(true)}
                  >
                    更换
                  </Button>
                </div>
              </div>
            ) : (
              <Button
                type="button"
                variant="outline"
                className="w-full h-20 text-gray-500"
                onClick={() => setShowAgentSelector(true)}
              >
                <div className="flex flex-col items-center gap-1">
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                  </svg>
                  <span>点击选择 Agent</span>
                </div>
              </Button>
            )}

            {/* Agent Selector Modal */}
            {showAgentSelector && (
              <div className="mt-4 p-4 bg-gray-50 rounded-lg border border-gray-200 max-h-96 overflow-y-auto">
                <div className="grid gap-3">
                  {mockAgents.map((agent) => (
                    <button
                      key={agent.agent_id}
                      type="button"
                      onClick={() => {
                        setSelectedAgent(agent);
                        setShowAgentSelector(false);
                      }}
                      className="flex items-center justify-between p-3 bg-white rounded-lg border border-gray-200 hover:border-sky-300 hover:bg-sky-50 transition-colors text-left"
                    >
                      <div>
                        <p className="font-medium text-gray-900">{agent.name}</p>
                        <p className="text-sm text-gray-500">{agent.category}</p>
                      </div>
                      <span className="font-medium text-sky-600">¥{agent.price}/次</span>
                    </button>
                  ))}
                </div>
              </div>
            )}
          </CardContent>
        </Card>

        {/* Task Input */}
        <Card>
          <CardHeader>
            <CardTitle>任务描述</CardTitle>
          </CardHeader>
          <CardContent className="space-y-4">
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                提示词 (Prompt)
              </label>
              <textarea
                value={prompt}
                onChange={(e) => setPrompt(e.target.value)}
                placeholder="详细描述您需要完成的任务...
                
示例:
请帮我审查以下代码的潜在问题:
def calculate(a, b):
    return a / b"
                rows={8}
                className="w-full px-3 py-2 rounded-lg border border-gray-300 focus:outline-none focus:ring-2 focus:ring-sky-500 focus:border-transparent resize-none font-mono text-sm"
              />
              <p className="mt-2 text-sm text-gray-500">
                {prompt.length} / 10,000 字符
              </p>
            </div>

            {/* Priority */}
            <div>
              <label className="block text-sm font-medium text-gray-700 mb-2">
                优先级 ({priority})
              </label>
              <input
                type="range"
                min="1"
                max="10"
                value={priority}
                onChange={(e) => setPriority(parseInt(e.target.value))}
                className="w-full"
              />
              <div className="flex justify-between text-xs text-gray-500 mt-1">
                <span>低 (1)</span>
                <span>高 (10)</span>
              </div>
            </div>
          </CardContent>
        </Card>

        {/* Cost Summary */}
        <Card className="bg-gray-50">
          <CardContent className="p-6">
            <div className="flex items-center justify-between mb-4">
              <span className="text-gray-600">Agent 费用</span>
              <span className="font-medium">¥{estimatedCost.toFixed(2)}</span>
            </div>
            <div className="flex items-center justify-between mb-4">
              <span className="text-gray-600">当前余额</span>
              <span className="font-medium text-green-600">¥100.00</span>
            </div>
            <div className="border-t border-gray-200 pt-4 flex items-center justify-between">
              <span className="text-lg font-medium">预估费用</span>
              <span className="text-2xl font-bold text-sky-600">¥{estimatedCost.toFixed(2)}</span>
            </div>
          </CardContent>
        </Card>

        {/* Actions */}
        <div className="flex gap-4">
          <Link href="/tasks" className="flex-1">
            <Button type="button" variant="outline" className="w-full" size="lg">
              取消
            </Button>
          </Link>
          <Button
            type="submit"
            className="flex-1"
            size="lg"
            isLoading={isLoading}
            disabled={!selectedAgent || !prompt.trim()}
          >
            提交任务
          </Button>
        </div>
      </form>
    </div>
  );
}
