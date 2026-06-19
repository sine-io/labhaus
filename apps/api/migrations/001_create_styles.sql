-- 样式库表
CREATE TABLE IF NOT EXISTS styles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id INTEGER UNIQUE NOT NULL,
    title TEXT NOT NULL,
    prompt TEXT NOT NULL,
    prompt_preview TEXT,
    category TEXT NOT NULL,
    styles TEXT[] DEFAULT '{}',
    scenes TEXT[] DEFAULT '{}',
    image_url TEXT,
    source_label TEXT,
    source_url TEXT,
    github_url TEXT,
    featured BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- 索引
CREATE INDEX IF NOT EXISTS idx_styles_category ON styles(category);
CREATE INDEX IF NOT EXISTS idx_styles_featured ON styles(featured);
CREATE INDEX IF NOT EXISTS idx_styles_case_id ON styles(case_id);

-- 全文搜索索引
CREATE INDEX IF NOT EXISTS idx_styles_title_search ON styles USING gin(to_tsvector('english', title));
CREATE INDEX IF NOT EXISTS idx_styles_prompt_search ON styles USING gin(to_tsvector('english', prompt));

-- 更新时间戳触发器
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_styles_updated_at BEFORE UPDATE ON styles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
