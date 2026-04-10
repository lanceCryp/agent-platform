"use client";

import * as React from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";

type Tab = "profile" | "security" | "notifications" | "api";

export default function SettingsPage() {
  const [activeTab, setActiveTab] = React.useState<Tab>("profile");
  const [isSaving, setIsSaving] = React.useState(false);

  // Profile form state
  const [profile, setProfile] = React.useState({
    username: "johndoe",
    email: "john@example.com",
    phone: "+86 138 8888 8888",
    avatar: "",
  });

  // Security form state
  const [security, setSecurity] = React.useState({
    currentPassword: "",
    newPassword: "",
    confirmPassword: "",
  });

  // Notification settings
  const [notifications, setNotifications] = React.useState({
    emailTaskComplete: true,
    emailTaskFailed: true,
    emailMarketing: false,
    pushTaskComplete: true,
    pushTaskFailed: true,
    pushNewsletter: false,
  });

  // API Keys
  const [apiKeys, setApiKeys] = React.useState([
    { id: "1", name: "Production Key", key: "sk_live_xxxx...xxxx", created: "2026-01-15" },
    { id: "2", name: "Development Key", key: "sk_test_xxxx...xxxx", created: "2026-03-20" },
  ]);

  const handleSaveProfile = async () => {
    setIsSaving(true);
    await new Promise((resolve) => setTimeout(resolve, 1000));
    setIsSaving(false);
  };

  const tabs = [
    { id: "profile" as const, label: "个人信息", icon: "👤" },
    { id: "security" as const, label: "安全设置", icon: "🔒" },
    { id: "notifications" as const, label: "通知设置", icon: "🔔" },
    { id: "api" as const, label: "API Keys", icon: "🔑" },
  ];

  return (
    <div className="max-w-4xl mx-auto">
      {/* Header */}
      <div className="mb-8">
        <h1 className="text-2xl font-bold text-gray-900 mb-2">个人设置</h1>
        <p className="text-gray-500">管理您的账户信息和偏好设置</p>
      </div>

      <div className="flex gap-8">
        {/* Sidebar */}
        <div className="w-48 shrink-0">
          <nav className="space-y-1">
            {tabs.map((tab) => (
              <button
                key={tab.id}
                onClick={() => setActiveTab(tab.id)}
                className={`w-full flex items-center gap-3 px-3 py-2 rounded-lg text-left transition-colors ${
                  activeTab === tab.id
                    ? "bg-sky-50 text-sky-600"
                    : "text-gray-600 hover:bg-gray-100"
                }`}
              >
                <span>{tab.icon}</span>
                <span className="font-medium">{tab.label}</span>
              </button>
            ))}
          </nav>
        </div>

        {/* Content */}
        <div className="flex-1">
          {/* Profile Tab */}
          {activeTab === "profile" && (
            <Card>
              <CardHeader>
                <CardTitle>个人信息</CardTitle>
                <CardDescription>更新您的个人资料信息</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                {/* Avatar */}
                <div className="flex items-center gap-6">
                  <div className="w-20 h-20 rounded-full bg-gradient-to-br from-sky-400 to-purple-500 flex items-center justify-center">
                    <span className="text-white text-2xl font-bold">J</span>
                  </div>
                  <div>
                    <Button variant="outline" size="sm">
                      上传头像
                    </Button>
                    <p className="text-sm text-gray-500 mt-1">支持 JPG、PNG，最大 2MB</p>
                  </div>
                </div>

                <div className="grid gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">用户名</label>
                    <Input
                      value={profile.username}
                      onChange={(e) => setProfile({ ...profile, username: e.target.value })}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">邮箱</label>
                    <Input
                      type="email"
                      value={profile.email}
                      onChange={(e) => setProfile({ ...profile, email: e.target.value })}
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">手机号</label>
                    <Input
                      type="tel"
                      value={profile.phone}
                      onChange={(e) => setProfile({ ...profile, phone: e.target.value })}
                    />
                  </div>
                </div>

                <Button onClick={handleSaveProfile} isLoading={isSaving}>
                  保存更改
                </Button>
              </CardContent>
            </Card>
          )}

          {/* Security Tab */}
          {activeTab === "security" && (
            <Card>
              <CardHeader>
                <CardTitle>安全设置</CardTitle>
                <CardDescription>管理您的密码和安全设置</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div className="grid gap-4">
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">当前密码</label>
                    <Input
                      type="password"
                      value={security.currentPassword}
                      onChange={(e) => setSecurity({ ...security, currentPassword: e.target.value })}
                      placeholder="输入当前密码"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">新密码</label>
                    <Input
                      type="password"
                      value={security.newPassword}
                      onChange={(e) => setSecurity({ ...security, newPassword: e.target.value })}
                      placeholder="输入新密码"
                    />
                  </div>
                  <div>
                    <label className="block text-sm font-medium text-gray-700 mb-2">确认新密码</label>
                    <Input
                      type="password"
                      value={security.confirmPassword}
                      onChange={(e) => setSecurity({ ...security, confirmPassword: e.target.value })}
                      placeholder="再次输入新密码"
                    />
                  </div>
                </div>

                <Button>更新密码</Button>

                <div className="pt-6 border-t border-gray-200">
                  <h3 className="font-medium text-gray-900 mb-4">两步验证</h3>
                  <div className="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
                    <div>
                      <p className="font-medium">未启用</p>
                      <p className="text-sm text-gray-500">启用后登录需要输入手机验证码</p>
                    </div>
                    <Button variant="outline">启用</Button>
                  </div>
                </div>
              </CardContent>
            </Card>
          )}

          {/* Notifications Tab */}
          {activeTab === "notifications" && (
            <Card>
              <CardHeader>
                <CardTitle>通知设置</CardTitle>
                <CardDescription>选择您希望接收的通知类型</CardDescription>
              </CardHeader>
              <CardContent className="space-y-6">
                <div>
                  <h3 className="font-medium text-gray-900 mb-4">邮件通知</h3>
                  <div className="space-y-4">
                    {[
                      { key: "emailTaskComplete", label: "任务完成", desc: "当任务执行完成时发送邮件" },
                      { key: "emailTaskFailed", label: "任务失败", desc: "当任务执行失败时发送邮件" },
                      { key: "emailMarketing", label: "营销邮件", desc: "接收产品更新和优惠信息" },
                    ].map((item) => (
                      <div key={item.key} className="flex items-center justify-between">
                        <div>
                          <p className="font-medium">{item.label}</p>
                          <p className="text-sm text-gray-500">{item.desc}</p>
                        </div>
                        <label className="relative inline-flex items-center cursor-pointer">
                          <input
                            type="checkbox"
                            checked={notifications[item.key as keyof typeof notifications]}
                            onChange={(e) => setNotifications({ ...notifications, [item.key]: e.target.checked })}
                            className="sr-only peer"
                          />
                          <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-sky-100 rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
                        </label>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="pt-6 border-t border-gray-200">
                  <h3 className="font-medium text-gray-900 mb-4">推送通知</h3>
                  <div className="space-y-4">
                    {[
                      { key: "pushTaskComplete", label: "任务完成" },
                      { key: "pushTaskFailed", label: "任务失败" },
                      { key: "pushNewsletter", label: "Newsletter" },
                    ].map((item) => (
                      <div key={item.key} className="flex items-center justify-between">
                        <p className="font-medium">{item.label}</p>
                        <label className="relative inline-flex items-center cursor-pointer">
                          <input
                            type="checkbox"
                            checked={notifications[item.key as keyof typeof notifications]}
                            onChange={(e) => setNotifications({ ...notifications, [item.key]: e.target.checked })}
                            className="sr-only peer"
                          />
                          <div className="w-11 h-6 bg-gray-200 peer-focus:outline-none rounded-full peer peer-checked:after:translate-x-full rtl:peer-checked:after:-translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:start-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-sky-500"></div>
                        </label>
                      </div>
                    ))}
                  </div>
                </div>

                <Button>保存设置</Button>
              </CardContent>
            </Card>
          )}

          {/* API Tab */}
          {activeTab === "api" && (
            <Card>
              <CardHeader>
                <CardTitle>API Keys</CardTitle>
                <CardDescription>管理您的 API Keys，用于集成 AgentHub 到您的应用</CardDescription>
              </CardHeader>
              <CardContent className="space-y-4">
                <div className="flex justify-end">
                  <Button>
                    <svg className="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                    </svg>
                    创建新 Key
                  </Button>
                </div>

                <div className="space-y-3">
                  {apiKeys.map((apiKey) => (
                    <div key={apiKey.id} className="p-4 bg-gray-50 rounded-lg border border-gray-200">
                      <div className="flex items-center justify-between mb-2">
                        <p className="font-medium">{apiKey.name}</p>
                        <div className="flex gap-2">
                          <Button variant="ghost" size="sm">复制</Button>
                          <Button variant="ghost" size="sm" className="text-red-600">删除</Button>
                        </div>
                      </div>
                      <p className="text-sm text-gray-500 font-mono mb-2">{apiKey.key}</p>
                      <p className="text-xs text-gray-400">创建于 {apiKey.created}</p>
                    </div>
                  ))}
                </div>

                <div className="p-4 bg-sky-50 rounded-lg border border-sky-200">
                  <h4 className="font-medium text-sky-900 mb-2">使用指南</h4>
                  <p className="text-sm text-sky-700 mb-3">
                    在请求头中包含您的 API Key：
                  </p>
                  <code className="block p-3 bg-sky-100 rounded text-sm font-mono">
                    Authorization: Bearer YOUR_API_KEY
                  </code>
                </div>
              </CardContent>
            </Card>
          )}
        </div>
      </div>
    </div>
  );
}
