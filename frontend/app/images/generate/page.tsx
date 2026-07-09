'use client';

import { useState } from 'react';
import { ImageIcon, Loader2, Download, X } from 'lucide-react';
import { buildFrontendHeaders, getStoredBearerToken } from '../../lib/auth-token.mjs';

interface GeneratedImage {
  id?: string;
  url: string;
  prompt: string;
}

interface GenerateImagesResponse {
  images?: GeneratedImage[];
}

export default function GenerateImagesPage() {
  const [prompts, setPrompts] = useState('');
  const [loading, setLoading] = useState(false);
  const [images, setImages] = useState<GeneratedImage[]>([]);
  const [error, setError] = useState('');

  const handleGenerate = async () => {
    if (!prompts.trim()) {
      setError('请输入至少一个 prompt');
      return;
    }

    setLoading(true);
    setError('');
    setImages([]);

    try {
      const promptList = prompts
        .split('\n')
        .map((p) => p.trim())
        .filter((p) => p.length > 0);

      const response = await fetch('/api/images/generate', {
        method: 'POST',
        headers: buildFrontendHeaders(getStoredBearerToken()),
        body: JSON.stringify({
          prompts: promptList,
          width: 1024,
          height: 1024,
          quality: 'standard',
        }),
      });

      if (!response.ok) {
        const data = await response.json();
        throw new Error(data.error || '生成失败');
      }

      const data = (await response.json()) as GenerateImagesResponse;
      setImages(
        (data.images ?? []).map((img, idx) => ({
          id: img.id,
          url: img.url,
          prompt: img.prompt || promptList[idx],
        }))
      );
    } catch (err: unknown) {
      setError(err instanceof Error ? err.message : '生成图像时出错');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="mx-auto w-full max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
      {/* Header */}
      <div className="mb-8">
        <div className="flex items-center gap-3 mb-3">
          <div className="flex h-10 w-10 items-center justify-center rounded-lg glass-card">
            <ImageIcon className="h-5 w-5 text-accent-cyan" />
          </div>
          <h1 className="text-3xl font-black">批量生图</h1>
        </div>
        <p className="text-text-secondary">输入多个 prompts（每行一个），一次性批量生成图像</p>
      </div>

      {/* Input Section */}
      <div className="glass-card rounded-2xl p-6 mb-6">
        <label className="block text-sm font-bold mb-3">
          Prompts <span className="text-text-muted">(每行一个)</span>
        </label>
        <textarea
          value={prompts}
          onChange={(e) => setPrompts(e.target.value)}
          placeholder={`a serene mountain landscape at sunset\na futuristic city with flying cars\na cute robot gardening flowers`}
          className="w-full h-48 px-4 py-3 rounded-xl glass-input resize-none font-mono text-sm"
          disabled={loading}
        />

        <div className="flex items-center justify-between mt-4">
          <p className="text-xs text-text-muted">
            {prompts.split('\n').filter((p) => p.trim()).length} 个 prompts
          </p>
          <button
            onClick={handleGenerate}
            disabled={loading || !prompts.trim()}
            className="glass-button px-6 py-3 rounded-xl font-bold flex items-center gap-2 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            {loading ? (
              <>
                <Loader2 className="h-4 w-4 animate-spin" />
                生成中...
              </>
            ) : (
              <>
                <ImageIcon className="h-4 w-4" />
                生成图像
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
              <p className="font-bold text-red-400 mb-1">生成失败</p>
              <p className="text-sm text-red-300">{error}</p>
            </div>
          </div>
        </div>
      )}

      {/* Results Grid */}
      {images.length > 0 && (
        <div>
          <h2 className="text-xl font-black mb-4">生成结果 ({images.length})</h2>
          <div className="grid gap-6 sm:grid-cols-2 lg:grid-cols-3">
            {images.map((image, idx) => (
              <div key={idx} className="glass-card rounded-xl overflow-hidden">
                <div className="aspect-square bg-bg-glass relative">
                  <img src={image.url} alt={image.prompt} className="w-full h-full object-cover" />
                </div>
                <div className="p-4">
                  <p className="text-sm text-text-secondary mb-3 line-clamp-2">{image.prompt}</p>
                  <a
                    href={image.url}
                    download
                    target="_blank"
                    rel="noopener noreferrer"
                    className="glass-button w-full px-4 py-2 rounded-lg text-sm font-bold flex items-center justify-center gap-2"
                  >
                    <Download className="h-4 w-4" />
                    下载
                  </a>
                </div>
              </div>
            ))}
          </div>
        </div>
      )}

      {/* Empty State */}
      {!loading && images.length === 0 && !error && (
        <div className="glass-card rounded-2xl p-12 text-center">
          <ImageIcon className="h-16 w-16 text-text-muted mx-auto mb-4 opacity-50" />
          <p className="text-text-muted">输入 prompts 并点击生成按钮开始</p>
        </div>
      )}
    </div>
  );
}
