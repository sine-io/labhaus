CREATE TABLE IF NOT EXISTS styles (
  id varchar(36) PRIMARY KEY,
  name varchar(100) NOT NULL,
  description varchar(500),
  prompt text NOT NULL,
  category varchar(50),
  tags text,
  created_at timestamptz NOT NULL,
  updated_at timestamptz NOT NULL,
  deleted_at timestamptz
);

CREATE INDEX IF NOT EXISTS idx_styles_name ON styles (name);
CREATE INDEX IF NOT EXISTS idx_styles_category ON styles (category);
CREATE INDEX IF NOT EXISTS idx_styles_deleted_at ON styles (deleted_at);

INSERT INTO styles (
  id,
  name,
  description,
  prompt,
  category,
  tags,
  created_at,
  updated_at,
  deleted_at
) VALUES
  (
    '00000000-0000-4000-8000-000000000101',
    'Minimal UI Dashboard',
    'Clean modern dashboard style for SaaS and product interface images.',
    'modern minimal UI dashboard, clean layout, soft shadows, spacious composition, polished product design',
    'ui',
    '["ui","dashboard","minimal","saas","modern"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000102',
    'SaaS Landing Hero',
    'High-conversion landing page hero style with crisp product visuals.',
    'SaaS landing page hero, gradient background, floating product cards, clean typography, conversion focused',
    'ui',
    '["ui","landing","hero","saas","product"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000201',
    'Retro Poster',
    'Warm vintage poster style for nostalgic campaign visuals.',
    'retro poster illustration, warm grain texture, bold headline space, vintage print palette, nostalgic advertising',
    'retro',
    '["retro","poster","vintage","grain","advertising"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000202',
    'Pixel Product Scene',
    'Playful retro pixel-art scene for product and social content.',
    'pixel art product scene, 16-bit color palette, playful composition, crisp blocky details, retro game mood',
    'retro',
    '["retro","pixel","game","product","playful"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000301',
    'Botanical Product',
    'Natural botanical style for wellness and lifestyle products.',
    'botanical product photography, natural leaves, soft daylight, organic materials, calm wellness mood',
    'nature',
    '["nature","botanical","wellness","organic","product"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000302',
    'Soft Landscape',
    'Atmospheric nature landscape style with calm editorial framing.',
    'soft landscape scene, misty mountains, natural light, calm cinematic framing, peaceful outdoor atmosphere',
    'nature',
    '["nature","landscape","cinematic","outdoor","calm"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000401',
    'Neon Cyberpunk Interface',
    'High-contrast cyberpunk UI style with neon accents.',
    'cyberpunk interface, neon cyan and magenta lighting, dark glass panels, futuristic data display, high contrast',
    'cyberpunk',
    '["cyberpunk","neon","interface","futuristic","dark"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000402',
    'Cyberpunk City Ad',
    'Futuristic city campaign style with cinematic street lighting.',
    'cyberpunk city advertisement, rain reflections, neon signs, cinematic street scene, futuristic brand placement',
    'cyberpunk',
    '["cyberpunk","city","neon","ad","cinematic"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000501',
    'Watercolor Editorial',
    'Soft watercolor art direction for editorial and story visuals.',
    'watercolor editorial illustration, gentle paper texture, expressive brushwork, refined color harmony, storybook mood',
    'art',
    '["art","watercolor","editorial","illustration","soft"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000502',
    'Editorial Collage',
    'Layered art collage style for expressive campaign compositions.',
    'editorial collage, layered paper cutouts, mixed media texture, bold composition, contemporary art direction',
    'art',
    '["art","collage","editorial","mixed-media","campaign"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000601',
    'Luxury Editorial',
    'Premium editorial style for high-end brand campaigns.',
    'luxury editorial photography, refined lighting, elegant composition, premium materials, sophisticated brand mood',
    'luxury',
    '["luxury","editorial","premium","brand","elegant"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  ),
  (
    '00000000-0000-4000-8000-000000000602',
    'Minimal Luxury Product',
    'Quiet luxury product style with clean negative space.',
    'minimal luxury product photography, warm neutral backdrop, precise highlights, negative space, premium packaging',
    'luxury',
    '["luxury","minimal","product","packaging","premium"]',
    '2026-07-07T00:00:00Z',
    '2026-07-07T00:00:00Z',
    NULL
  )
ON CONFLICT (id) DO UPDATE SET
  name = EXCLUDED.name,
  description = EXCLUDED.description,
  prompt = EXCLUDED.prompt,
  category = EXCLUDED.category,
  tags = EXCLUDED.tags,
  created_at = EXCLUDED.created_at,
  updated_at = EXCLUDED.updated_at,
  deleted_at = NULL;
