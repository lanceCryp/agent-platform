import Link from "next/link";

export default function HomePage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-gray-50 to-gray-100">
      {/* Header */}
      <header className="border-b border-gray-200 bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <div className="flex items-center gap-2">
              <div className="w-8 h-8 bg-gradient-to-br from-sky-500 to-purple-600 rounded-lg flex items-center justify-center">
                <span className="text-white font-bold text-sm">A</span>
              </div>
              <span className="text-xl font-bold bg-gradient-to-r from-sky-600 to-purple-600 bg-clip-text text-transparent">
                AgentHub
              </span>
            </div>
            <nav className="hidden md:flex items-center gap-8">
              <Link href="/agents" className="text-gray-600 hover:text-gray-900 transition-colors">
                Agent 市场
              </Link>
              <Link href="/pricing" className="text-gray-600 hover:text-gray-900 transition-colors">
                定价
              </Link>
              <Link href="/docs" className="text-gray-600 hover:text-gray-900 transition-colors">
                文档
              </Link>
            </nav>
            <div className="flex items-center gap-4">
              <Link
                href="/login"
                className="text-gray-600 hover:text-gray-900 transition-colors"
              >
                登录
              </Link>
              <Link
                href="/register"
                className="px-4 py-2 bg-sky-500 text-white rounded-lg hover:bg-sky-600 transition-colors font-medium"
              >
                免费试用
              </Link>
            </div>
          </div>
        </div>
      </header>

      {/* Hero Section */}
      <section className="py-24 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto text-center">
          <h1 className="text-5xl sm:text-6xl lg:text-7xl font-bold tracking-tight mb-6">
            <span className="bg-gradient-to-r from-gray-900 via-gray-800 to-gray-900 bg-clip-text text-transparent">
              AI Agent
            </span>
            <br />
            <span className="bg-gradient-to-r from-sky-500 via-purple-500 to-pink-500 bg-clip-text text-transparent">
              按需调用
            </span>
          </h1>
          <p className="text-xl sm:text-2xl text-gray-600 max-w-3xl mx-auto mb-12">
            191个专业AI Agent，覆盖工程、设计、营销等全领域。
            <br />
            <span className="text-gray-500">像雇佣专家一样，按需使用，成本降低90%。</span>
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/register"
              className="px-8 py-4 bg-sky-500 text-white rounded-xl hover:bg-sky-600 transition-all font-semibold text-lg shadow-lg shadow-sky-500/25"
            >
              立即开始免费试用
            </Link>
            <Link
              href="/agents"
              className="px-8 py-4 bg-white text-gray-700 border border-gray-200 rounded-xl hover:bg-gray-50 transition-all font-semibold text-lg"
            >
              浏览 Agent 市场 →
            </Link>
          </div>
        </div>
      </section>

      {/* Features Section */}
      <section className="py-24 px-4 sm:px-6 lg:px-8 bg-white">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              为什么选择 AgentHub？
            </h2>
            <p className="text-xl text-gray-600 max-w-2xl mx-auto">
              专业、高效、低成本的企业级AI服务解决方案
            </p>
          </div>
          
          <div className="grid md:grid-cols-3 gap-8">
            {/* Feature 1 */}
            <div className="p-8 rounded-2xl border border-gray-200 hover:border-sky-200 hover:shadow-lg transition-all">
              <div className="w-12 h-12 bg-sky-100 rounded-xl flex items-center justify-center mb-6">
                <svg className="w-6 h-6 text-sky-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19.428 15.428a2 2 0 00-1.022-.547l-2.387-.477a6 6 0 00-3.86.517l-.318.158a6 6 0 01-3.86.517L6.05 15.21a2 2 0 00-1.806.547M8 4h8l-1 1v5.172a2 2 0 00.586 1.414l5 5c1.26 1.26.367 3.414-1.415 3.414H4.828c-1.782 0-2.674-2.154-1.414-3.414l5-5A2 2 0 009 10.172V5L8 4z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">191+ 专业 Agent</h3>
              <p className="text-gray-600">
                覆盖工程开发、品牌设计、营销策划、数据分析等全领域，每个Agent都是该领域的专家。
              </p>
            </div>

            {/* Feature 2 */}
            <div className="p-8 rounded-2xl border border-gray-200 hover:border-purple-200 hover:shadow-lg transition-all">
              <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center mb-6">
                <svg className="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">秒级响应</h3>
              <p className="text-gray-600">
                7×24小时在线，无需等待。提交任务后立即开始执行，快速获得专业结果。
              </p>
            </div>

            {/* Feature 3 */}
            <div className="p-8 rounded-2xl border border-gray-200 hover:border-green-200 hover:shadow-lg transition-all">
              <div className="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center mb-6">
                <svg className="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 3s1.343 3 3 3 3 .895 3 3-1.343 3-3 3-3-.895-3-3 1.343-3 3-3zm0 2c.828 0 1.5-.672 1.5-1.5S12.828 9 12 9s-1.5.672-1.5 1.5S11.172 12 12 12zm0 6c.828 0 1.5-.672 1.5-1.5S12.828 15 12 15s-1.5.672-1.5 1.5S11.172 18 12 18z" />
                </svg>
              </div>
              <h3 className="text-xl font-bold text-gray-900 mb-3">成本降低90%</h3>
              <p className="text-gray-600">
                相比雇佣全职专家或外包服务，按需订阅模式大幅降低成本，创业公司也能享受顶级AI服务。
              </p>
            </div>
          </div>
        </div>
      </section>

      {/* Categories Section */}
      <section className="py-24 px-4 sm:px-6 lg:px-8">
        <div className="max-w-7xl mx-auto">
          <div className="text-center mb-16">
            <h2 className="text-4xl font-bold text-gray-900 mb-4">
              全领域覆盖
            </h2>
            <p className="text-xl text-gray-600">
              从代码开发到品牌营销，满足企业全方位需求
            </p>
          </div>

          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { name: "工程开发", count: 45, icon: "💻", color: "bg-sky-500" },
              { name: "设计创意", count: 32, icon: "🎨", color: "bg-purple-500" },
              { name: "市场营销", count: 38, icon: "📈", color: "bg-green-500" },
              { name: "游戏开发", count: 18, icon: "🎮", color: "bg-pink-500" },
              { name: "数据分析", count: 15, icon: "📊", color: "bg-orange-500" },
              { name: "财务管理", count: 12, icon: "💰", color: "bg-emerald-500" },
              { name: "法务合规", count: 8, icon: "⚖️", color: "bg-indigo-500" },
              { name: "学术研究", count: 23, icon: "📚", color: "bg-cyan-500" },
            ].map((category) => (
              <Link
                key={category.name}
                href={`/agents?category=${category.name}`}
                className="group p-6 bg-white rounded-xl border border-gray-200 hover:border-gray-300 hover:shadow-md transition-all"
              >
                <div className={`w-12 h-12 ${category.color} rounded-xl flex items-center justify-center text-2xl mb-4 group-hover:scale-110 transition-transform`}>
                  {category.icon}
                </div>
                <h3 className="font-bold text-gray-900 mb-1">{category.name}</h3>
                <p className="text-sm text-gray-500">{category.count} 个 Agent</p>
              </Link>
            ))}
          </div>
        </div>
      </section>

      {/* CTA Section */}
      <section className="py-24 px-4 sm:px-6 lg:px-8 bg-gradient-to-br from-sky-500 to-purple-600">
        <div className="max-w-4xl mx-auto text-center text-white">
          <h2 className="text-4xl sm:text-5xl font-bold mb-6">
            准备好提升效率了吗？
          </h2>
          <p className="text-xl text-white/90 mb-12">
            加入 thousands of 企业，开始使用 AI Agent 服务
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Link
              href="/register"
              className="px-8 py-4 bg-white text-sky-600 rounded-xl hover:bg-gray-100 transition-all font-semibold text-lg"
            >
              立即免费开始
            </Link>
            <Link
              href="/contact"
              className="px-8 py-4 bg-white/10 backdrop-blur-sm border border-white/20 text-white rounded-xl hover:bg-white/20 transition-all font-semibold text-lg"
            >
              联系我们
            </Link>
          </div>
        </div>
      </section>

      {/* Footer */}
      <footer className="py-12 px-4 sm:px-6 lg:px-8 bg-gray-900 text-gray-400">
        <div className="max-w-7xl mx-auto">
          <div className="grid md:grid-cols-4 gap-8">
            <div>
              <div className="flex items-center gap-2 mb-4">
                <div className="w-8 h-8 bg-gradient-to-br from-sky-500 to-purple-600 rounded-lg flex items-center justify-center">
                  <span className="text-white font-bold text-sm">A</span>
                </div>
                <span className="text-xl font-bold text-white">AgentHub</span>
              </div>
              <p className="text-sm">
                企业级AI Agent服务平台，让专业AI能力触手可及。
              </p>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">产品</h4>
              <ul className="space-y-2 text-sm">
                <li><Link href="/agents" className="hover:text-white transition-colors">Agent 市场</Link></li>
                <li><Link href="/pricing" className="hover:text-white transition-colors">定价方案</Link></li>
                <li><Link href="/enterprise" className="hover:text-white transition-colors">企业版</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">资源</h4>
              <ul className="space-y-2 text-sm">
                <li><Link href="/docs" className="hover:text-white transition-colors">文档</Link></li>
                <li><Link href="/api" className="hover:text-white transition-colors">API</Link></li>
                <li><Link href="/blog" className="hover:text-white transition-colors">博客</Link></li>
              </ul>
            </div>
            <div>
              <h4 className="font-semibold text-white mb-4">公司</h4>
              <ul className="space-y-2 text-sm">
                <li><Link href="/about" className="hover:text-white transition-colors">关于我们</Link></li>
                <li><Link href="/contact" className="hover:text-white transition-colors">联系我们</Link></li>
                <li><Link href="/privacy" className="hover:text-white transition-colors">隐私政策</Link></li>
              </ul>
            </div>
          </div>
          <div className="border-t border-gray-800 mt-12 pt-8 text-center text-sm">
            <p>© 2026 AgentHub. All rights reserved.</p>
          </div>
        </div>
      </footer>
    </div>
  );
}
