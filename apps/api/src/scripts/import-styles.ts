#!/usr/bin/env node
import fs from 'fs/promises';
import path from 'path';
import { fileURLToPath } from 'url';
import { StyleRepository } from '../repositories/style.repository.js';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

interface CaseData {
  id: number;
  title: string;
  image: string;
  imageAlt: string;
  sourceLabel: string;
  sourceUrl: string;
  prompt: string;
  promptPreview: string;
  category: string;
  styles: string[];
  scenes: string[];
  featured: boolean;
  githubUrl: string;
}

interface CasesJson {
  repository: string;
  totalCases: number;
  cases: CaseData[];
}

async function importStylesFromAwesomeGptImage2() {
  console.log('Starting style import from awesome-gpt-image-2...');

  const casesPath = process.env.CASES_JSON_PATH || 
    path.join(__dirname, '../../../..', 'awesome-gpt-image-2/dist/cases.json');

  console.log(`Reading cases from: ${casesPath}`);

  try {
    const data = await fs.readFile(casesPath, 'utf-8');
    const casesJson: CasesJson = JSON.parse(data);

    console.log(`Found ${casesJson.totalCases} cases`);

    const styleRepo = new StyleRepository();
    const styles = casesJson.cases.map((c) => ({
      case_id: c.id,
      title: c.title,
      prompt: c.prompt,
      prompt_preview: c.promptPreview || null,
      category: c.category,
      styles: c.styles || [],
      scenes: c.scenes || [],
      image_url: c.image || null,
      source_label: c.sourceLabel || null,
      source_url: c.sourceUrl || null,
      github_url: c.githubUrl || null,
      featured: c.featured || false,
    }));

    console.log('Inserting styles into database...');
    const inserted = await styleRepo.batchCreate(styles);

    console.log(`✓ Successfully imported ${inserted} styles`);
    process.exit(0);
  } catch (error) {
    console.error('Import failed:', error);
    process.exit(1);
  }
}

importStylesFromAwesomeGptImage2();
