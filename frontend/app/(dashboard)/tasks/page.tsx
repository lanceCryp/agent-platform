"use client";

import * as React from "react";
import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";

// Mock tasks data
const mockTasks = [
  {
    id: "1",
    agent_name: "代码审查专家",
    prompt: "请帮我审查以下Python代码...",
    status: "completed",
    cost: 9.90,
    created_at: "2026-04-10 10:30",
    completed_at: "2026-04-10 10:31",
  },
  {
    id: "2",
    agent_name: "营销文案专家",
    prompt: "为智能手表撰写产品描述...",
    status: "processing",
    cost: 14.90,
    created_at: "2026-04-10 10:25",
    completed_at: null,
  },
  {
    id: "3",
    agent_name: "品牌视觉设计师",
    prompt: "检查LOGO设计的品牌一致性...",
    status: "pending",
    cost: 29.90,
    created_at: "2026-04-10 10:20",
    completed_at: null,
  },
  {
    id: "4",
    agent_name: "数据分析专家",
    prompt: "分析本季度的销售数据...",
    status: "failed",
    cost: 39.90,
    created_at: "2026-04-10 09:15",
    completed_at: null,
    error: "任务超时，请重试",
  },
  {
    id: "5",
    agent_name: "代码审查专家",
    prompt: "代码安全审查...",
    status: "completed",
    cost: 9.90,
    created_at: "2026-04-09 15:30",
    completed_at: "2026-04-09 15:32",
  },
];

const statusConfig: Record<string, { color: string; label: string; bg: string }> = {
  pending: { color: "text-amber-600", label: "等待中", bg: "bg-amber-100" },
  processing: { color: "text-sky-600", label: "处理中", bg: "bg-sky-100" },
  completed: { color: "text-green-600", label: "已完成", bg: "bg-green-100" },
  failed: { color: "text-red-600", label: "失败", bg: "bg-red-100" },
  cancelled: { color: "text-gray-600", label: "已取消", bg: "bg-gray-100" },
};

export default function TasksPage() {
  const [statusFilter, setStatusFilter] = React.useState<string>("all");
  const [tasks, setTasks] = React.useState(mockTasks);

  const filteredTasks = statusFilter === "all"
    ? tasks
    : tasks.filter((t) => t.status === statusFilter);

  const getStatusActions = (task: typeof mockTasks[0]) => {
    switch (task.status) {
      case "pending":
        return (
          <Button variant="outline" size="sm">
            取消
          </Button>
        );
      case "failed":
        return (
          <Button variant="outline" size="sm">
            重试
          </Button>
        );
      case "completed":
        return (
          <Link href={`/tasks/${task.id}`}>
            <Button variant="outline" size="sm">
              查看结果
            </Button>
          </Link>
        );
      default:
        return null;
    }
  };

  return (
    <div className="max-w-6xl mx-auto">
      {/* Header */}
      <div className="flex items-center justify-between mb-6">
        <div>
          <h1 className="text-2xl font-bold text-gray-900">我的任务</h1>
          <p className="text-gray-500">查看和管理您的所有任务</p>
        </div>
        <Link href="/tasks/new">
          <Button>
            <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
            </svg>
            新建任务
          </Button>
        </Link>
      </div>

      {/* Filters */}
      <div className="flex gap-2 mb-6">
        {[
          { value: "all", label: "全部" },
          { value: "pending", label: "等待中" },
          { value: "processing", label: "处理中" },
          { value: "completed", label: "已完成" },
          { value: "failed", label: "失败" },
        ].map((filter) => (
          <button
            key={filter.value}
            onClick={() => setStatusFilter(filter.value)}
            className={`px-4 py-2 rounded-lg text-sm font-medium transition-colors ${
              statusFilter === filter.value
                ? "bg-sky-500 text-white"
                : "bg-white text-gray-600 hover:bg-gray-100 border border-gray-200"
            }`}
          >
            {filter.label}
          </button>
        ))}
      </div>

      {/* Task List */}
      <div className="space-y-4">
        {filteredTasks.length === 0 ? (
          <Card>
            <CardContent className="py-16 text-center">
              <div className="text-6xl mb-4">📋</div>
              <h3 className="text-xl font-semibold text-gray-900 mb-2">暂无任务</h3>
              <p className="text-gray-500 mb-6">开始创建您的第一个任务吧</p>
              <Link href="/tasks/new">
                <Button>创建任务</Button>
              </Link>
            </CardContent>
          </Card>
        ) : (
          filteredTasks.map((task) => {
            const status = statusConfig[task.status];
            return (
              <Card key={task.id} className="hover:shadow-md transition-shadow">
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex-1">
                      {/* Agent & Status */}
                      <div className="flex items-center gap-3 mb-2">
                        <span className="font-medium text-gray-900">{task.agent_name}</span>
                        <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${status.bg} ${status.color}`}>
                          {status.label}
                        </span>
                      </div>

                      {/* Prompt */}
                      <p className="text-gray-600 text-sm mb-3 line-clamp-2">
                        {task.prompt}
                      </p>

                      {/* Meta */}
                      <div className="flex items-center gap-4 text-sm text-gray-500">
                        <span className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          {task.created_at}
                        </span>
                        <span className="flex items-center gap-1">
                          <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          ¥{task.cost}
                        </span>
                      </div>

                      {/* Error Message */}
                      {task.error && (
                        <p className="mt-2 text-sm text-red-600">{task.error}</p>
                      )}
                    </div>

                    {/* Actions */}
                    <div className="ml-4">
                      {getStatusActions(task)}
                    </div>
                  </div>
                </CardContent>
              </Card>
            );
          })
        )}
      </div>

      {/* Pagination */}
      {filteredTasks.length > 0 && (
        <div className="flex items-center justify-center gap-2 mt-8">
          <Button variant="outline" size="sm" disabled>
            上一页
          </Button>
          <span className="px-4 py-2 text-sm text-gray-600">
            第 1 / 10 页
          </span>
          <Button variant="outline" size="sm">
            下一页
          </Button>
        </div>
      )}
    </div>
  );
}
