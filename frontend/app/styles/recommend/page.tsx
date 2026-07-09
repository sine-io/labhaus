'use client';

import { useState } from 'react';
import { Sparkles, Search, Loader2, Tag, X } from 'lucide-react';
import { buildFrontendHeaders, getStoredBearerToken } from '../../lib/auth-token.mjs';

interface StyleRecommendation {
  name: string;
  description: string;
  prompt: string;
  category: string;
  tags: string[];
  score: number;
}

export default function RecommendStylesPage() {
  const [query, setQuery] = useState('');
  const [loading, setLoading] = useState(false);
  const [styles, setStyles] = useState<StyleRecommendation[]>([]);
  const [error, setError] = useState('');

  const handleSearch = async () => {
    if (!query.trim()) {
      setError('请输入搜索关键词');
      return;
    }

    setLoading(true);
    setError('');
    setStyles([]);

    try {
      const response = await fetch('/api/styles/recommend', {
        method: 'POST',
        headers: buildFrontendHeaders(getStoredBearerToken()),
        body: JSON.stringify({
          query: query.trim(),
          limit: 10,
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || '推荐失败');
      }

      const data = await response.json();
      setStyles(data.recommendations || []);
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '获取推荐时出错');
    } finally {
      setLoading(false);
    }
  };

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !loading) {
      handleSearch();
    }
  };

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg glass-card">
            <Sparkles className="h-5 w-5 text-accent-cyan" />
          </div>
          <h1 className="text-3xl font-black">样式推荐</h1>
        </div>
        <p className="text-text-secondary">基于 TF-IDF 和余弦相似度，智能推荐最匹配的图像样式</p>
      </div>

      {/* Search Section */}
      <div className="glass-card rounded-2xl p-6 mb-6">
        <label className="block text-sm font-bold mb-3">搜索关键词</label>
        <div className="flex gap-3">
          <input
            type="text"
            value={query}
            onChange={(e) => setQuery(e.target.value)}
            onKeyPress={handleKeyPress}
            placeholder="例如：cyberpunk city, anime character, landscape..."
            className="flex-1 px-4 py-3 rounded-xl glass-input text-sm"
            disabled={loading}
          />
          <button
            onClick={handleSearch}
            disabled={loading || !query.trim()}
            className="glass-button px-6 py-3 rounded-xl font-bold flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                搜索中
              </>
            ) : (
              <>
                <Search className="h-4 w-4" />
                搜索
              </>
            )}
          </button>
        </div>
      </div>

      {/* Error Message */}
      {error && (
        <div className="glass-card rounded-xl p-4 mb-6 border-red-500/50 bg-red-500/10">
          <div className="flex items-start gap-3">
            <X className="h-5 w-5 text-red-400 flex-shrink-0 mt-0.5" />
            <div>
              <p className="font-bold text-red-400 mb-1">推荐失败</p>
              <p className="text-sm text-red-300">{error}</p>
            </div>
          </div>
        </div>
      )}

      {/* Results */}
      {styles.length > 0 && (
        <div>
          <h2 className="text-xl font-black mb-4">推荐结果 ({styles.length})</h2>
          <div className="space-y-4">
            {styles.map((style, idx) => (
              <div key={idx} className="glass-card rounded-xl p-6">
                <div className="flex items-start justify-between mb-3">
                  <div className="flex-1">
                    <h3 className="text-xl font-black mb-1">{style.name}</h3>
                    <p className="text-sm text-text-secondary">{style.description}</p>
                  </div>
                  <div className="flex items-center gap-2 ml-4">
                    <div className="text-right">
                      <div className="text-2xl font-black text-accent-cyan">
                        {(style.score * 100).toFixed(0)}%
                      </div>
                      <div className="text-xs text-text-muted">相似度</div>
                    </div>
                  </div>
                </div>

                {/* Prompt */}
                <div className="mb-4 p-4 rounded-lg bg-white/5">
                  <p className="text-xs font-bold text-text-muted mb-2">PROMPT</p>
                  <p className="text-sm font-mono text-text-primary">{style.prompt}</p>
                </div>

                {/* Tags */}
                <div className="flex flex-wrap gap-2">
                  <span className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-accent-cyan/10 text-accent-cyan text-xs font-bold">
                    {style.category}
                  </span>
                  {style.tags.map((tag, tagIdx) => (
                    <span
                      key={tagIdx}
                      className="inline-flex items-center gap-1 px-3 py-1 rounded-full bg-white/5 text-text-secondary text-xs font-bold"
                    >
                      <Tag className="h-3 w-3" />
                      {tag}
                    </span>
                  ))}
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {!loading && styles.length === 0 && !error && (
        <div className="glass-card rounded-2xl p-12 text-center">
          <Sparkles className="h-16 w-16 text-text-muted mx-auto mb-4 opacity-50" />
          <p className="text-text-muted">输入关键词并点击搜索开始</p>
        </div>
      )}
    </div>
  );
}
