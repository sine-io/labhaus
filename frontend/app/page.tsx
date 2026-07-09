import Link from 'next/link';
import { ImageIcon, Sparkles, ArrowRight } from 'lucide-react';

export default function HomePage() {
  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-16 sm:px-6 lg:px-8">
      {/* Hero Section */}
      <div className="text-center">
        <div className="inline-flex items-center gap-2 rounded-full glass-card px-4 py-2 text-sm font-bold text-accent-cyan-light mb-6">
          <Sparkles className="h-4 w-4" />
          <span>AI 内容生产工作流实验室</span>
        </div>

        <h1 className="text-5xl md:text-7xl font-black tracking-tight mb-6">
          批量生产
          <br />
          <span className="text-accent-cyan">AI 创意内容</span>
        </h1>

        <p className="text-lg text-text-secondary max-w-2xl mx-auto mb-12">
          LabHaus 提供可视化工作流编排，让专业团队高效批量生产 AI
          图像和视频内容。探索我们的核心功能：
        </p>

        {/* Feature Cards */}
        <div className="grid md:grid-cols-2 gap-6 max-w-4xl mx-auto">
          <Link
            href="/images/generate"
            className="group glass-card rounded-2xl p-8 hover:border-border-glow transition-all"
          >
            <div className="flex items-center justify-center w-14 h-14 rounded-xl bg-accent-cyan/10 mb-4 mx-auto group-hover:bg-accent-cyan/20 transition-colors">
              <ImageIcon className="h-7 w-7 text-accent-cyan" />
            </div>
            <h2 className="text-2xl font-black mb-3">批量生图</h2>
            <p className="text-text-secondary mb-4">
              一次输入多个 prompts，批量生成高质量图像，自动存储到 MinIO 对象存储
            </p>
            <div className="flex items-center justify-center gap-2 text-accent-cyan font-bold group-hover:gap-3 transition-all">
              <span>开始生成</span>
              <ArrowRight className="h-4 w-4" />
            </div>
          </Link>

          <Link
            href="/styles/recommend"
            className="group glass-card rounded-2xl p-8 hover:border-border-glow transition-all"
          >
            <div className="flex items-center justify-center w-14 h-14 rounded-xl bg-accent-cyan/10 mb-4 mx-auto group-hover:bg-accent-cyan/20 transition-colors">
              <Sparkles className="h-7 w-7 text-accent-cyan" />
            </div>
            <h2 className="text-2xl font-black mb-3">样式推荐</h2>
            <p className="text-text-secondary mb-4">
              基于 TF-IDF 和余弦相似度算法，智能推荐最匹配的图像样式和风格
            </p>
            <div className="flex items-center justify-center gap-2 text-accent-cyan font-bold group-hover:gap-3 transition-all">
              <span>探索样式</span>
              <ArrowRight className="h-4 w-4" />
            </div>
          </Link>
        </div>
      </div>

      {/* Status Badge */}
      <div className="mt-16 text-center">
        <div className="inline-flex items-center gap-2 text-sm text-text-muted">
          <div className="h-2 w-2 rounded-full bg-green-400 animate-pulse"></div>
          <span>Phase 2.1 已完成 · 后端 API 已就绪</span>
        </div>
      </div>
    </div>
  );
}
