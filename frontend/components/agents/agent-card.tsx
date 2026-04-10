"use client";

import * as React from "react";
import Link from "next/link";
import type { Agent } from "@/types/api";
import { Button } from "@/components/ui/button";

interface AgentCardProps {
  agent: Agent;
}

export function AgentCard({ agent }: AgentCardProps) {
  const tierColors = {
    1: "bg-green-100 text-green-700",
    2: "bg-blue-100 text-blue-700",
    3: "bg-purple-100 text-purple-700",
    4: "bg-amber-100 text-amber-700",
  };

  const tierNames = {
    1: "基础",
    2: "专业",
    3: "高级",
    4: "旗舰",
  };

  const runtimeIcons = {
    openclaw: "🦞",
    claude: "🧠",
    openai: "🤖",
  };

  return (
    <Link href={`/agents/${agent.agent_id}`}>
      <div className="group h-full bg-white rounded-xl border border-gray-200 p-5 hover:border-sky-200 hover:shadow-lg transition-all cursor-pointer">
        {/* Header */}
        <div className="flex items-start justify-between mb-4">
          <div className="flex-1">
            <div className="flex items-center gap-2 mb-1">
              <span className="text-2xl">{runtimeIcons[agent.runtime_type]}</span>
              <span className={`px-2 py-0.5 rounded-full text-xs font-medium ${tierColors[agent.tier]}`}>
                {tierNames[agent.tier]}
              </span>
            </div>
            <h3 className="font-semibold text-gray-900 group-hover:text-sky-600 transition-colors">
              {agent.name}
            </h3>
            {agent.name_en && (
              <p className="text-sm text-gray-500">{agent.name_en}</p>
            )}
          </div>
        </div>

        {/* Description */}
        <p className="text-sm text-gray-600 line-clamp-2 mb-4">
          {agent.description}
        </p>

        {/* Tags */}
        <div className="flex flex-wrap gap-1.5 mb-4">
          {agent.tags.slice(0, 3).map((tag) => (
            <span
              key={tag}
              className="px-2 py-0.5 bg-gray-100 text-gray-600 rounded text-xs"
            >
              {tag}
            </span>
          ))}
          {agent.tags.length > 3 && (
            <span className="px-2 py-0.5 bg-gray-100 text-gray-500 rounded text-xs">
              +{agent.tags.length - 3}
            </span>
          )}
        </div>

        {/* Stats */}
        <div className="flex items-center justify-between text-sm border-t border-gray-100 pt-4">
          <div className="flex items-center gap-4">
            <div className="flex items-center gap-1 text-gray-500">
              <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
              <span className="font-medium">{agent.rating.toFixed(1)}</span>
            </div>
            <div className="text-gray-500">
              <span className={agent.success_rate >= 95 ? "text-green-600" : "text-amber-600"}>
                {agent.success_rate}%
              </span>
              成功率
            </div>
          </div>
        </div>

        {/* Price */}
        <div className="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
          <div>
            <span className="text-2xl font-bold text-sky-600">¥{agent.price_per_request.toFixed(2)}</span>
            <span className="text-sm text-gray-500">/次</span>
          </div>
          <Button size="sm" variant="outline" className="group-hover:bg-sky-50 group-hover:border-sky-300">
            立即使用 →
          </Button>
        </div>
      </div>
    </Link>
  );
}

// Agent Grid Component
interface AgentGridProps {
  agents: Agent[];
  loading?: boolean;
}

export function AgentGrid({ agents, loading }: AgentGridProps) {
  if (loading) {
    return (
      <div className="grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
        {Array.from({ length: 8 }).map((_, i) => (
          <div key={i} className="bg-white rounded-xl border border-gray-200 p-5 animate-pulse">
            <div className="flex items-center gap-2 mb-4">
              <div className="w-8 h-8 bg-gray-200 rounded-lg" />
              <div className="flex-1">
                <div className="h-4 bg-gray-200 rounded w-3/4 mb-2" />
                <div className="h-3 bg-gray-200 rounded w-1/2" />
              </div>
            </div>
            <div className="space-y-2 mb-4">
              <div className="h-3 bg-gray-200 rounded" />
              <div className="h-3 bg-gray-200 rounded w-2/3" />
            </div>
            <div className="flex gap-2 mb-4">
              <div className="h-6 bg-gray-200 rounded w-16" />
              <div className="h-6 bg-gray-200 rounded w-16" />
            </div>
            <div className="h-8 bg-gray-200 rounded" />
          </div>
        ))}
      </div>
    );
  }

  if (agents.length === 0) {
    return (
      <div className="text-center py-16">
        <div className="text-6xl mb-4">🔍</div>
        <h3 className="text-xl font-semibold text-gray-900 mb-2">没有找到匹配的 Agent</h3>
        <p className="text-gray-500">尝试调整筛选条件或搜索关键词</p>
      </div>
    );
  }

  return (
    <div className="grid sm:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
      {agents.map((agent) => (
        <AgentCard key={agent.id} agent={agent} />
      ))}
    </div>
  );
}
