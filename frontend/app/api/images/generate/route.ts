import { NextRequest, NextResponse } from 'next/server';
import {
  buildBackendHeaders,
  normalizeGenerateImageResponse,
} from '../../_lib/backend-contract.mjs';

const BACKEND_URL = process.env.BACKEND_URL || 'http://localhost:8080';

export async function POST(request: NextRequest) {
  try {
    const body = await request.json();

    const response = await fetch(`${BACKEND_URL}/api/images/generate`, {
      method: 'POST',
      headers: buildBackendHeaders(request),
      body: JSON.stringify(body),
    });

    const data = await response.json();

    if (!response.ok) {
      return NextResponse.json(
        { error: data.error || '生成图像失败' },
        { status: response.status }
      );
    }

    return NextResponse.json(normalizeGenerateImageResponse(data));
  } catch (error: unknown) {
    console.error('API proxy error:', error);
    return NextResponse.json(
      { error: error instanceof Error ? error.message : '服务器错误' },
      { status: 500 }
    );
  }
}
