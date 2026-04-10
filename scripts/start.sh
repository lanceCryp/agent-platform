#!/bin/bash
# AgentHub 启动脚本

set -e

echo "🚀 启动 AgentHub..."

# 检查端口
check_port() {
    if lsof -Pi :$1 -sTCP:LISTEN -t >/dev/null 2>&1; then
        echo "⚠️  端口 $1 已被占用"
        return 1
    fi
    return 0
}

# 1. 启动后端 (Go API Gateway)
start_backend() {
    echo "📦 启动后端服务..."
    cd /workspace/projects/workspace/backend
    
    # 检查 Go 是否安装
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装"
        exit 1
    fi
    
    # 下载依赖
    echo "📥 下载 Go 依赖..."
    go mod download
    
    # 启动服务 (后台运行)
    nohup go run cmd/api-gateway/main.go > /tmp/agenthub-backend.log 2>&1 &
    BACKEND_PID=$!
    echo "✅ 后端服务已启动 (PID: $BACKEND_PID)"
    
    # 等待服务启动
    sleep 3
    
    # 检查是否启动成功
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ 后端服务健康检查通过"
    else
        echo "⚠️  后端服务可能未正常启动，查看日志: /tmp/agenthub-backend.log"
    fi
}

# 2. 启动前端 (Next.js)
start_frontend() {
    echo "📦 启动前端服务..."
    cd /workspace/projects/workspace/frontend
    
    # 检查 Node 是否安装
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装"
        exit 1
    fi
    
    # 安装依赖
    echo "📥 安装前端依赖..."
    npm install
    
    # 启动开发服务器 (后台运行)
    nohup npm run dev > /tmp/agenthub-frontend.log 2>&1 &
    FRONTEND_PID=$!
    echo "✅ 前端服务已启动 (PID: $FRONTEND_PID)"
    
    # 等待服务启动
    sleep 5
    
    # 检查是否启动成功
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ 前端服务健康检查通过"
    else
        echo "⚠️  前端服务可能未正常启动，查看日志: /tmp/agenthub-frontend.log"
    fi
}

# 显示帮助
show_help() {
    echo "用法: ./scripts/start.sh [选项]"
    echo ""
    echo "选项:"
    echo "  --backend    只启动后端"
    echo "  --frontend   只启动前端"
    echo "  --all        启动全部 (默认)"
    echo "  --status     查看服务状态"
    echo "  --stop       停止所有服务"
    echo "  --help       显示帮助"
}

# 查看状态
show_status() {
    echo "📊 服务状态:"
    echo ""
    
    # 后端状态
    if curl -s http://localhost:8080/health > /dev/null 2>&1; then
        echo "✅ 后端 API: http://localhost:8080 (运行中)"
    else
        echo "❌ 后端 API: http://localhost:8080 (未运行)"
    fi
    
    # 前端状态
    if curl -s http://localhost:3000 > /dev/null 2>&1; then
        echo "✅ 前端 Web: http://localhost:3000 (运行中)"
    else
        echo "❌ 前端 Web: http://localhost:3000 (未运行)"
    fi
}

# 停止服务
stop_services() {
    echo "🛑 停止所有服务..."
    
    # 停止后端
    pkill -f "go run cmd/api-gateway" && echo "✅ 后端已停止" || echo "⚠️  后端未运行"
    
    # 停止前端
    pkill -f "next dev" && echo "✅ 前端已停止" || echo "⚠️  前端未运行"
    
    echo "✅ 所有服务已停止"
}

# 主逻辑
case "${1:-all}" in
    --backend)
        start_backend
        ;;
    --frontend)
        start_frontend
        ;;
    --all)
        start_backend
        start_frontend
        echo ""
        echo "🎉 AgentHub 启动完成!"
        echo ""
        show_status
        ;;
    --status)
        show_status
        ;;
    --stop)
        stop_services
        ;;
    --help)
        show_help
        ;;
    *)
        echo "❌ 未知选项: $1"
        show_help
        exit 1
        ;;
esac
