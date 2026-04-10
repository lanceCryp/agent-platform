"use client";

import * as React from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { AgentGrid } from "@/components/agents/agent-card";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

// Mock data for demonstration
const mockAgents = [
  {
    id: "1",
    agent_id: "engineering-code-reviewer",
    name: "代码审查专家",
    name_en: "Code Review Expert",
    description: "专业代码审查，发现潜在bug、安全漏洞和性能问题。提供详细的改进建议。",
    category: "engineering",
    tags: ["代码审查", "质量", "安全"],
    tier: 1,
    runtime_type: "claude" as const,
    price_per_request: 9.90,
    avg_duration_seconds: 45,
    success_rate: 98.5,
    rating: 4.8,
    total_tasks: 15234,
    total_success: 15005,
    total_failed: 229,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "2",
    agent_id: "design-brand-guardian",
    name: "品牌视觉设计师",
    name_en: "Brand Guardian",
    description: "确保品牌视觉一致性，提供设计建议和品牌规范指南。",
    category: "design",
    tags: ["品牌设计", "视觉", "规范"],
    tier: 2,
    runtime_type: "openclaw" as const,
    price_per_request: 29.90,
    avg_duration_seconds: 120,
    success_rate: 96.2,
    rating: 4.9,
    total_tasks: 8756,
    total_success: 8421,
    total_failed: 335,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "3",
    agent_id: "marketing-content-writer",
    name: "营销文案专家",
    name_en: "Marketing Content Writer",
    description: "创作高质量营销文案，包括广告、社交媒体、产品描述等。",
    category: "marketing",
    tags: ["文案创作", "营销", "SEO"],
    tier: 1,
    runtime_type: "claude" as const,
    price_per_request: 14.90,
    avg_duration_seconds: 60,
    success_rate: 97.8,
    rating: 4.7,
    total_tasks: 23456,
    total_success: 22940,
    total_failed: 516,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "4",
    agent_id: "data-analysis-expert",
    name: "数据分析专家",
    name_en: "Data Analysis Expert",
    description: "深度数据分析，挖掘业务洞察，提供数据驱动的决策建议。",
    category: "data",
    tags: ["数据分析", "BI", "洞察"],
    tier: 2,
    runtime_type: "openai" as const,
    price_per_request: 39.90,
    avg_duration_seconds: 180,
    success_rate: 95.5,
    rating: 4.6,
    total_tasks: 6543,
    total_success: 6249,
    total_failed: 294,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "5",
    agent_id: "game-narrative-designer",
    name: "游戏叙事设计师",
    name_en: "Game Narrative Designer",
    description: "创作引人入胜的游戏故事、角色背景和世界观设定。",
    category: "gaming",
    tags: ["游戏设计", "叙事", "创意"],
    tier: 3,
    runtime_type: "claude" as const,
    price_per_request: 49.90,
    avg_duration_seconds: 240,
    success_rate: 94.2,
    rating: 4.9,
    total_tasks: 3421,
    total_success: 3223,
    total_failed: 198,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "6",
    agent_id: "finance-risk-analyst",
    name: "财务风险分析师",
    name_en: "Financial Risk Analyst",
    description: "评估投资风险，分析财务报表，提供风险管理建议。",
    category: "finance",
    tags: ["财务管理", "风险", "投资"],
    tier: 4,
    runtime_type: "openclaw" as const,
    price_per_request: 99.90,
    avg_duration_seconds: 300,
    success_rate: 93.8,
    rating: 4.8,
    total_tasks: 2134,
    total_success: 2002,
    total_failed: 132,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "7",
    agent_id: "legal-contract-reviewer",
    name: "合同审查律师",
    name_en: "Contract Reviewer",
    description: "专业合同审查，识别法律风险，提供修改建议。",
    category: "legal",
    tags: ["法务", "合同", "合规"],
    tier: 3,
    runtime_type: "claude" as const,
    price_per_request: 69.90,
    avg_duration_seconds: 180,
    success_rate: 96.5,
    rating: 4.7,
    total_tasks: 4567,
    total_success: 4407,
    total_failed: 160,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
  {
    id: "8",
    agent_id: "research-paper-writer",
    name: "学术论文助手",
    name_en: "Research Paper Assistant",
    description: "协助学术论文写作，包括文献综述、研究方法、结论建议。",
    category: "academic",
    tags: ["学术", "论文", "研究"],
    tier: 2,
    runtime_type: "openai" as const,
    price_per_request: 34.90,
    avg_duration_seconds: 600,
    success_rate: 94.0,
    rating: 4.5,
    total_tasks: 9876,
    total_success: 9283,
    total_failed: 593,
    is_active: true,
    created_at: "2026-01-01",
    updated_at: "2026-04-01",
  },
];

const categories = [
  { name: "全部", count: 191, icon: "🌐" },
  { name: "工程开发", count: 45, icon: "💻" },
  { name: "设计创意", count: 32, icon: "🎨" },
  { name: "市场营销", count: 38, icon: "📈" },
  { name: "数据分析", count: 15, icon: "📊" },
  { name: "游戏开发", count: 18, icon: "🎮" },
  { name: "财务管理", count: 12, icon: "💰" },
  { name: "法务合规", count: 8, icon: "⚖️" },
  { name: "学术研究", count: 23, icon: "📚" },
];

const tiers = [
  { value: 0, label: "全部等级" },
  { value: 1, label: "🥉 基础" },
  { value: 2, label: "🥈 专业" },
  { value: 3, label: "🥇 高级" },
  { value: 4, label: "👑 旗舰" },
];

export default function AgentsPage() {
  const searchParams = useSearchParams();
  const categoryParam = searchParams.get("category");
  
  const [searchQuery, setSearchQuery] = React.useState("");
  const [selectedCategory, setSelectedCategory] = React.useState(categoryParam || "全部");
  const [selectedTier, setSelectedTier] = React.useState(0);
  const [sortBy, setSortBy] = React.useState<"popular" | "rating" | "price-low" | "price-high">("popular");
  const [isLoading, setIsLoading] = React.useState(false);

  // Filter agents
  const filteredAgents = React.useMemo(() => {
    let result = [...mockAgents];
    
    // Category filter
    if (selectedCategory !== "全部") {
      result = result.filter((agent) => agent.category === selectedCategory);
    }
    
    // Tier filter
    if (selectedTier > 0) {
      result = result.filter((agent) => agent.tier === selectedTier);
    }
    
    // Search filter
    if (searchQuery) {
      const query = searchQuery.toLowerCase();
      result = result.filter(
        (agent) =>
          agent.name.toLowerCase().includes(query) ||
          agent.description.toLowerCase().includes(query) ||
          agent.tags.some((tag) => tag.toLowerCase().includes(query))
      );
    }
    
    // Sort
    switch (sortBy) {
      case "rating":
        result.sort((a, b) => b.rating - a.rating);
        break;
      case "price-low":
        result.sort((a, b) => a.price_per_request - b.price_per_request);
        break;
      case "price-high":
        result.sort((a, b) => b.price_per_request - a.price_per_request);
        break;
      case "popular":
      default:
        result.sort((a, b) => b.total_tasks - a.total_tasks);
    }
    
    return result;
  }, [selectedCategory, selectedTier, searchQuery, sortBy]);

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Header */}
      <header className="bg-white border-b border-gray-200 sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-sky-500 to-purple-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">A</span>
              </div>
              <span className="text-xl font-bold bg-gradient-to-r from-sky-600 to-purple-600 bg-clip-text text-transparent">
                AgentHub
              </span>
            </Link>
            <div className="flex items-center gap-4">
              <Link href="/login">
                <Button variant="ghost">登录</Button>
              </Link>
              <Link href="/register">
                <Button>免费试用</Button>
              </Link>
            </div>
          </div>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold text-gray-900 mb-2">Agent 市场</h1>
          <p className="text-gray-600">
            浏览 191 个专业 AI Agent，找到适合您业务的解决方案
          </p>
        </div>

        {/* Filters */}
        <div className="bg-white rounded-xl border border-gray-200 p-4 mb-8">
          <div className="flex flex-col lg:flex-row gap-4">
            {/* Search */}
            <div className="flex-1">
              <Input
                placeholder="搜索 Agent 名称、描述或标签..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                icon={
                  <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                  </svg>
                }
              />
            </div>

            {/* Sort */}
            <div className="flex gap-2">
              <select
                value={sortBy}
                onChange={(e) => setSortBy(e.target.value as typeof sortBy)}
                className="h-10 px-3 rounded-lg border border-gray-300 bg-white text-sm focus:outline-none focus:ring-2 focus:ring-sky-500"
              >
                <option value="popular">🔥 热门</option>
                <option value="rating">⭐ 评分最高</option>
                <option value="price-low">💰 价格从低到高</option>
                <option value="price-high">💎 价格从高到低</option>
              </select>
            </div>
          </div>

          {/* Category Pills */}
          <div className="flex flex-wrap gap-2 mt-4">
            {categories.map((cat) => (
              <button
                key={cat.name}
                onClick={() => setSelectedCategory(cat.name)}
                className={`px-4 py-2 rounded-full text-sm font-medium transition-all ${
                  selectedCategory === cat.name
                    ? "bg-sky-500 text-white"
                    : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                }`}
              >
                {cat.icon} {cat.name} ({cat.count})
              </button>
            ))}
          </div>

          {/* Tier Filter */}
          <div className="flex items-center gap-4 mt-4 pt-4 border-t border-gray-100">
            <span className="text-sm text-gray-500">等级筛选:</span>
            <div className="flex gap-2">
              {tiers.map((tier) => (
                <button
                  key={tier.value}
                  onClick={() => setSelectedTier(tier.value)}
                  className={`px-3 py-1.5 rounded-lg text-sm transition-all ${
                    selectedTier === tier.value
                      ? "bg-purple-500 text-white"
                      : "bg-gray-100 text-gray-600 hover:bg-gray-200"
                  }`}
                >
                  {tier.label}
                </button>
              ))}
            </div>
          </div>
        </div>

        {/* Results */}
        <div className="flex items-center justify-between mb-6">
          <p className="text-gray-600">
            找到 <span className="font-semibold text-gray-900">{filteredAgents.length}</span> 个 Agent
          </p>
        </div>

        <AgentGrid agents={filteredAgents} loading={isLoading} />

        {/* Load More */}
        {filteredAgents.length > 0 && (
          <div className="text-center mt-8">
            <Button variant="outline" size="lg">
              加载更多 Agent
            </Button>
          </div>
        )}
      </div>
    </div>
  );
}
