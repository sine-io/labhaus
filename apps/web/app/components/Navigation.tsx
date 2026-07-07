'use client';

import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ImageIcon, Sparkles, Home, KeyRound } from 'lucide-react';

export function Navigation() {
  const pathname = usePathname();

  const navItems = [
    { href: '/', label: '首页', icon: Home },
    { href: '/auth', label: '认证', icon: KeyRound },
    { href: '/images/generate', label: '批量生图', icon: ImageIcon },
    { href: '/styles/recommend', label: '样式推荐', icon: Sparkles },
  ];

  return (
    <header className="topbar">
      <div className="mx-auto w-full max-w-7xl px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          {/* Logo */}
          <Link href="/" className="flex items-center gap-2">
            <div className="flex h-9 w-9 items-center justify-center rounded-lg glass-card">
              <Sparkles className="h-5 w-5 text-accent-cyan" />
            </div>
            <span className="text-lg font-black tracking-tight">LabHaus</span>
          </Link>

          {/* Navigation */}
          <nav className="flex items-center gap-1 rounded-lg glass-card px-1 py-1">
            {navItems.map((item) => {
              const Icon = item.icon;
              const isActive = pathname === item.href;
              return (
                <Link
                  key={item.href}
                  href={item.href}
                  className={`
                    flex items-center gap-2 rounded-md px-3 py-2 text-sm font-bold transition-colors
                    ${
                      isActive
                        ? 'bg-white/10 text-white'
                        : 'text-text-secondary hover:bg-white/8 hover:text-white'
                    }
                  `}
                >
                  <Icon className="h-4 w-4" />
                  {item.label}
                </Link>
              );
            })}
          </nav>
        </div>
      </div>
    </header>
  );
}
