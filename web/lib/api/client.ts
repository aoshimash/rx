import ky, { type KyInstance } from 'ky';
import { authStore } from '@/stores/auth';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080/api/v1';

/**
 * API Client configured with authentication and error handling
 */
export const api: KyInstance = ky.create({
  prefixUrl: API_URL,
  hooks: {
    beforeRequest: [
      (request) => {
        const token = authStore.getToken();
        if (token) {
          request.headers.set('Authorization', `Bearer ${token}`);
        }
      },
    ],
    afterResponse: [
      async (_request, _options, response) => {
        if (response.status === 401) {
          // Clear invalid token
          authStore.clearToken();
          
          // Trigger auth modal
          if (typeof window !== 'undefined') {
            window.dispatchEvent(new CustomEvent('auth:required'));
          }
        }
        return response;
      },
    ],
  },
  retry: {
    limit: 1,
    methods: ['get'],
    statusCodes: [408, 413, 429, 500, 502, 503, 504],
  },
  timeout: 10000,
});

/**
 * API Error class for typed error handling
 */
export class ApiError extends Error {
  constructor(
    public code: string,
    message: string,
    public details?: Record<string, unknown>,
    public status?: number
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

/**
 * Parse API error response
 */
export async function parseApiError(response: Response): Promise<ApiError> {
  try {
    const data = await response.json();
    return new ApiError(
      data.code || 'UNKNOWN_ERROR',
      data.message || 'An error occurred',
      data.details,
      response.status
    );
  } catch {
    return new ApiError(
      'UNKNOWN_ERROR',
      `Request failed with status ${response.status}`,
      undefined,
      response.status
    );
  }
}
